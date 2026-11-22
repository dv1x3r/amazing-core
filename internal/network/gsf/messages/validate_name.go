package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// ValidateNameRequest contains the Zing name to check.
type ValidateNameRequest struct {
	Name string
}

func (req *ValidateNameRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Name = reader.ReadString()
}

// ValidateNameResponse contains the filtered Zing name.
type ValidateNameResponse struct {
	// The client does not really care about the value.
	//
	// A non-empty value means the name was rejected.
	FilterName string
}

func (res *ValidateNameResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteString(res.FilterName)
}
