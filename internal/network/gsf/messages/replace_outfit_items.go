package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// ReplaceOutfitItemsRequest requests the outfit item swap.
type ReplaceOutfitItemsRequest struct {
	PlayerAvatarOutfitOID types.OID
	OldInventoryOIDs      []types.OID
	NewInventoryOIDs      []types.OID
}

func (req *ReplaceOutfitItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarOutfitOID)
	req.OldInventoryOIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
	req.NewInventoryOIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
}

// ReplaceOutfitItemsResponse contains the swap status.
type ReplaceOutfitItemsResponse struct {
	IsUpdated bool
}

func (res *ReplaceOutfitItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsUpdated)
}
