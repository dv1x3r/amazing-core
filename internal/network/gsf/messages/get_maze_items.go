package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetMazeItemsRequest contains the player maze of the items.
type GetMazeItemsRequest struct {
	PlayerMazeOID types.OID
	PlayerOID     types.OID
}

func (req *GetMazeItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerMazeOID)
	reader.ReadObject(&req.PlayerOID)
}

// GetMazeItemsResponse contains the items placed in the player maze.
type GetMazeItemsResponse struct {
	MazeItems []types.PlayerItem
}

func (res *GetMazeItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.MazeItems, func(value types.PlayerItem) {
		writer.WriteObject(&value)
	})
}
