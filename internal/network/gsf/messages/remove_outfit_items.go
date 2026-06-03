package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// RemoveOutfitItemsRequest requests the outfit item removal.
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

// RemoveOutfitItemsResponse contains the removal status.
type RemoveOutfitItemsResponse struct {
	IsUpdated bool
}

func (res *RemoveOutfitItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsUpdated)
}
