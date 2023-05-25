package packets

import (
	"math/big"
	"strings"

	"github.com/paalgyula/summit/pkg/blizzard/auth/srp"
	"github.com/paalgyula/summit/pkg/wow"
)

// ClientLoginChallenge received login challenge packet
type ClientLoginChallenge struct {
	GameName        string
	Version         [3]byte
	Build           uint16
	Platform        string
	OS              string
	Locale          string
	WorldRegionBias uint32
	IP              [4]uint8
	AccountName     string
}

func NewClientLoginChallenge(accName string) *ClientLoginChallenge {
	return &ClientLoginChallenge{
		GameName:        "\x00WoW",
		Version:         [3]byte{3, 3, 5},
		Build:           12340,
		Platform:        "\x0068x",
		OS:              "\x00niW",
		Locale:          "SUne",
		WorldRegionBias: 0,
		IP:              [4]uint8{89, 51, 25, 12},
		AccountName:     strings.ToUpper(accName),
	}
}

func (p *ClientLoginChallenge) UnmarshalPacket(bb wow.PacketData) error {
	r := bb.Reader()
	p.GameName = r.ReadStringFixed(4)
	r.ReadL(&p.Version)
	r.ReadL(&p.Build)
	p.Platform = r.ReadStringFixed(4)
	p.OS = r.ReadStringFixed(4)
	p.Locale = r.ReadStringFixed(4)
	r.ReadL(&p.WorldRegionBias)
	r.ReadL(&p.IP)

	var len uint8
	r.ReadB(&len)

	p.AccountName = r.ReadStringFixed(int(len))

	return nil
}

func (p *ClientLoginChallenge) MarshalPacket() []byte {
	w := wow.NewPacketWriter()
	w.WriteStringFixed(p.GameName, 4)
	w.Write(p.Version[:])
	w.WriteL(p.Build)
	w.WriteStringFixed(p.Platform, 4)
	w.WriteStringFixed(p.OS, 4)
	w.WriteStringFixed(p.Locale, 4)
	w.WriteL(p.WorldRegionBias)
	w.WriteL(p.IP)

	w.WriteL(uint8(len(p.AccountName)))
	w.WriteStringFixed(p.AccountName, len(p.AccountName))

	return w.Bytes()
}

type ChallengeStatus uint8

const (
	ChallengeStatusSuccess ChallengeStatus = iota
	// Account not found
	ChallengeStatusFailed
	// Account has been banned
	ChallengeStatusFailBanned
	// This <game> account has been closed and is no longer available for use. Please
	// go to <site>/banned.html for further information.
	ChallengeStatusFailUnknownAccount
	// The information you have entered is not valid. Please check the spelling
	// of the account name and password. If you need help in retrieving a lost or
	// stolen password, see <site> for more information
	ChallengeStatusFailUnknown0
	// The information you have entered is not valid. Please check the spelling
	// of the account name and password. If you need help in retrieving a lost
	// or stolen password, see <site> for more information
	ChallengeStatusFailIncorrectPassword
	// This account is already logged into <game>. Please check the spelling and try again.
	ChallengeStatusFailAlreadyOnline
	// You have used up your prepaid time for this account. Please purchase more to continue playing
	ChallengeStatusFailNoTime
	// Could not log in to <game> at this time. Please try again later.
	ChallengeStatusFailDbBusy
	// Unable to validate game version. This may be caused by file corruption or
	// interference of another program. Please visit <site> for more information
	// and possible solutions to this issue.
	ChallengeStatusFailVersionInvalid
	// Downloading
	ChallengeStatusFailVersionUpdate
	// Unable to connect
	ChallengeStatusFailInvalidServer
	// This <game> account has been temporarily suspended. Please go to <site>/banned.html for further information
	ChallengeStatusFailSuspended
	// Unable to connect
	ChallengeStatusFailFailNoaccess
	// Connected.
	ChallengeStatusSuccessSurvey
	// Access to this account has been blocked by parental controls. Your settings may be changed in your account preferences at <site>
	ChallengeStatusFailParentcontrol
	// You have applied a lock to your account. You can change your locked status by calling your account lock phone number.
	ChallengeStatusFailLockedEnforced
	// Your trial subscription has expired. Please visit <site> to upgrade your account.
	ChallengeStatusFailTrialEnded
)

// ServerLoginChallenge is the server's response to a client's challenge. It contains
// some SRP information used for handshaking.
type ServerLoginChallenge struct {
	Status  ChallengeStatus
	B       big.Int
	Salt    big.Int
	SaltCRC big.Int
}

func (pkt *ServerLoginChallenge) UnmarshalPacket(data wow.PacketData) {

}

// Bytes writes out the packet to an array of bytes.
func (pkt *ServerLoginChallenge) MarshalPacket() []byte {
	w := wow.NewPacketWriter()

	w.WriteByte(0) // unk1
	w.WriteByte(uint8(pkt.Status))

	if pkt.Status == ChallengeStatusSuccess {
		w.Write(PadBigIntBytes(wow.ReverseBytes(pkt.B.Bytes()), 32))
		w.WriteByte(1)
		w.WriteByte(srp.G)
		w.WriteByte(32)
		w.WriteReverse(srp.N().Bytes())
		w.Write(PadBigIntBytes(wow.ReverseBytes(pkt.Salt.Bytes()), 32))
		w.Write(PadBigIntBytes(wow.ReverseBytes(pkt.SaltCRC.Bytes()), 16))
		w.WriteByte(0) // unk2
	}

	return w.Bytes()
}
