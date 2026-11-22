package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// ReloginRequest contains the session to restore after a reconnect.
type ReloginRequest struct {
	// SessionOID and ConversationID come from the LOGIN response.
	SessionOID     types.OID
	ConversationID int64
	LoginName      string
}

func (req *ReloginRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.SessionOID)
	req.ConversationID = reader.ReadInt64()
	req.LoginName = reader.ReadString()
}

// ReloginResponse is empty.
type ReloginResponse struct {
}

func (res *ReloginResponse) Serialize(writer gsf.ProtocolWriter) {
}
