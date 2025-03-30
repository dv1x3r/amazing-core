package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type AssetPackage struct {
	AssetContainer
	PTag        string
	CreatedDate time.Time
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
