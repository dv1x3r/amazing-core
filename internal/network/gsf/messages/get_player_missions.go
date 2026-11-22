package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetPlayerMissionsRequest contains the page and the hierarchies of the player missions.
type GetPlayerMissionsRequest struct {
	Limit       int32
	Offset      int32
	Action      int32
	Hierarchies []string
}

func (req *GetPlayerMissionsRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Limit = reader.ReadInt32()
	req.Offset = reader.ReadInt32()
	req.Action = reader.ReadInt32()
	req.Hierarchies = gsf.ReadSlice(reader, reader.ReadString)
}

// GetPlayerMissionsResponse contains the player missions.
type GetPlayerMissionsResponse struct {
	PlayerMissions []types.PlayerQuest
}

func (res *GetPlayerMissionsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.PlayerMissions, func(value types.PlayerQuest) {
		writer.WriteObject(&value)
	})
}
