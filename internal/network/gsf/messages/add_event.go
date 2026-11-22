package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// AddEventRequest contains a rule engine event of the player and its objects.
type AddEventRequest struct {
	EventOID               types.OID
	EventObjectOID         types.OID
	TargetObjectOID        types.OID
	TargetObjectTemplateID types.OID
	ContextOIDs            []types.OID
	DynamicAttributes      []string
	ReturnRewards          gsf.Null[bool]
	Count                  int32
}

func (req *AddEventRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.EventOID)
	reader.ReadObject(&req.EventObjectOID)
	reader.ReadObject(&req.TargetObjectOID)
	reader.ReadObject(&req.TargetObjectTemplateID)
	req.ContextOIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
	req.DynamicAttributes = gsf.ReadSlice(reader, reader.ReadString)
	req.ReturnRewards = gsf.ReadNullable(reader, reader.ReadBool)
	req.Count = reader.ReadInt32()
}

// AddEventResponse contains the stored player event and its rewards.
type AddEventResponse struct {
	PlayerEventOID     *types.OID
	QuestAwardElements []types.QuestAwardElement
}

func (res *AddEventResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(res.PlayerEventOID)
	gsf.WriteSlice(writer, res.QuestAwardElements, func(value types.QuestAwardElement) {
		writer.WriteObject(&value)
	})
}
