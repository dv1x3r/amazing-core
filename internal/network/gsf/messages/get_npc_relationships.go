package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetNPCRelationshipsRequest contains the zone of the NPCs, nil for every NPC.
type GetNPCRelationshipsRequest struct {
	ZoneOID *types.OID
}

func (req *GetNPCRelationshipsRequest) Deserialize(reader gsf.ProtocolReader) {
	var zoneOID types.OID
	if reader.ReadObject(&zoneOID) {
		req.ZoneOID = &zoneOID
	}
}

// GetNPCRelationshipsResponse contains the relationships of the player with the NPCs.
type GetNPCRelationshipsResponse struct {
	// Relationships is never nil, NPCManager iterates it.
	Relationships []types.NPCRelationship
}

func (res *GetNPCRelationshipsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Relationships, func(value types.NPCRelationship) {
		writer.WriteObject(&value)
	})
}
