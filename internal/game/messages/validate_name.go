package messages

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type ValidateNameRequest struct {
	Name string
}

func (req *ValidateNameRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Name = reader.ReadString()
}

type ValidateNameResponse struct {
	FilterName string
}

func (res *ValidateNameResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteString(res.FilterName)
}
