package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// SelectPlayerNameRequest submits the selected family name during registration.
type SelectPlayerNameRequest struct {
	Name string
}

func (req *SelectPlayerNameRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Name = reader.ReadString()
}

// SelectPlayerNameResponse success is implied by a non-error response.
//
// On AppCode 71 (duplicate) the client shows an error dialog.
type SelectPlayerNameResponse struct {
}

func (res *SelectPlayerNameResponse) Serialize(writer gsf.ProtocolWriter) {
}
