package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetPlayerStatsRequest requests the player statistics, the stores read the player money from them.
type GetPlayerStatsRequest struct{}

func (req *GetPlayerStatsRequest) Deserialize(reader gsf.ProtocolReader) {}

// GetPlayerStatsResponse contains the player statistics.
type GetPlayerStatsResponse struct {
	PlayerStats []types.PlayerStats
}

func (res *GetPlayerStatsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.PlayerStats, func(value types.PlayerStats) {
		writer.WriteObject(&value)
	})
}
