package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

type LogoutRequest struct {
}

func (req *LogoutRequest) Deserialize(reader gsf.ProtocolReader) {
}

type LogoutResponse struct {
}

func (res *LogoutResponse) Serialize(writer gsf.ProtocolWriter) {
}
