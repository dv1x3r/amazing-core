package types

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/notificationtype"
)

// Notification is a notification definition.
type Notification struct {
	OID              OID
	Category         NotificationCategory
	NotificationType notificationtype.NotificationType
	RequiresEmail    bool
}

func (n *Notification) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&n.OID)
	writer.WriteObject(&n.Category)
	writer.WriteInt32(int32(n.NotificationType))
	writer.WriteBool(n.RequiresEmail)
}
