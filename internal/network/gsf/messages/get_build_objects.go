package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetBuildObjectsRequest is empty.
type GetBuildObjectsRequest struct {
}

func (req *GetBuildObjectsRequest) Deserialize(reader gsf.ProtocolReader) {
}

// GetBuildObjectsResponse contains the maze pieces of the player.
type GetBuildObjectsResponse struct {
	PlayerBuildObjects []types.PlayerBuildObject
}

func (res *GetBuildObjectsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.PlayerBuildObjects, func(value types.PlayerBuildObject) {
		writer.WriteObject(&value)
	})
}
