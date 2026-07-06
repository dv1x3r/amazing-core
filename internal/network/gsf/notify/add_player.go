package notify

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// AddPlayer is the payload for GSFAddPlayerNotify.
type AddPlayer struct {
	// POID is the player OID.
	POID types.OID

	PlayerVillagerOID types.OID

	// Ver is the player-details cache version used by clients to refetch appearance data.
	Ver int64

	// LOID is the location OID.
	LOID types.OID

	// LCP is an unused client flag.
	LCP bool

	TimeOffset int64
	Pos        types.Position
	WPos       []types.Position
	QTH        types.QTH
	SecondQTH  types.QTH
	Weight     int32

	// Seq is the movement sequence byte; the current client sends 0.
	Seq byte

	Type        int32
	ActionState []types.PlayerActionNotify
}

func (n *AddPlayer) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&n.POID)
	writer.WriteObject(&n.PlayerVillagerOID)
	writer.WriteInt64(n.Ver)
	writer.WriteObject(&n.LOID)
	writer.WriteBool(n.LCP)
	writer.WriteInt64(n.TimeOffset)
	writer.WriteObject(&n.Pos)
	gsf.WriteSlice(writer, n.WPos, func(value types.Position) {
		writer.WriteObject(&value)
	})
	writer.WriteObject(&n.QTH)
	writer.WriteObject(&n.SecondQTH)
	writer.WriteInt32(n.Weight)
	writer.PutByte(n.Seq)
	writer.WriteInt32(n.Type)
	gsf.WriteSlice(writer, n.ActionState, func(value types.PlayerActionNotify) {
		writer.WriteObject(&value)
	})
}
