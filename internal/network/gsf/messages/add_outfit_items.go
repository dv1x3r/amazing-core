package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// AddOutfitItemsRequest contains the items to put into the outfit slots.
type AddOutfitItemsRequest struct {
	PlayerAvatarOutfitOID types.OID
	InventoryOIDs         []types.OID
	SlotOIDs              []types.OID
}

func (req *AddOutfitItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarOutfitOID)
	req.InventoryOIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
	req.SlotOIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
}

// AddOutfitItemsResponse contains whether the outfit was updated.
type AddOutfitItemsResponse struct {
	IsUpdated bool
}

func (res *AddOutfitItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsUpdated)
}
