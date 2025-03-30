package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type InventoryPosition struct {
	ObjectPosition ObjectPosition
	Rotation       string
}

func (ip *InventoryPosition) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&ip.ObjectPosition)
	writer.WriteString(ip.Rotation)
}

func (ip *InventoryPosition) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&ip.ObjectPosition)
	ip.Rotation = reader.ReadString()
}
