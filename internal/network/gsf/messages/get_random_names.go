package messages

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// GetRandomNamesRequest requests generated name parts by name-part type.
type GetRandomNamesRequest struct {
	// Number of names to return.
	Amount int32

	// "second_name" for Zing names.
	// "Family_1", "Family_2", "Family_3" for family name parts.
	NamePartType string
}

func (req *GetRandomNamesRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Amount = reader.ReadInt32()
	req.NamePartType = reader.ReadString()
}

// GetRandomNamesResponse contains generated name parts.
type GetRandomNamesResponse struct {
	Names []string
}

func (res *GetRandomNamesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Names, writer.WriteString)
}
