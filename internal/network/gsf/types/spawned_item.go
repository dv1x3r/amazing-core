package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// SpawnedItem is an item instance spawned at the specific position.
type SpawnedItem struct {
	InstanceUUID      string
	ItemOID           OID
	SpawnPoint        string
	Position          ObjectPosition
	CollectedSequence int32
}

func (si *SpawnedItem) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteString(si.InstanceUUID)
	writer.WriteObject(&si.ItemOID)
	writer.WriteString(si.SpawnPoint)
	writer.WriteObject(&si.Position)
	writer.WriteInt32(si.CollectedSequence)
}

func (si *SpawnedItem) Deserialize(reader gsf.ProtocolReader) {
	si.InstanceUUID = reader.ReadString()
	reader.ReadObject(&si.ItemOID)
	si.SpawnPoint = reader.ReadString()
	reader.ReadObject(&si.Position)
	si.CollectedSequence = reader.ReadInt32()
}
