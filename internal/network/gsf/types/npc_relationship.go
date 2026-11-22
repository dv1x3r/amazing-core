package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// NPCRelationship is the relationship of the player with an NPC, NPCManager.listRelationships keeps it by NPCOID.
// The level picks the name of the NPC Relationship_Text, the meter of NPCRadialMenu is (Points - MinPoints) / (MaxPoints - MinPoints).
type NPCRelationship struct {
	OID             OID
	CreateDate      gsf.UnixTime
	PlayerOID       OID
	NPCOID          OID
	Level           int16
	Points          int32
	LastInteraction gsf.UnixTime
	UnlockDate      gsf.UnixTime
	Friend          bool
	Ordinal         int32
	MinPoints       int32
	MaxPoints       int32
	IsNew           bool
}

func (r *NPCRelationship) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&r.OID)
	writer.WriteUtcDate(r.CreateDate)
	writer.WriteObject(&r.PlayerOID)
	writer.WriteObject(&r.NPCOID)
	writer.WriteInt16(r.Level)
	writer.WriteInt32(r.Points)
	writer.WriteUtcDate(r.LastInteraction)
	writer.WriteUtcDate(r.UnlockDate)
	writer.WriteBool(r.Friend)
	writer.WriteInt32(r.Ordinal)
	writer.WriteInt32(r.MinPoints)
	writer.WriteInt32(r.MaxPoints)
	writer.WriteBool(r.IsNew)
}

func (r *NPCRelationship) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&r.OID)
	r.CreateDate = reader.ReadUtcDate()
	reader.ReadObject(&r.PlayerOID)
	reader.ReadObject(&r.NPCOID)
	r.Level = reader.ReadInt16()
	r.Points = reader.ReadInt32()
	r.LastInteraction = reader.ReadUtcDate()
	r.UnlockDate = reader.ReadUtcDate()
	r.Friend = reader.ReadBool()
	r.Ordinal = reader.ReadInt32()
	r.MinPoints = reader.ReadInt32()
	r.MaxPoints = reader.ReadInt32()
	r.IsNew = reader.ReadBool()
}

// NPCRelationshipLevel is the points of the relationship level of an NPC.
type NPCRelationshipLevel struct {
	OID               OID
	Level             int32
	RelationshipValue int32
	NPCOID            OID
}

func (l *NPCRelationshipLevel) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&l.OID)
	writer.WriteInt32(l.Level)
	writer.WriteInt32(l.RelationshipValue)
	writer.WriteObject(&l.NPCOID)
}

func (l *NPCRelationshipLevel) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&l.OID)
	l.Level = reader.ReadInt32()
	l.RelationshipValue = reader.ReadInt32()
	reader.ReadObject(&l.NPCOID)
}
