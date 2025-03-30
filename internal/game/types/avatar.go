package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type Avatar struct {
	AssetContainer
	Dimensions string
	Weight     float64
	Height     float64
	MaxOutfits int16
	Name       string
}

func (a *Avatar) Serialize(writer gsf.ProtocolWriter) {
	a.AssetContainer.Serialize(writer)
	writer.WriteString(a.Dimensions)
	writer.WriteFloat64(a.Weight)
	writer.WriteFloat64(a.Height)
	writer.WriteInt16(a.MaxOutfits)
	writer.WriteString(a.Name)
}

func (a *Avatar) Deserialize(reader gsf.ProtocolReader) {
	a.AssetContainer.Deserialize(reader)
	a.Dimensions = reader.ReadString()
	a.Weight = reader.ReadFloat64()
	a.Height = reader.ReadFloat64()
	a.MaxOutfits = reader.ReadInt16()
	a.Name = reader.ReadString()
}
