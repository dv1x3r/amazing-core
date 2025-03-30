package messages

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type SelectPlayerNameRequest struct {
	Name string
}

func (req *SelectPlayerNameRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Name = reader.ReadString()
}

type SelectPlayerNameResponse struct {
}

func (res *SelectPlayerNameResponse) Serialize(writer gsf.ProtocolWriter) {
}
