package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetRuleCountListRequest contains the objects and the expressions to count.
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

// GetRuleCountListResponse contains the counts of the requested objects.
type GetRuleCountListResponse struct {
	ItemPairs []types.OIDLongPair
}

func (res *GetRuleCountListResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.ItemPairs, func(value types.OIDLongPair) {
		writer.WriteObject(&value)
	})
}
