package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetOnlineStatusesRequest requests friend online status values.
type GetOnlineStatusesRequest struct {
}

func (req *GetOnlineStatusesRequest) Deserialize(reader gsf.ProtocolReader) {
}

// GetOnlineStatusesResponse contains friend online status values.
type GetOnlineStatusesResponse struct {
	OnlineStatuses []types.OnlineStatus
}

func (res *GetOnlineStatusesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.OnlineStatuses, func(value types.OnlineStatus) {
		writer.WriteObject(&value)
	})
}
