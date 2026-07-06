package notify

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// PlayerMove is the payload for GSFPlayerMoveNotify.
type PlayerMove struct {
	// POID is the player OID.
	POID types.OID

	PlayerVillagerOID types.OID

	// LCP is an unused client flag.
	LCP bool

	Pos       []types.Position
	QTH       types.QTH
	SecondQTH types.QTH

	// Seq is the movement sequence byte; the current client sends 0.
	Seq byte
}

func (n *PlayerMove) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&n.POID)
	writer.WriteObject(&n.PlayerVillagerOID)
	writer.WriteBool(n.LCP)
	gsf.WriteSlice(writer, n.Pos, func(value types.Position) {
		writer.WriteObject(&value)
	})
	writer.WriteObject(&n.QTH)
	writer.WriteObject(&n.SecondQTH)
	writer.PutByte(n.Seq)
}
