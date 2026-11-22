package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// CheckUsernameRequest contains the username to check at registration.
type CheckUsernameRequest struct {
	Username string
	Password string
}

func (req *CheckUsernameRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Username = reader.ReadString()
	req.Password = reader.ReadString()
}

// CheckUsernameResponse is empty, AppCode 301 means the username is taken.
type CheckUsernameResponse struct {
}

func (res *CheckUsernameResponse) Serialize(writer gsf.ProtocolWriter) {
}
