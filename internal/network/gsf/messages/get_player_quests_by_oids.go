package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetPlayerQuestsByOIDsRequest contains the OIDs of the quests.
type GetPlayerQuestsByOIDsRequest struct {
	QuestOIDs []types.OID
}

func (req *GetPlayerQuestsByOIDsRequest) Deserialize(reader gsf.ProtocolReader) {
	req.QuestOIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
}

// GetPlayerQuestsByOIDsResponse contains the player quests of the quests.
type GetPlayerQuestsByOIDsResponse struct {
	PlayerQuests []types.PlayerQuest
}

func (res *GetPlayerQuestsByOIDsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.PlayerQuests, func(value types.PlayerQuest) {
		writer.WriteObject(&value)
	})
}
