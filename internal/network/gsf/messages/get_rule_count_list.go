package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetRuleCountListRequest requests rule counter values.
type GetRuleCountListRequest struct {
	Items       []types.OID
	Expressions []string
	OwnerOID    types.OID
	Evaluate    bool
}

func (req *GetRuleCountListRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Items = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
	req.Expressions = gsf.ReadSlice(reader, reader.ReadString)
	reader.ReadObject(&req.OwnerOID)
	req.Evaluate = reader.ReadBool()
}

// GetRuleCountListResponse contains rule counter values paired with the requested OIDs.
type GetRuleCountListResponse struct {
	ItemPairs []types.OIDLongPair
}

func (res *GetRuleCountListResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.ItemPairs, func(value types.OIDLongPair) {
		writer.WriteObject(&value)
	})
}
