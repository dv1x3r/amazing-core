package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// AbandonQuestRequest requests to abandon or fail a quest.
type AbandonQuestRequest struct {
	PlayerQuestOID types.OID
	QuestOID       types.OID
	LocationOID    types.OID
	IsFailed       bool
}

func (req *AbandonQuestRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerQuestOID)
	reader.ReadObject(&req.QuestOID)
	reader.ReadObject(&req.LocationOID)
	req.IsFailed = reader.ReadBool()
}

// AbandonQuestResponse contains abandoned quest status.
type AbandonQuestResponse struct {
	IsAbandoned bool
}

func (res *AbandonQuestResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsAbandoned)
}
