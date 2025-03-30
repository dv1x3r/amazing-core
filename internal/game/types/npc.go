package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type NPC struct {
	RuleContainer
	ZoneID             OID
	SpawnPoint         string
	StartQuestIDs      []OID
	RelationshipPoints int32
	Friend             bool
	Ordinal            int32
	PnrCreatedDate     time.Time
	PlayerItem         PlayerItem
}

func (npc *NPC) Serialize(writer gsf.ProtocolWriter) {
	npc.RuleContainer.Serialize(writer)
	writer.WriteObject(&npc.ZoneID)
	writer.WriteString(npc.SpawnPoint)
	gsf.WriteSlice(writer, npc.StartQuestIDs, func(value OID) {
		writer.WriteObject(&value)
	})
	writer.WriteInt32(npc.RelationshipPoints)
	writer.WriteBool(npc.Friend)
	writer.WriteInt32(npc.Ordinal)
	writer.WriteUtcDate(npc.PnrCreatedDate)
	writer.WriteObject(&npc.PlayerItem)
}

func (npc *NPC) Deserialize(reader gsf.ProtocolReader) {
	npc.RuleContainer.Deserialize(reader)
	reader.ReadObject(&npc.ZoneID)
	npc.SpawnPoint = reader.ReadString()
	npc.StartQuestIDs = gsf.ReadSlice(reader, func() OID {
		var value OID
		reader.ReadObject(&value)
		return value
	})
	npc.RelationshipPoints = reader.ReadInt32()
	npc.Friend = reader.ReadBool()
	npc.Ordinal = reader.ReadInt32()
	npc.PnrCreatedDate = reader.ReadUtcDate()
	reader.ReadObject(&npc.PlayerItem)
}
