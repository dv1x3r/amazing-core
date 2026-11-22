package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// GetStatsTypeRequest is empty.
type GetStatsTypeRequest struct{}

func (req *GetStatsTypeRequest) Deserialize(reader gsf.ProtocolReader) {}

// GetStatsTypeResponse contains the statistic types, always an empty list.
type GetStatsTypeResponse struct{}

func (res *GetStatsTypeResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(0)
}
