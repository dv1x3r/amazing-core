package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// QuestItem is a quest item instance.
type QuestItem struct {
	RuleContainer
	Item Item
}

func (qi *QuestItem) Serialize(writer gsf.ProtocolWriter) {
	qi.RuleContainer.Serialize(writer)
	writer.WriteObject(&qi.Item)
}

func (qi *QuestItem) Deserialize(reader gsf.ProtocolReader) {
	qi.RuleContainer.Deserialize(reader)
	reader.ReadObject(&qi.Item)
}
