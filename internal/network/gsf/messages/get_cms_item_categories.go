package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetCMSItemCategoriesRequest requests the CMS item category lookup.
type GetCMSItemCategoriesRequest struct {
}

func (req *GetCMSItemCategoriesRequest) Deserialize(reader gsf.ProtocolReader) {
}

// GetCMSItemCategoriesResponse contains CMS item category definitions.
type GetCMSItemCategoriesResponse struct {
	ItemCategories []types.ItemCategory
}

func (res *GetCMSItemCategoriesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.ItemCategories, func(value types.ItemCategory) {
		writer.WriteObject(&value)
	})
}
