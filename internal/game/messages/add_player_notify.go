package messages

import (
	"github.com/dv1x3r/amazing-core/internal/game/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
)

type AddPlayerNotify struct {
	PID              types.OID
	PlayerVillagerID types.OID
	Ver              int64
	LID              types.OID
	LCP              bool
	TimeOffset       int64
	Pos              types.Position
	WPos             []types.Position
	Qth              types.Qth
	SecondQth        types.Qth
	Weight           int32
	Seq              byte
	Type             int32
	ActionState      []types.PlayerActionNotify
}

func (apn *AddPlayerNotify) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&apn.PID)
	writer.WriteObject(&apn.PlayerVillagerID)
	writer.WriteInt64(apn.Ver)
	writer.WriteObject(&apn.LID)
	writer.WriteBool(apn.LCP)
	writer.WriteInt64(apn.TimeOffset)
	writer.WriteObject(&apn.Pos)
	gsf.WriteSlice(writer, apn.WPos, func(value types.Position) {
		writer.WriteObject(&value)
	})
	writer.WriteObject(&apn.Qth)
	writer.WriteObject(&apn.SecondQth)
	writer.WriteInt32(apn.Weight)
	writer.PutByte(apn.Seq)
	writer.WriteInt32(apn.Type)
	gsf.WriteSlice(writer, apn.ActionState, func(value types.PlayerActionNotify) {
		writer.WriteObject(&value)
	})
}

func (apn *AddPlayerNotify) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&apn.PID)
	reader.ReadObject(&apn.PlayerVillagerID)
	apn.Ver = reader.ReadInt64()
	reader.ReadObject(&apn.LID)
	apn.LCP = reader.ReadBool()
	apn.TimeOffset = reader.ReadInt64()
	reader.ReadObject(&apn.Pos)
	apn.WPos = gsf.ReadSlice(reader, func() types.Position {
		var value types.Position
		reader.ReadObject(&value)
		return value
	})
	reader.ReadObject(&apn.Qth)
	reader.ReadObject(&apn.SecondQth)
	apn.Weight = reader.ReadInt32()
	apn.Seq = reader.GetByte()
	apn.Type = reader.ReadInt32()
	apn.ActionState = gsf.ReadSlice(reader, func() types.PlayerActionNotify {
		var value types.PlayerActionNotify
		reader.ReadObject(&value)
		return value
	})
}
