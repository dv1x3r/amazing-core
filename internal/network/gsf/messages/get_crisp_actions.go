package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// GetCrispActionsRequest is empty.
type GetCrispActionsRequest struct{}

func (req *GetCrispActionsRequest) Deserialize(reader gsf.ProtocolReader) {}

// GetCrispActionsResponse contains the Crisp actions, always an empty list.
type GetCrispActionsResponse struct{}

func (res *GetCrispActionsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(0)
}
