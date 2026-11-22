package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetStoreItemsForItemsRequest contains the store and the item to sell.
type GetStoreItemsForItemsRequest struct {
	StoreOID types.OID
	ItemOID  types.OID
}

func (req *GetStoreItemsForItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.StoreOID)
	reader.ReadObject(&req.ItemOID)
}

// GetStoreItemsForItemsResponse contains the buyback product of the item, nil when it can't be sold.
type GetStoreItemsForItemsResponse struct {
	StoreItem    *types.StoreItem
	SellResponse int32
}

func (res *GetStoreItemsForItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(res.StoreItem)
	writer.WriteInt32(res.SellResponse)
}
