package messages

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types"
)

type GetPublicItemCategoriesRequest struct {
	LangLocalePairID types.OID
	TierID           types.OID
	BirthDate        time.Time
	RegistrationDate time.Time
	PreviewDate      time.Time
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
