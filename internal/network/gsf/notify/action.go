package notify

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// Action is the payload for GSFActionNotify.
// The client sends it from SyncManager.SyncActionNotify while in a zone: jumps, animations, fishing, held items, quest states.
type Action struct {
	ActionOID *types.OID
	TargetOID *types.OID
	SourceOID *types.OID

	// ActionType is SyncManager.ActionNotifyType.
	ActionType byte

	Parameters string
	Pos        *types.Position
	QTH        *types.QTH
	IsState    bool
	ClearState bool
	StateIndex int32

	// UID is the player OID from the client uid field.
	UID *types.OID

	// VID is the village/context OID from the client vid field.
	VID *types.OID
}

func (n *Action) Deserialize(reader gsf.ProtocolReader) {
	n.ActionOID = &types.OID{}
	if !reader.ReadObject(n.ActionOID) {
		n.ActionOID = nil
	}
	n.TargetOID = &types.OID{}
	if !reader.ReadObject(n.TargetOID) {
		n.TargetOID = nil
	}
	n.SourceOID = &types.OID{}
	if !reader.ReadObject(n.SourceOID) {
		n.SourceOID = nil
	}
	n.ActionType = reader.GetByte()
	n.Parameters = reader.ReadString()
	n.Pos = &types.Position{}
	if !reader.ReadObject(n.Pos) {
		n.Pos = nil
	}
	n.QTH = &types.QTH{}
	if !reader.ReadObject(n.QTH) {
		n.QTH = nil
	}
	n.IsState = reader.ReadBool()
	n.ClearState = reader.ReadBool()
	n.StateIndex = reader.ReadInt32()
	n.UID = &types.OID{}
	if !reader.ReadObject(n.UID) {
		n.UID = nil
	}
	n.VID = &types.OID{}
	if !reader.ReadObject(n.VID) {
		n.VID = nil
	}
}
