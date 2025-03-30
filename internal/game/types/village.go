package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type Village struct {
	OID                OID
	NextVillageID      OID
	PrevVillageID      OID
	VillageShardID     OID
	VillageTemplateID  OID
	VillageFlagID      OID
	MayorPlayerID      OID
	VillageTheme       AssetContainer
	VillagePlots       []VillagePlot
	VillageRolePlayers []VillageRolePlayer
	VillageItems       []VillageItem
	IsOpen             bool
	IsPublic           bool
	VillageName        string
	ShardNo            gsf.Null[int32]
	Rating             int32
}

func (v *Village) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&v.OID)
	writer.WriteObject(&v.NextVillageID)
	writer.WriteObject(&v.PrevVillageID)
	writer.WriteObject(&v.VillageShardID)
	writer.WriteObject(&v.VillageTemplateID)
	writer.WriteObject(&v.VillageFlagID)
	writer.WriteObject(&v.MayorPlayerID)
	writer.WriteObject(&v.VillageTheme)
	gsf.WriteSlice(writer, v.VillagePlots, func(value VillagePlot) {
		writer.WriteObject(&value)
	})
	gsf.WriteSlice(writer, v.VillageRolePlayers, func(value VillageRolePlayer) {
		writer.WriteObject(&value)
	})
	gsf.WriteSlice(writer, v.VillageItems, func(value VillageItem) {
		writer.WriteObject(&value)
	})
	writer.WriteBool(v.IsOpen)
	writer.WriteBool(v.IsPublic)
	writer.WriteString(v.VillageName)
	gsf.WriteNullable(writer, v.ShardNo, writer.WriteInt32)
	writer.WriteInt32(v.Rating)
}

func (v *Village) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&v.OID)
	reader.ReadObject(&v.NextVillageID)
	reader.ReadObject(&v.PrevVillageID)
	reader.ReadObject(&v.VillageShardID)
	reader.ReadObject(&v.VillageTemplateID)
	reader.ReadObject(&v.VillageFlagID)
	reader.ReadObject(&v.MayorPlayerID)
	reader.ReadObject(&v.VillageTheme)
	v.VillagePlots = gsf.ReadSlice(reader, func() VillagePlot {
		var value VillagePlot
		reader.ReadObject(&value)
		return value
	})
	v.VillageRolePlayers = gsf.ReadSlice(reader, func() VillageRolePlayer {
		var value VillageRolePlayer
		reader.ReadObject(&value)
		return value
	})
	v.VillageItems = gsf.ReadSlice(reader, func() VillageItem {
		var value VillageItem
		reader.ReadObject(&value)
		return value
	})
	v.IsOpen = reader.ReadBool()
	v.IsPublic = reader.ReadBool()
	v.VillageName = reader.ReadString()
	v.ShardNo = gsf.ReadNullable(reader, reader.ReadInt32)
	v.Rating = reader.ReadInt32()
}
