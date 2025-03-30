package messages

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type GetRandomNamesRequest struct {
	Amount       int32
	NamePartType string
}

func (req *GetRandomNamesRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Amount = reader.ReadInt32()
	req.NamePartType = reader.ReadString()
}

type GetRandomNamesResponse struct {
	Names []string
}

func (res *GetRandomNamesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Names, writer.WriteString)
}
