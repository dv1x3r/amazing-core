package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// AddQuestItemRequest contains the collected quest items and where they were collected.
type AddQuestItemRequest struct {
	PlayerQuestOID   types.OID
	QuestOID         types.OID
	ItemInstanceUUID string
	// ItemOID is the Quest Item container of the quest Treenode.
	ItemOID              types.OID
	Count                int32
	SpawnPoint           string
	ZoneOID              types.OID
	ObjectPosition       types.ObjectPosition
	PlacedNPCOID         types.OID
	PlacedZoneOID        types.OID
	PlacedText           string
	PlacedObjectPosition types.ObjectPosition
}

func (req *AddQuestItemRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerQuestOID)
	reader.ReadObject(&req.QuestOID)
	req.ItemInstanceUUID = reader.ReadString()
	reader.ReadObject(&req.ItemOID)
	req.Count = reader.ReadInt32()
	req.SpawnPoint = reader.ReadString()
	reader.ReadObject(&req.ZoneOID)
	reader.ReadObject(&req.ObjectPosition)
	reader.ReadObject(&req.PlacedNPCOID)
	reader.ReadObject(&req.PlacedZoneOID)
	req.PlacedText = reader.ReadString()
	reader.ReadObject(&req.PlacedObjectPosition)
}

// AddQuestItemResponse contains the quest inventory OIDs of the collected items.
type AddQuestItemResponse struct {
	QuestInventoryOIDs []types.OID
}

func (res *AddQuestItemResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.QuestInventoryOIDs, func(value types.OID) {
		writer.WriteObject(&value)
	})
}
