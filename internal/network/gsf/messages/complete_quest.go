package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// CompleteQuestRequest requests a quest completion.
type CompleteQuestRequest struct {
	PlayerQuestOID types.OID
	QuestOID       types.OID
	LocationOID    types.OID
}

func (req *CompleteQuestRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerQuestOID)
	reader.ReadObject(&req.QuestOID)
	reader.ReadObject(&req.LocationOID)
}

// CompleteQuestResponse contains completed quest status.
type CompleteQuestResponse struct {
	IsCompleted bool
	GiftBox     types.GiftBox
	AwardSet    []types.QuestAwardElement
}

func (res *CompleteQuestResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsCompleted)
	writer.WriteObject(&res.GiftBox)
	gsf.WriteSlice(writer, res.AwardSet, func(value types.QuestAwardElement) {
		writer.WriteObject(&value)
	})
}
