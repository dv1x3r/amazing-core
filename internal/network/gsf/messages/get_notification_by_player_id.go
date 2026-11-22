package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// GetNotificationByPlayerIDRequest is empty.
type GetNotificationByPlayerIDRequest struct{}

func (req *GetNotificationByPlayerIDRequest) Deserialize(reader gsf.ProtocolReader) {}

// GetNotificationByPlayerIDResponse contains the stored notifications, always an empty list.
type GetNotificationByPlayerIDResponse struct{}

func (res *GetNotificationByPlayerIDResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(0)
}
