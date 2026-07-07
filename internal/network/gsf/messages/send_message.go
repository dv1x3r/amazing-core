package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// SendMessageRequest is the client chat send payload.
type SendMessageRequest struct {
	Message           string
	TypeID            types.OID
	RecipientID       types.OID
	PlayerChatGroupID types.OID
}

func (req *SendMessageRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Message = reader.ReadString()
	reader.ReadObject(&req.TypeID)
	reader.ReadObject(&req.RecipientID)
	reader.ReadObject(&req.PlayerChatGroupID)
}

// SendMessageResponse returns the moderated chat message payload.
type SendMessageResponse struct {
	CrispDataTO     *types.CrispDataTO
	FilteredMessage string
	Categories      []string
	OffenceGroup    int32
	RecipientID     types.OID
}

func (res *SendMessageResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(res.CrispDataTO)
	writer.WriteString(res.FilteredMessage)
	gsf.WriteSlice(writer, res.Categories, writer.WriteString)
	writer.WriteInt32(res.OffenceGroup)
	writer.WriteObject(&res.RecipientID)
}
