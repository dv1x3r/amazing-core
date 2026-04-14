package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

type GetClientVersionInfoRequest struct {
	ClientName string
}

func (req *GetClientVersionInfoRequest) Deserialize(reader gsf.ProtocolReader) {
	req.ClientName = reader.ReadString()
}

type GetClientVersionInfoResponse struct {
	ClientVersionInfo string
}

func (res *GetClientVersionInfoResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteString(res.ClientVersionInfo)
}
