package packets

import (
	"math/big"

	"github.com/paalgyula/summit/pkg/wow"
)

// ClientLoginProof encodes proof that the client has the correct information.
type ClientLoginProof struct {
	A             big.Int
	M             big.Int
	CRCHash       []byte // 20byte long
	NumberOfKeys  uint8
	SecurityFlags uint8
}

func (ClientLoginProof) OpCode() AuthCmd {
	return AuthLoginProof
}

func (pkt ClientLoginProof) MarshalPacket() []byte {
	w := wow.NewPacket(wow.OpCode(AuthLoginProof))

	w.WriteZeroPadded(wow.ReverseBytes(pkt.A.Bytes()), 32)
	w.WriteZeroPadded(wow.ReverseBytes(pkt.M.Bytes()), 20)
	// crc := make([]byte, 20)
	w.WriteZeroPadded(pkt.CRCHash, 20)
	w.WriteOne(0x00)

	return w.Bytes()
}

// Read will load a ClientLoginProof packet from a buffer.
// An error will be returned if at least one of the fields didn't load correctly.
func (pkt *ClientLoginProof) UnmarshalPacket(bb wow.PacketData) error {
	r := wow.NewPacketReader(bb)

	pkt.A.SetBytes(r.ReadReverseBytes(32))
	pkt.M.SetBytes(r.ReadReverseBytes(20))
	pkt.CRCHash, _ = r.ReadNBytes(20)

	r.ReadL(&pkt.NumberOfKeys)
	return r.ReadL(&pkt.SecurityFlags)
}

// ServerLoginProof is the server's response to a client's challenge. It contains
// some SRP information used for handshaking.
type ServerLoginProof struct {
	StatusCode uint8
	Proof      big.Int
}

// Bytes writes out the packet to an array of bytes.
func (pkt *ServerLoginProof) MarshalPacket() []byte {
	w := wow.NewPacket(wow.OpCode(AuthLoginProof))

	w.Write(pkt.StatusCode)

	if pkt.StatusCode == 0 {
		w.WriteZeroPadded(wow.ReverseBytes(pkt.Proof.Bytes()), 32)
		// buffer.Write([]byte("\x00\x00\x00\x00")) // unk1
	}

	return w.Bytes()
}
