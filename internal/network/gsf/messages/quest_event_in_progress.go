package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// QuestEventInProgressRequest is empty.
type QuestEventInProgressRequest struct{}

func (req *QuestEventInProgressRequest) Deserialize(reader gsf.ProtocolReader) {}

// QuestEventInProgressResponse is empty.
type QuestEventInProgressResponse struct{}

func (res *QuestEventInProgressResponse) Serialize(writer gsf.ProtocolWriter) {}
