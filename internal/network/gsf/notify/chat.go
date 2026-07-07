package notify

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// Chat is the payload for GSFChatNotify.
type Chat struct {
	Msg          string
	TypeID       types.OID
	GroupID      types.OID
	SenderID     types.OID
	Categories   []string
	OffenceGroup int32
}

func (n *Chat) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteString(n.Msg)
	writer.WriteObject(&n.TypeID)
	writer.WriteObject(&n.GroupID)
	writer.WriteObject(&n.SenderID)
	gsf.WriteSlice(writer, n.Categories, writer.WriteString)
	writer.WriteInt32(n.OffenceGroup)
}
