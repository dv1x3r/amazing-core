package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// Quest is a quest definition instance.
type Quest struct {
	RuleContainer
	SingleSession   bool
	RetainHistory   bool
	EPoints         int32
	SpawnPoint      string
	Ordinal         int32
	QuestTypeOID    OID
	MinPlayers      int32
	MaxPlayers      int32
	DifficultyLevel gsf.Null[int32]
	QuestItems      []QuestItem
	ParentOID       OID
	ChildrenOIDs    []OID
	Children        []Quest
}

func (q *Quest) Serialize(writer gsf.ProtocolWriter) {
	q.RuleContainer.Serialize(writer)
	writer.WriteBool(q.SingleSession)
	writer.WriteBool(q.RetainHistory)
	writer.WriteInt32(q.EPoints)
	writer.WriteString(q.SpawnPoint)
	writer.WriteInt32(q.Ordinal)
	writer.WriteObject(&q.QuestTypeOID)
	writer.WriteInt32(q.MinPlayers)
	writer.WriteInt32(q.MaxPlayers)
	gsf.WriteNullable(writer, q.DifficultyLevel, writer.WriteInt32)
	gsf.WriteSlice(writer, q.QuestItems, func(value QuestItem) {
		writer.WriteObject(&value)
	})
	writer.WriteObject(&q.ParentOID)
	gsf.WriteSlice(writer, q.ChildrenOIDs, func(value OID) {
		writer.WriteObject(&value)
	})
	gsf.WriteSlice(writer, q.Children, func(value Quest) {
		writer.WriteObject(&value)
	})
}

func (q *Quest) Deserialize(reader gsf.ProtocolReader) {
	q.RuleContainer.Deserialize(reader)
	q.SingleSession = reader.ReadBool()
	q.RetainHistory = reader.ReadBool()
	q.EPoints = reader.ReadInt32()
	q.SpawnPoint = reader.ReadString()
	q.Ordinal = reader.ReadInt32()
	reader.ReadObject(&q.QuestTypeOID)
	q.MinPlayers = reader.ReadInt32()
	q.MaxPlayers = reader.ReadInt32()
	q.DifficultyLevel = gsf.ReadNullable(reader, reader.ReadInt32)
	q.QuestItems = gsf.ReadSlice(reader, func() QuestItem {
		var value QuestItem
		reader.ReadObject(&value)
		return value
	})
	reader.ReadObject(&q.ParentOID)
	q.ChildrenOIDs = gsf.ReadSlice(reader, func() OID {
		var value OID
		reader.ReadObject(&value)
		return value
	})
	q.Children = gsf.ReadSlice(reader, func() Quest {
		var value Quest
		reader.ReadObject(&value)
		return value
	})
}
