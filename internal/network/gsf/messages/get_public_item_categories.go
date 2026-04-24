package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

type GetPublicItemCategoriesRequest struct {
	LangLocalePairID types.OID
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

type GetPublicItemCategoriesResponse struct {
	ItemCategories []types.ItemCategory
}

func (res *GetPublicItemCategoriesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.ItemCategories, func(value types.ItemCategory) {
		writer.WriteObject(&value)
	})
}
