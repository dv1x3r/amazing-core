package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// GetFriendshipRequestCountsRequest is empty.
type GetFriendshipRequestCountsRequest struct{}

func (req *GetFriendshipRequestCountsRequest) Deserialize(reader gsf.ProtocolReader) {}

// GetFriendshipRequestCountsResponse contains the numbers of the sent and the received friendship requests.
type GetFriendshipRequestCountsResponse struct {
	SentCount     int32
	ReceivedCount int32
}

func (res *GetFriendshipRequestCountsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(res.SentCount)
	writer.WriteInt32(res.ReceivedCount)
}
