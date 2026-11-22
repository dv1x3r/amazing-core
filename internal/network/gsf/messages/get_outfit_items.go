package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetOutfitItemsRequest contains the outfit of the items.
type GetOutfitItemsRequest struct {
	PlayerAvatarOutfitOID types.OID
	PlayerOID             types.OID
}

func (req *GetOutfitItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarOutfitOID)
	reader.ReadObject(&req.PlayerOID)
}

// GetOutfitItemsResponse contains the items of the outfit.
type GetOutfitItemsResponse struct {
	OutfitItems []types.PlayerItem
}

func (res *GetOutfitItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.OutfitItems, func(value types.PlayerItem) {
		writer.WriteObject(&value)
	})
}
