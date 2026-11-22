package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetStoreThemesRequest requests the categories and subcategories of a store.
type GetStoreThemesRequest struct {
	StoreOID types.OID
}

func (req *GetStoreThemesRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.StoreOID)
}

// GetStoreThemesResponse contains the categories and subcategories, the client reads their names from the Text assets.
type GetStoreThemesResponse struct {
	StoreThemes []types.StoreTheme
}

func (res *GetStoreThemesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.StoreThemes, func(value types.StoreTheme) {
		writer.WriteObject(&value)
	})
}
