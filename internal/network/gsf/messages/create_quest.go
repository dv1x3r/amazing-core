package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// CreateQuestRequest contains the quest the player takes.
type CreateQuestRequest struct {
	QuestOID types.OID
}

func (req *CreateQuestRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.QuestOID)
}

// CreateQuestResponse contains the new player quest and the player quests of its steps.
type CreateQuestResponse struct {
	PlayerQuestOID    types.OID
	IsLocation        bool
	ParentLocationOID types.OID
	SubQuestMap       map[string]types.OID
}

func (res *CreateQuestResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.PlayerQuestOID)
	writer.WriteBool(res.IsLocation)
	writer.WriteObject(&res.ParentLocationOID)
	gsf.WriteMap(writer, res.SubQuestMap, func(value types.OID) {
		writer.WriteObject(&value)
	})
}
