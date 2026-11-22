package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetStoreThemesRequest contains the store of the categories.
type GetStoreThemesRequest struct {
	StoreOID types.OID
}

func (req *GetStoreThemesRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.StoreOID)
}

// GetStoreThemesResponse contains the store categories and subcategories.
type GetStoreThemesResponse struct {
	StoreThemes []types.StoreTheme
}

func (res *GetStoreThemesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.StoreThemes, func(value types.StoreTheme) {
		writer.WriteObject(&value)
	})
}
