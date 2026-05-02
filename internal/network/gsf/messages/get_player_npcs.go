package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetPlayerNPCsRequest requests NPCs for a player and zone.
type GetPlayerNPCsRequest struct {
	PlayerID types.OID
	ZoneID   types.OID
}

func (req *GetPlayerNPCsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerID)
	reader.ReadObject(&req.ZoneID)
}

// GetPlayerNPCsResponse contains NPC records for the requested zone.
type GetPlayerNPCsResponse struct {
	NPCs []types.NPC
}

func (res *GetPlayerNPCsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.NPCs, func(value types.NPC) {
		writer.WriteObject(&value)
	})
}
