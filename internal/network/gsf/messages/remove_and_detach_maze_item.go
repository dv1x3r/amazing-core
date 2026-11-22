package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// RemoveAndDetachMazeItemRequest takes an item on a mount of its parent back to the inventory.
type RemoveAndDetachMazeItemRequest struct {
	InventoryOID types.OID
	// Action is always 0 in MazeManager.RemoveMazeItem.
	Action int32
}

func (req *RemoveAndDetachMazeItemRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.InventoryOID)
	req.Action = reader.ReadInt32()
}

// RemoveAndDetachMazeItemResponse tells whether the item was taken back.
type RemoveAndDetachMazeItemResponse struct {
	IsUpdated bool
}

func (res *RemoveAndDetachMazeItemResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsUpdated)
}
