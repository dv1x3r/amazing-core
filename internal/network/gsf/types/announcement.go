package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// Announcement is a client-visible announcement with asset and ordering data.
type Announcement struct {
	AssetContainer
	CreateTS gsf.UnixTime
	Ordinal  int32
}

func (a *Announcement) Serialize(writer gsf.ProtocolWriter) {
	a.AssetContainer.Serialize(writer)
	writer.WriteUtcDate(a.CreateTS)
	writer.WriteInt32(a.Ordinal)
}

func (a *Announcement) Deserialize(reader gsf.ProtocolReader) {
	a.AssetContainer.Deserialize(reader)
	a.CreateTS = reader.ReadUtcDate()
	a.Ordinal = reader.ReadInt32()
}
