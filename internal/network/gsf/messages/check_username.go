package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// CheckUsernameRequest is the registration username availability check payload.
type CheckUsernameRequest struct {
	Username string
	Password string
}

func (req *CheckUsernameRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Username = reader.ReadString()
	req.Password = reader.ReadString()
}

// CheckUsernameResponse success is implied by a non-error response.
//
// On AppCode 301 the server signals the username is already in use.
type CheckUsernameResponse struct {
}

func (res *CheckUsernameResponse) Serialize(writer gsf.ProtocolWriter) {
}
