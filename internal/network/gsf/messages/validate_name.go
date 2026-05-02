package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// ValidateNameRequest submits a Zing name for validation.
type ValidateNameRequest struct {
	Name string
}

func (req *ValidateNameRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Name = reader.ReadString()
}

// ValidateNameResponse contains the validation filter result for a Zing name.
type ValidateNameResponse struct {
	// The client does not really care about the value.
	//
	// A non-empty value means the name was rejected.
	FilterName string
}

func (res *ValidateNameResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteString(res.FilterName)
}
