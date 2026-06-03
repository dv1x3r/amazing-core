package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetInventoryObjectsRequest requests player-owned items for the inventory grid.
type GetInventoryObjectsRequest struct {
	ContainerOID     types.OID
	ItemCategoryOIDs []types.OID
	PlayerItemOIDs   []types.OID
}

func (req *GetInventoryObjectsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.ContainerOID)
	req.ItemCategoryOIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
	req.PlayerItemOIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
}

// GetInventoryObjectsResponse contains player-owned items for the inventory grid.
type GetInventoryObjectsResponse struct {
	PlayerItems []types.PlayerItem
}

func (res *GetInventoryObjectsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.PlayerItems, func(value types.PlayerItem) {
		writer.WriteObject(&value)
	})
}
