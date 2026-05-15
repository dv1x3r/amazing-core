package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// VillageItem is an item placement within a village.
type VillageItem struct {
	OID        OID
	PlayerOID  OID
	X          int32
	Y          int32
	Z          int32
	Rotation   string
	Ordinal    int32
	ItemOID    OID
	ItemType   int32
	VillageOID OID
}

func (vi *VillageItem) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&vi.OID)
	writer.WriteObject(&vi.PlayerOID)
	writer.WriteInt32(vi.X)
	writer.WriteInt32(vi.Y)
	writer.WriteInt32(vi.Z)
	writer.WriteString(vi.Rotation)
	writer.WriteInt32(vi.Ordinal)
	writer.WriteObject(&vi.ItemOID)
	writer.WriteInt32(vi.ItemType)
	writer.WriteObject(&vi.VillageOID)
}

func (vi *VillageItem) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&vi.OID)
	reader.ReadObject(&vi.PlayerOID)
	vi.X = reader.ReadInt32()
	vi.Y = reader.ReadInt32()
	vi.Z = reader.ReadInt32()
	vi.Rotation = reader.ReadString()
	vi.Ordinal = reader.ReadInt32()
	reader.ReadObject(&vi.ItemOID)
	vi.ItemType = reader.ReadInt32()
	reader.ReadObject(&vi.VillageOID)
}
