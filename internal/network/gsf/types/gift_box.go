package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// GiftBox is the asset container send as a reward.
type GiftBox struct {
	AssetContainer
	CreateDate gsf.UnixTime
}

func (gb *GiftBox) Serialize(writer gsf.ProtocolWriter) {
	gb.AssetContainer.Serialize(writer)
	writer.WriteUtcDate(gb.CreateDate)
}

func (gb *GiftBox) Deserialize(reader gsf.ProtocolReader) {
	gb.AssetContainer.Deserialize(reader)
	gb.CreateDate = reader.ReadUtcDate()
}
