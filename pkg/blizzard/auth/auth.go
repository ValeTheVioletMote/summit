package auth

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"strings"
	"sync"

	"github.com/paalgyula/summit/pkg/blizzard/auth/packets"
	"github.com/paalgyula/summit/pkg/blizzard/auth/srp"
	"github.com/paalgyula/summit/pkg/db"
	"github.com/paalgyula/summit/pkg/wow"
	"github.com/paalgyula/summit/server/auth/data/static"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func StartServer(listenAddress string) error {

	l, err := net.Listen("tcp4", listenAddress)
	if err != nil {
		return fmt.Errorf("auth.StartServer: %w", err)
	}
	defer l.Close()

	log.Info().Msgf("auth server is listening on: %s", listenAddress)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Error().Err(err).Msgf("cannot accept connection")
		}

		NewClient(c)
	}
}

type RealmClient struct {
	c net.Conn

	outLock sync.Mutex

	log zerolog.Logger

	account *db.Account

	// Keys for authentication
	PrivateEphemeral *big.Int
	PublicEphemeral  *big.Int
}

func NewClient(c net.Conn) *RealmClient {
	rc := &RealmClient{
		c:       c,
		log:     log.With().Str("addr", c.RemoteAddr().String()).Logger(),
		account: nil,

		PrivateEphemeral: big.NewInt(0),
		PublicEphemeral:  big.NewInt(0),
	}

	go rc.listen()

	return rc
}

func (rc *RealmClient) HandleLogin(pkt *packets.ClientLoginChallenge) error {
	res := new(packets.ServerLoginChallenge)

	// TODO: is this safe?
	res.Status = packets.ChallengeStatusSuccess

	// Validate the packet.
	gameName := strings.TrimRight(pkt.GameName, "\x00")
	if gameName != static.SupportedGameName {
		res.Status = packets.ChallengeStatusFailed
	} else if pkt.Version != static.SupportedGameVersion || pkt.Build != static.SupportedGameBuild {
		res.Status = packets.ChallengeStatusFailVersionInvalid
	} else {
		rc.account = db.GetInstance().FindAccount(pkt.AccountName)

		if rc.account == nil {
			res.Status = packets.ChallengeStatusFailUnknownAccount
			rc.c.Close()
		}
	}

	if res.Status == packets.ChallengeStatusSuccess {
		b, B := srp.GenerateEphemeralPair(rc.account.Verifier())
		rc.PrivateEphemeral.Set(b)
		rc.PublicEphemeral.Set(B)

		res.B.Set(B)
		res.Salt.Set(rc.account.Salt())
		res.SaltCRC.SetInt64(0)
	}

	// Send out the packet
	return rc.Send(packets.AuthLoginChallenge, res.MarshalPacket())
}

func (rc *RealmClient) HandleProof(pkt *packets.ClientLoginProof) error {
	response := packets.ServerLoginProof{}

	K, M := srp.CalculateSessionKey(
		&pkt.A,
		rc.PublicEphemeral,
		rc.PrivateEphemeral,
		rc.account.Verifier(),
		rc.account.Salt(),
		rc.account.Name)

	if M.Cmp(&pkt.M) != 0 {
		response.Error = 4 // TODO(jeshua): make these constants
		rc.Send(packets.AuthLoginProof, response.MarshalPacket())
		rc.c.Close()

		return nil
	} else {
		response.Error = 0
		response.Proof.Set(srp.CalculateServerProof(&pkt.A, M, K))
		rc.log = rc.log.With().Str("account", rc.account.Name).Logger()

		rc.account.Session = K.Text(16)
	}

	return rc.Send(packets.AuthLoginProof, response.MarshalPacket())
}

func (rc *RealmClient) HandleRealmList() error {
	rc.log.Debug().Msg("handling realmlist request")

	srl := packets.ServerRealmlist{}
	srl.Realms = []packets.Realm{{
		Icon:          6,
		Lock:          0,
		Flags:         packets.RealmFlagNewPlayers | packets.RealmFlagRecommended,
		Name:          "Summit Pilot",
		Address:       "127.0.0.1:5000",
		Population:    0,
		NumCharacters: 5,
		Timezone:      1,
	}, {
		Icon:          1,
		Lock:          0,
		Flags:         packets.RealmFlagNewPlayers | packets.RealmFlagRecommended,
		Name:          "Summit Pilot PVP",
		Address:       "127.0.0.1:5000",
		Population:    0,
		NumCharacters: 2,
		Timezone:      2,
	}}

	return rc.Send(packets.RealmList, srl.MarshalPacket())
}

func (rc *RealmClient) Send(opcode packets.AuthCmd, payload []byte) error {
	size := len(payload)

	rc.log.Debug().
		Str("opcode", fmt.Sprintf("0x%04x", opcode)).
		Int("size", size).
		Hex("data", payload).
		Msg("sending packet to client")

	w := wow.NewPacketWriter()
	w.WriteByte(byte(opcode))
	// w.WriteByte(byte(size))
	w.Write(payload)

	return rc.Write(w.Bytes())
}

func (rc *RealmClient) Write(bb []byte) error {
	rc.outLock.Lock()
	defer rc.outLock.Unlock()

	w, err := rc.c.Write(bb)
	if err != nil {
		return err
	}

	if w != len(bb) {
		return errors.New("the written and sent bytes are not equal")
	}

	return nil
}

func (rc *RealmClient) listen() {
	defer rc.c.Close()
	rc.log.Info().Msgf("accepting messages from a new login connection")

	for {
		// Read packets infinitely :)
		pkt, err := rc.read(rc.c)
		if err != nil {
			log.Error().Err(err).Msg("error while reading from client")

			return
		}

		switch pkt.Command {
		case packets.AuthLoginChallenge:
			var clc packets.ClientLoginChallenge
			pkt.Unmarshal(&clc)
			rc.HandleLogin(&clc)
		case packets.AuthLoginProof:
			var clp packets.ClientLoginProof
			pkt.Unmarshal(&clp)
			rc.HandleProof(&clp)
		case packets.RealmList:
			var rlp packets.ClientRealmlist
			pkt.Unmarshal(&rlp)
			rc.HandleRealmList()
		default:
			rc.log.Fatal().Msgf("unhandled command: %T(0x%02x)", pkt.Command, pkt.Command)
		}
	}
}

// read reads the packet from the auth socket
func (rc *RealmClient) read(r io.Reader) (*packets.RData, error) {
	opCodeData := make([]byte, 1)
	n, err := r.Read(opCodeData)
	if err != nil {
		return nil, fmt.Errorf("erorr while reading opcode: %v", err)
	}

	if n != 1 {
		return nil, errors.New("short read when reading opcode data")
	}

	// In the auth server, the length is based on the packet type.
	opCode := packets.AuthCmd(opCodeData[0])
	length := 0

	switch opCode {
	case packets.AuthLoginChallenge:
		lenData, err := ReadBytes(r, 3)
		if err != nil {
			return nil, fmt.Errorf("error while reading header length: %v", err)
		}

		length = int(binary.LittleEndian.Uint16(lenData[1:]))
	case packets.AuthLoginProof:
		length = 74
	case packets.RealmList:
		length = 4
	default:
		rc.log.Error().
			Hex("packet", opCodeData).
			Msg("packet is not handled yet")

		return nil, err
	}

	ret := packets.RData{Command: opCode}
	bb, err := ReadBytes(r, length)
	if err != nil {
		return nil, err
	}

	ret.Data = bb

	return &ret, nil
}

// ReadBytes will read a specified number of bytes from a given buffer. If not all
// of the data is read (or there was an error), an error will be returned.
func ReadBytes(buffer io.Reader, length int) ([]byte, error) {
	data := make([]byte, length)

	if length > 0 {
		n, err := buffer.Read(data)
		if err != nil {
			return nil, fmt.Errorf("error while reading bytes: %v", err)
		}

		if n != length {
			fmt.Printf("%s\n", hex.Dump(data[:n]))
			return nil, fmt.Errorf("short read: wanted %v bytes, got %v", length, n)
		}
	}

	return data, nil
}
