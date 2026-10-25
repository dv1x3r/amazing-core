package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// ReloginRequest restores the USER session after a reconnect,
// with the session OID and conversation ID from the LOGIN response.
type ReloginRequest struct {
	SessionOID     types.OID
	ConversationID int64
	LoginName      string
}

func (req *ReloginRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.SessionOID)
	req.ConversationID = reader.ReadInt64()
	req.LoginName = reader.ReadString()
}

// ReloginResponse acknowledges the restored USER session.
type ReloginResponse struct {
}

func (res *ReloginResponse) Serialize(writer gsf.ProtocolWriter) {
}
