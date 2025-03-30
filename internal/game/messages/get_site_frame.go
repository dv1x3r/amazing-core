package messages

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types"
)

type GetSiteFrameRequest struct {
	TypeValue        int32
	LangLocalePairID types.OID
	TierID           types.OID
	BirthDate        time.Time
	RegistrationDate time.Time
	PreviewDate      time.Time
	IsPreviewEnabled bool
}

func (req *GetSiteFrameRequest) Deserialize(reader gsf.ProtocolReader) {
	req.TypeValue = reader.ReadInt32()
	reader.ReadObject(&req.LangLocalePairID)
	reader.ReadObject(&req.TierID)
	req.BirthDate = reader.ReadUtcDate()
	req.RegistrationDate = reader.ReadUtcDate()
	req.PreviewDate = reader.ReadUtcDate()
	req.IsPreviewEnabled = reader.ReadBool()
}

type GetSiteFrameResponse struct {
	SiteFrame        types.SiteFrame
	AssetDeliveryURL string
}

func (res *GetSiteFrameResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.SiteFrame)
	writer.WriteString(res.AssetDeliveryURL)
}
