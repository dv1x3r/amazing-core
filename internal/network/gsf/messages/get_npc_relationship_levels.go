package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetNPCRelationshipLevelsRequest requests the points of the relationship levels.
type GetNPCRelationshipLevelsRequest struct{}

func (req *GetNPCRelationshipLevelsRequest) Deserialize(reader gsf.ProtocolReader) {}

// GetNPCRelationshipLevelsResponse contains the points of the relationship levels, NPCStoreSVCS iterates the list.
type GetNPCRelationshipLevelsResponse struct {
	Levels []types.NPCRelationshipLevel
}

func (res *GetNPCRelationshipLevelsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Levels, func(value types.NPCRelationshipLevel) {
		writer.WriteObject(&value)
	})
}
