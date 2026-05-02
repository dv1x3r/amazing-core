package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetPublicItemCategoriesRequest requests public item category definitions.
type GetPublicItemCategoriesRequest struct {
	// The client always sends 293578400718237473.
	LangLocalePairID types.OID

	// The following fields are serialized but not used.
	TierID           types.OID
	BirthDate        gsf.UnixTime
	RegistrationDate gsf.UnixTime
	PreviewDate      gsf.UnixTime
	IsPreviewEnabled bool
}

func (req *GetPublicItemCategoriesRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.LangLocalePairID)
	reader.ReadObject(&req.TierID)
	req.BirthDate = reader.ReadUtcDate()
	req.RegistrationDate = reader.ReadUtcDate()
	req.PreviewDate = reader.ReadUtcDate()
	req.IsPreviewEnabled = reader.ReadBool()
}

// GetPublicItemCategoriesResponse contains public item category definitions.
type GetPublicItemCategoriesResponse struct {
	ItemCategories []types.ItemCategory
}

func (res *GetPublicItemCategoriesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.ItemCategories, func(value types.ItemCategory) {
		writer.WriteObject(&value)
	})
}
