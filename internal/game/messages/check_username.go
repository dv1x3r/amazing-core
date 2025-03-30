package messages

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type CheckUsernameRequest struct {
	Username string
	Password string
}

func (req *CheckUsernameRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Username = reader.ReadString()
	req.Password = reader.ReadString()
}

type CheckUsernameResponse struct {
}

func (res *CheckUsernameResponse) Serialize(writer gsf.ProtocolWriter) {
}
