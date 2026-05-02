package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetMazeItemsRequest requests item instances placed in a player maze.
type GetMazeItemsRequest struct {
	PlayerMazeID types.OID
	PlayerID     types.OID
}

func (req *GetMazeItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerMazeID)
	reader.ReadObject(&req.PlayerID)
}

// GetMazeItemsResponse contains item instances placed in a player maze.
type GetMazeItemsResponse struct {
	MazeItems []types.PlayerItem
}

func (res *GetMazeItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.MazeItems, func(value types.PlayerItem) {
		writer.WriteObject(&value)
	})
}
