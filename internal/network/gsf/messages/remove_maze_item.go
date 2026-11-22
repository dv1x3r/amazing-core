package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// RemoveMazeItemRequest takes a placed item back to the inventory.
type RemoveMazeItemRequest struct {
	InventoryOID types.OID
}

func (req *RemoveMazeItemRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.InventoryOID)
}

// RemoveMazeItemResponse tells whether the item was taken back.
type RemoveMazeItemResponse struct {
	IsUpdated bool
}

func (res *RemoveMazeItemResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsUpdated)
}
