package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetPlayerQuestsByOIDsRequest requests quest definitions by OIDs.
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

// GetPlayerQuestsByOIDsResponse contains quest definitions for the requested OIDs.
type GetPlayerQuestsByOIDsResponse struct {
	PlayerQuests []types.PlayerQuest
}

func (res *GetPlayerQuestsByOIDsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.PlayerQuests, func(value types.PlayerQuest) {
		writer.WriteObject(&value)
	})
}
