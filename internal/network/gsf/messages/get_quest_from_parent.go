package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetQuestFromParentRequest contains the parent, e.g. an NPC, and the quest type.
type GetQuestFromParentRequest struct {
	ParentOID       types.OID
	ParentHierarchy string
	QuestTypeOID    types.OID
}

func (req *GetQuestFromParentRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.ParentOID)
	req.ParentHierarchy = reader.ReadString()
	reader.ReadObject(&req.QuestTypeOID)
}

// GetQuestFromParentResponse contains the quest the parent offers.
type GetQuestFromParentResponse struct {
	Quest *types.Quest
}

func (res *GetQuestFromParentResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(res.Quest)
}
