package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetPlayerQuestsRequest contains the NPC and the page of the player quests.
type GetPlayerQuestsRequest struct {
	Limit  int32
	Offset int32
	Action int32
	NPCOID types.OID
}

func (req *GetPlayerQuestsRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Limit = reader.ReadInt32()
	req.Offset = reader.ReadInt32()
	req.Action = reader.ReadInt32()
	reader.ReadObject(&req.NPCOID)
}

// GetPlayerQuestsResponse contains the player quests of the NPC.
type GetPlayerQuestsResponse struct {
	PlayerQuests []types.PlayerQuest
}

func (res *GetPlayerQuestsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.PlayerQuests, func(value types.PlayerQuest) {
		writer.WriteObject(&value)
	})
}
