package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetStoreItemsRequest requests item instances placed in a store.
type GetStoreItemsRequest struct {
	StoreOID            types.OID
	CategoryOID         types.OID
	Flags               int32
	CompletionEventName string
}

func (req *GetStoreItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.StoreOID)
	reader.ReadObject(&req.CategoryOID)
	req.Flags = reader.ReadInt32()
	req.CompletionEventName = reader.ReadString()
}

// GetStoreItemsResponse contains item instances placed in a store.
type GetStoreItemsResponse struct {
	StoreItems []types.StoreItem
}

func (res *GetStoreItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.StoreItems, func(value types.StoreItem) {
		writer.WriteObject(&value)
	})
}
