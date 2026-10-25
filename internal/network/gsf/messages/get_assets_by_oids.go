package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetAssetsByOIDsRequest requests asset containers by OIDs.
type GetAssetsByOIDsRequest struct {
	OIDs []types.OID
}

func (req *GetAssetsByOIDsRequest) Deserialize(reader gsf.ProtocolReader) {
	req.OIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
}

// GetAssetsByOIDsResponse contains the asset containers.
type GetAssetsByOIDsResponse struct {
	Assets []types.AssetContainer
}

func (res *GetAssetsByOIDsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Assets, func(value types.AssetContainer) {
		writer.WriteObject(&value)
	})
}
