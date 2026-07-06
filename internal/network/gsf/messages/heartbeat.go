package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// HeartbeatRequest is the periodic USER server keepalive payload.
type HeartbeatRequest struct {
	RunPercent     float32
	XCoordinate    float32
	YCoordinate    float32
	ZCoordinate    float32
	InactivityTime int64
}

func (req *HeartbeatRequest) Deserialize(reader gsf.ProtocolReader) {
	req.RunPercent = reader.ReadFloat32()
	req.XCoordinate = reader.ReadFloat32()
	req.YCoordinate = reader.ReadFloat32()
	req.ZCoordinate = reader.ReadFloat32()
	req.InactivityTime = reader.ReadInt64()
}

// HeartbeatResponse refreshes player stats, server time, and queued notifications.
type HeartbeatResponse struct {
	PlayerStats       []types.PlayerStats
	CurrentServerTime gsf.UnixTime
	QueueNotify       []types.PlayerNotify
}

func (res *HeartbeatResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.PlayerStats, func(value types.PlayerStats) {
		writer.WriteObject(&value)
	})
	writer.WriteUtcDate(res.CurrentServerTime)
	gsf.WriteSlice(writer, res.QueueNotify, func(value types.PlayerNotify) {
		writer.WriteObject(&value)
	})
}
