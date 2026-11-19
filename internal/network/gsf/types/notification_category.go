package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// NotificationCategory groups the notification definitions.
type NotificationCategory struct {
	OID               OID
	Ordinal           int32
	DefaultExpiryDays int32
}

func (nc *NotificationCategory) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&nc.OID)
	writer.WriteInt32(nc.Ordinal)
	writer.WriteInt32(nc.DefaultExpiryDays)
}
