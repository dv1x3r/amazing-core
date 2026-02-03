package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

type PlayerActionNotify struct {
	PID        OID
	ActionID   OID
	TargetID   OID
	SourceID   OID
	ActionType byte
	Parameters string
	Pos        Position
	Qth        Qth
	IsState    bool
	ClearState bool
	StateIndex int32
}

func (pan *PlayerActionNotify) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pan.PID)
	writer.WriteObject(&pan.ActionID)
	writer.WriteObject(&pan.TargetID)
	writer.WriteObject(&pan.SourceID)
	writer.PutByte(pan.ActionType)
	writer.WriteString(pan.Parameters)
	writer.WriteObject(&pan.Pos)
	writer.WriteObject(&pan.Qth)
	writer.WriteBool(pan.IsState)
	writer.WriteBool(pan.ClearState)
	writer.WriteInt32(pan.StateIndex)
}

func (pan *PlayerActionNotify) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pan.PID)
	reader.ReadObject(&pan.ActionID)
	reader.ReadObject(&pan.TargetID)
	reader.ReadObject(&pan.SourceID)
	pan.ActionType = reader.GetByte()
	pan.Parameters = reader.ReadString()
	reader.ReadObject(&pan.Pos)
	reader.ReadObject(&pan.Qth)
	pan.IsState = reader.ReadBool()
	pan.ClearState = reader.ReadBool()
	pan.StateIndex = reader.ReadInt32()
}
