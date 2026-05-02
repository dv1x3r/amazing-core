package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// AssetPackage is an asset container that is referenced by another asset container.
// Contains an additional tag used to conditionally load assets depending on which maze the player is entering.
type AssetPackage struct {
	AssetContainer
	PTag        string
	CreatedDate gsf.UnixTime
}

func (ap *AssetPackage) Serialize(writer gsf.ProtocolWriter) {
	ap.AssetContainer.Serialize(writer)
	writer.WriteString(ap.PTag)
	writer.WriteUtcDate(ap.CreatedDate)
}

func (ap *AssetPackage) Deserialize(reader gsf.ProtocolReader) {
	ap.AssetContainer.Deserialize(reader)
	ap.PTag = reader.ReadString()
	ap.CreatedDate = reader.ReadUtcDate()
}
