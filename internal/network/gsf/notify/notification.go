package notify

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// Notification is the payload for GSFNotificationNotify.
type Notification struct {
	// OID is the notification definition OID from GET_CMS_NOTIFICATIONS.
	OID types.OID

	// POID is the player the notification comes from.
	POID types.OID

	PlayerNotificationOID *types.OID

	// ObjectOID is the notification subject, e.g. the player quest of QuestCompleted.
	ObjectOID types.OID

	ObjectOID2        *types.OID
	Source            string
	SendingObjectOID  *types.OID
	Quantity          int32
	CurrentServerTime gsf.UnixTime
	TargetTime        gsf.UnixTime
}

func (n *Notification) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&n.OID)
	writer.WriteObject(&n.POID)
	writer.WriteObject(n.PlayerNotificationOID)
	writer.WriteObject(&n.ObjectOID)
	writer.WriteObject(n.ObjectOID2)
	writer.WriteString(n.Source)
	writer.WriteObject(n.SendingObjectOID)
	writer.WriteInt32(n.Quantity)
	writer.WriteUtcDate(n.CurrentServerTime)
	writer.WriteUtcDate(n.TargetTime)
}
