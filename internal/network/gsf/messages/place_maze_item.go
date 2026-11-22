package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// PlaceMazeItemRequest places an inventory item in a player maze, MazeManager sends it when an item is dropped or moved.
type PlaceMazeItemRequest struct {
	InventoryOID       types.OID
	PlayerMazeOID      types.OID
	PlayerMazePieceOID types.OID
	SlotOID            types.OID
	Rotation           string
	ObjectPosition     types.ObjectPosition
	// DynamicAttributes is the yard area name of the item, or "indoor".
	DynamicAttributes []string
}

func (req *PlaceMazeItemRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.InventoryOID)
	reader.ReadObject(&req.PlayerMazeOID)
	reader.ReadObject(&req.PlayerMazePieceOID)
	reader.ReadObject(&req.SlotOID)
	req.Rotation = reader.ReadString()
	reader.ReadObject(&req.ObjectPosition)
	req.DynamicAttributes = gsf.ReadSlice(reader, reader.ReadString)
}

// PlaceMazeItemResponse tells whether the item was placed.
type PlaceMazeItemResponse struct {
	IsUpdated bool
}

func (res *PlaceMazeItemResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsUpdated)
}
