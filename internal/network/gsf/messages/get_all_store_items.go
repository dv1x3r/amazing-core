package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetAllStoreItemsRequest requests the products of a store.
type GetAllStoreItemsRequest struct {
	StoreOID            types.OID
	CategoryOID         types.OID
	Flags               int32
	CompletionEventName string
}

func (req *GetAllStoreItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.StoreOID)
	reader.ReadObject(&req.CategoryOID)
	req.Flags = reader.ReadInt32()
	req.CompletionEventName = reader.ReadString()
}

// GetAllStoreItemsResponse contains the products of a store.
type GetAllStoreItemsResponse struct {
	StoreItems []types.StoreItem
}

func (res *GetAllStoreItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.StoreItems, func(value types.StoreItem) {
		writer.WriteObject(&value)
	})
}
