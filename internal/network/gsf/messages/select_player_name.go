package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// SelectPlayerNameRequest contains the chosen family name.
type SelectPlayerNameRequest struct {
	Name string
}

func (req *SelectPlayerNameRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Name = reader.ReadString()
}

// SelectPlayerNameResponse is empty, AppCode 71 means the name is taken.
type SelectPlayerNameResponse struct {
}

func (res *SelectPlayerNameResponse) Serialize(writer gsf.ProtocolWriter) {
}
