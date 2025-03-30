package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type SiteFrame struct {
	AssetContainer
	TypeValue int32
}

func (sf *SiteFrame) Serialize(writer gsf.ProtocolWriter) {
	sf.AssetContainer.Serialize(writer)
	writer.WriteInt32(sf.TypeValue)
}

func (sf *SiteFrame) Deserialize(reader gsf.ProtocolReader) {
	sf.AssetContainer.Deserialize(reader)
	sf.TypeValue = reader.ReadInt32()
}
