package messages

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types"
)

type GetPlayerNPCsRequest struct {
	PlayerID types.OID
	ZoneID   types.OID
}

func (req *GetPlayerNPCsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerID)
	reader.ReadObject(&req.ZoneID)
}

type GetPlayerNPCsResponse struct {
	NPCs []types.NPC
}

func (res *GetPlayerNPCsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.NPCs, func(value types.NPC) {
		writer.WriteObject(&value)
	})
}
