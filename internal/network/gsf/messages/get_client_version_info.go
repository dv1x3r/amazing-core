package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// GetClientVersionInfoRequest identifies the client requesting version policy.
type GetClientVersionInfoRequest struct {
	// The client always sends "AmazingWorld".
	ClientName string
}

func (req *GetClientVersionInfoRequest) Deserialize(reader gsf.ProtocolReader) {
	req.ClientName = reader.ReadString()
}

// GetClientVersionInfoResponse contains the accepted client version policy.
type GetClientVersionInfoResponse struct {
	// ClientVersionInfo string in "<version>.<forceUpdate>" format, e.g. "133852.true".
	//
	// The forceUpdate part is either true (blocking) or false (optional).
	ClientVersionInfo string
}

func (res *GetClientVersionInfoResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteString(res.ClientVersionInfo)
}
