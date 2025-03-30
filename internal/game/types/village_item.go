package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type VillageItem struct {
	OID       OID
	PlayerID  OID
	X         int32
	Y         int32
	Z         int32
	Rotation  string
	Ordinal   int32
	ItemID    OID
	ItemType  int32
	VillageID OID
}

func (vi *VillageItem) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&vi.OID)
	writer.WriteObject(&vi.PlayerID)
	writer.WriteInt32(vi.X)
	writer.WriteInt32(vi.Y)
	writer.WriteInt32(vi.Z)
	writer.WriteString(vi.Rotation)
	writer.WriteInt32(vi.Ordinal)
	writer.WriteObject(&vi.ItemID)
	writer.WriteInt32(vi.ItemType)
	writer.WriteObject(&vi.VillageID)
}

func (vi *VillageItem) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&vi.OID)
	reader.ReadObject(&vi.PlayerID)
	vi.X = reader.ReadInt32()
	vi.Y = reader.ReadInt32()
	vi.Z = reader.ReadInt32()
	vi.Rotation = reader.ReadString()
	vi.Ordinal = reader.ReadInt32()
	reader.ReadObject(&vi.ItemID)
	vi.ItemType = reader.ReadInt32()
	reader.ReadObject(&vi.VillageID)
}
