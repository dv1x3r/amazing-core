package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetNPCsWithQuestOfferRequest is empty.
type GetNPCsWithQuestOfferRequest struct{}

func (req *GetNPCsWithQuestOfferRequest) Deserialize(reader gsf.ProtocolReader) {}

// GetNPCsWithQuestOfferResponse contains the NPCs that offer a quest to the player.
type GetNPCsWithQuestOfferResponse struct {
	// NPCs is nil when none: TaskGUI takes a random element of a non-nil list and fails on an empty one.
	NPCs []types.NPC
}

func (res *GetNPCsWithQuestOfferResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.NPCs, func(value types.NPC) {
		writer.WriteObject(&value)
	})
}
