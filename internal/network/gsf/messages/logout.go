package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// LogoutRequest is empty.
type LogoutRequest struct {
}

func (req *LogoutRequest) Deserialize(reader gsf.ProtocolReader) {
}

// LogoutResponse is empty.
type LogoutResponse struct {
}

func (res *LogoutResponse) Serialize(writer gsf.ProtocolWriter) {
}
