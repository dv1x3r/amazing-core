package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// PlayerQuest is a quest instance for a player.
type PlayerQuest struct {
	OID                   OID
	CreateDate            gsf.UnixTime
	AcceptedDate          gsf.UnixTime
	StartedDate           gsf.UnixTime
	ExpiryDate            gsf.UnixTime
	CompletedDate         gsf.UnixTime
	PlayerOID             OID
	QuestStateOID         OID
	QuestStateName        string
	PlayerAvatarOID       OID
	ParentPlayerQuestOID  OID
	Quest                 Quest
	PlayerLevel           int32
	NPCRelationshipLevel  int32
	NPCRelationshipPoints int32
	PlayerMoney           int32
	PlayerXP              int32
	Unlocked              bool
	RuleProperty          RuleProperty
}

func (pq *PlayerQuest) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pq.OID)
	writer.WriteUtcDate(pq.CreateDate)
	writer.WriteUtcDate(pq.AcceptedDate)
	writer.WriteUtcDate(pq.StartedDate)
	writer.WriteUtcDate(pq.ExpiryDate)
	writer.WriteUtcDate(pq.CompletedDate)
	writer.WriteObject(&pq.PlayerOID)
	writer.WriteObject(&pq.QuestStateOID)
	writer.WriteString(pq.QuestStateName)
	writer.WriteObject(&pq.PlayerAvatarOID)
	writer.WriteObject(&pq.ParentPlayerQuestOID)
	writer.WriteObject(&pq.Quest)
	writer.WriteInt32(pq.PlayerLevel)
	writer.WriteInt32(pq.NPCRelationshipLevel)
	writer.WriteInt32(pq.NPCRelationshipPoints)
	writer.WriteInt32(pq.PlayerMoney)
	writer.WriteInt32(pq.PlayerXP)
	writer.WriteBool(pq.Unlocked)
	writer.WriteObject(&pq.RuleProperty)
}

func (pq *PlayerQuest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pq.OID)
	pq.CreateDate = reader.ReadUtcDate()
	pq.AcceptedDate = reader.ReadUtcDate()
	pq.StartedDate = reader.ReadUtcDate()
	pq.ExpiryDate = reader.ReadUtcDate()
	pq.CompletedDate = reader.ReadUtcDate()
	reader.ReadObject(&pq.PlayerOID)
	reader.ReadObject(&pq.QuestStateOID)
	pq.QuestStateName = reader.ReadString()
	reader.ReadObject(&pq.PlayerAvatarOID)
	reader.ReadObject(&pq.ParentPlayerQuestOID)
	reader.ReadObject(&pq.Quest)
	pq.PlayerLevel = reader.ReadInt32()
	pq.NPCRelationshipLevel = reader.ReadInt32()
	pq.NPCRelationshipPoints = reader.ReadInt32()
	pq.PlayerMoney = reader.ReadInt32()
	pq.PlayerXP = reader.ReadInt32()
	pq.Unlocked = reader.ReadBool()
	reader.ReadObject(&pq.RuleProperty)
}
