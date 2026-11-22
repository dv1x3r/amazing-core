package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetPlayerNPCsRequest contains the player and the zone of the NPCs.
type GetPlayerNPCsRequest struct {
	PlayerOID types.OID
	ZoneOID   types.OID
}

func (req *GetPlayerNPCsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerOID)
	reader.ReadObject(&req.ZoneOID)
}

// GetPlayerNPCsResponse contains the NPCs of the zone.
type GetPlayerNPCsResponse struct {
	NPCs []types.NPC
}

func (res *GetPlayerNPCsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.NPCs, func(value types.NPC) {
		writer.WriteObject(&value)
	})
}
