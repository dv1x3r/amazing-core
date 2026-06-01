package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// StartQuestRequest requests a quest start.
type StartQuestRequest struct {
	PlayerQuestOID types.OID
	QuestOID       types.OID
	LocationOID    types.OID
	SpawnedItems   []types.SpawnedItem
}

func (req *StartQuestRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerQuestOID)
	reader.ReadObject(&req.QuestOID)
	reader.ReadObject(&req.LocationOID)
	req.SpawnedItems = gsf.ReadSlice(reader, func() types.SpawnedItem {
		var value types.SpawnedItem
		reader.ReadObject(&value)
		return value
	})
}

// StartQuestResponse contains started quest status.
type StartQuestResponse struct {
	IsStarted        bool
	SpawnedItems     []types.SpawnedItem
	ConfirmedPlayers []types.OID
}

func (res *StartQuestResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsStarted)
	gsf.WriteSlice(writer, res.SpawnedItems, func(value types.SpawnedItem) {
		writer.WriteObject(&value)
	})
	gsf.WriteSlice(writer, res.ConfirmedPlayers, func(value types.OID) {
		writer.WriteObject(&value)
	})
}
