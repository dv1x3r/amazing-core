package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetCMSNotificationsRequest requests the notification definitions.
type GetCMSNotificationsRequest struct{}

func (req *GetCMSNotificationsRequest) Deserialize(reader gsf.ProtocolReader) {
}

// GetCMSNotificationsResponse contains the notification definitions.
type GetCMSNotificationsResponse struct {
	Notifications []types.Notification
}

func (res *GetCMSNotificationsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Notifications, func(value types.Notification) {
		writer.WriteObject(&value)
	})
}
