package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// LogoutRequest requests termination of the current USER server session.
type LogoutRequest struct {
}

func (req *LogoutRequest) Deserialize(reader gsf.ProtocolReader) {
}

// LogoutResponse acknowledges termination of the current USER server session.
type LogoutResponse struct {
}

func (res *LogoutResponse) Serialize(writer gsf.ProtocolWriter) {
}
