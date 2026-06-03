package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// AddOutfitItemsRequest requests the outfit item assign.
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

// AddOutfitItemsResponse contains the assign status.
type AddOutfitItemsResponse struct {
	IsUpdated bool
}

func (res *AddOutfitItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsUpdated)
}
