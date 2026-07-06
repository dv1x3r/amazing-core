package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// PlayerActionNotify describes a player action notification on the SYNC server.
type PlayerActionNotify struct {
	// POID is the player OID.
	POID OID

	ActionOID  OID
	TargetOID  OID
	SourceOID  OID
	ActionType byte
	Parameters string
	Pos        Position
	QTH        QTH
	IsState    bool
	ClearState bool
	StateIndex int32
}

func (pan *PlayerActionNotify) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pan.POID)
	writer.WriteObject(&pan.ActionOID)
	writer.WriteObject(&pan.TargetOID)
	writer.WriteObject(&pan.SourceOID)
	writer.PutByte(pan.ActionType)
	writer.WriteString(pan.Parameters)
	writer.WriteObject(&pan.Pos)
	writer.WriteObject(&pan.QTH)
	writer.WriteBool(pan.IsState)
	writer.WriteBool(pan.ClearState)
	writer.WriteInt32(pan.StateIndex)
}

func (pan *PlayerActionNotify) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pan.POID)
	reader.ReadObject(&pan.ActionOID)
	reader.ReadObject(&pan.TargetOID)
	reader.ReadObject(&pan.SourceOID)
	pan.ActionType = reader.GetByte()
	pan.Parameters = reader.ReadString()
	reader.ReadObject(&pan.Pos)
	reader.ReadObject(&pan.QTH)
	pan.IsState = reader.ReadBool()
	pan.ClearState = reader.ReadBool()
	pan.StateIndex = reader.ReadInt32()
}
