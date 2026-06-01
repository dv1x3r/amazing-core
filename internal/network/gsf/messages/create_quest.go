package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// CreateQuestRequest requests a quest initialization.
type CreateQuestRequest struct {
	QuestOID types.OID
}

func (req *CreateQuestRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.QuestOID)
}

// CreateQuestResponse contains new quest ID's.
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
