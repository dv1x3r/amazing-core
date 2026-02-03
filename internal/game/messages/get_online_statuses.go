package messages

import (
	"github.com/dv1x3r/amazing-core/internal/game/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
)

type GetOnlineStatusesRequest struct {
}

func (req *GetOnlineStatusesRequest) Deserialize(reader gsf.ProtocolReader) {
}

type GetOnlineStatusesResponse struct {
	OnlineStatuses []types.OnlineStatus
}

func (res *GetOnlineStatusesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.OnlineStatuses, func(value types.OnlineStatus) {
		writer.WriteObject(&value)
	})
}
