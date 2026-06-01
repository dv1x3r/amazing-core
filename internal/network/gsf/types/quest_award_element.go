package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// QuestAwardElement is a quest award instance.
type QuestAwardElement struct {
	OID            OID
	Delta          int32
	ObjectType     int32
	NPCOID         OID
	Asset          AssetContainer
	Source         string
	PlayerAvatar   PlayerAvatar
	PlayerItem     PlayerItem
	PlayerMaze     PlayerMaze
	RuleAttributes []string
}

func (qae *QuestAwardElement) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&qae.OID)
	writer.WriteInt32(qae.Delta)
	writer.WriteInt32(qae.ObjectType)
	writer.WriteObject(&qae.NPCOID)
	writer.WriteObject(&qae.Asset)
	writer.WriteString(qae.Source)
	writer.WriteObject(&qae.PlayerAvatar)
	writer.WriteObject(&qae.PlayerItem)
	writer.WriteObject(&qae.PlayerMaze)
	gsf.WriteSlice(writer, qae.RuleAttributes, writer.WriteString)
}

func (qae *QuestAwardElement) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&qae.OID)
	qae.Delta = reader.ReadInt32()
	qae.ObjectType = reader.ReadInt32()
	reader.ReadObject(&qae.NPCOID)
	reader.ReadObject(&qae.Asset)
	qae.Source = reader.ReadString()
	reader.ReadObject(&qae.PlayerAvatar)
	reader.ReadObject(&qae.PlayerItem)
	reader.ReadObject(&qae.PlayerMaze)
	qae.RuleAttributes = gsf.ReadSlice(reader, reader.ReadString)
}
