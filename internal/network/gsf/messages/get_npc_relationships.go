package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetNPCRelationshipsRequest requests the relationships of the player with the NPCs of a zone, of all NPCs without a zone.
type GetNPCRelationshipsRequest struct {
	ZoneOID *types.OID
}

func (req *GetNPCRelationshipsRequest) Deserialize(reader gsf.ProtocolReader) {
	var zoneOID types.OID
	if reader.ReadObject(&zoneOID) {
		req.ZoneOID = &zoneOID
	}
}

// GetNPCRelationshipsResponse contains the relationships, NPCManager iterates the list, so it is never null.
type GetNPCRelationshipsResponse struct {
	Relationships []types.NPCRelationship
}

func (res *GetNPCRelationshipsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Relationships, func(value types.NPCRelationship) {
		writer.WriteObject(&value)
	})
}
