package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetQuestFromParentRequest requests quest for parent.
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

// GetQuestFromParentResponse contains quest record for parent.
type GetQuestFromParentResponse struct {
	Quest types.Quest
}

func (res *GetQuestFromParentResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.Quest)
}
