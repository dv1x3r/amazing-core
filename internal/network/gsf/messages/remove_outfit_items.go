package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// RemoveOutfitItemsRequest contains the items to take out of the outfit.
type RemoveOutfitItemsRequest struct {
	PlayerAvatarOutfitOID types.OID
	InventoryOIDs         []types.OID
}

func (req *RemoveOutfitItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarOutfitOID)
	req.InventoryOIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
}

// RemoveOutfitItemsResponse contains whether the outfit was updated.
type RemoveOutfitItemsResponse struct {
	IsUpdated bool
}

func (res *RemoveOutfitItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsUpdated)
}
