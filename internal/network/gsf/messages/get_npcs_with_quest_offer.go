package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetNPCsWithQuestOfferRequest requests the NPCs that offer a quest to the player.
type GetNPCsWithQuestOfferRequest struct{}

func (req *GetNPCsWithQuestOfferRequest) Deserialize(reader gsf.ProtocolReader) {}

// GetNPCsWithQuestOfferResponse contains the NPCs that offer a quest, TaskGUI shows the portrait of a random one.
// No NPCs is a null list, TaskGUI takes a random element of a non-null list and fails on an empty one.
type GetNPCsWithQuestOfferResponse struct {
	NPCs []types.NPC
}

func (res *GetNPCsWithQuestOfferResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.NPCs, func(value types.NPC) {
		writer.WriteObject(&value)
	})
}
