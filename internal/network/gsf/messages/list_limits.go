package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// ListLimitsRequest is empty.
type ListLimitsRequest struct{}

func (req *ListLimitsRequest) Deserialize(reader gsf.ProtocolReader) {}

// ListLimitsResponse contains the friend list limits, always an empty list.
type ListLimitsResponse struct{}

func (res *ListLimitsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(0)
}
