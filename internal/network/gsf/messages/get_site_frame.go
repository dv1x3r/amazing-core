package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetSiteFrameRequest requests the core site frame asset container.
type GetSiteFrameRequest struct {
	// The client always sends 1.
	TypeValue int32

	// The client always sends 293578400718237473.
	LangLocalePairID types.OID

	// The following fields are serialized but not used.
	TierID           types.OID
	BirthDate        gsf.UnixTime
	RegistrationDate gsf.UnixTime
	PreviewDate      gsf.UnixTime
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

// GetSiteFrameResponse contains the core site frame and asset delivery URL.
type GetSiteFrameResponse struct {
	// Contains the core assets for the current platform.
	SiteFrame types.SiteFrame

	// Base URL used for asset downloads.
	AssetDeliveryURL string
}

func (res *GetSiteFrameResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.SiteFrame)
	writer.WriteString(res.AssetDeliveryURL)
}
