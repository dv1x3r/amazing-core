package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetRuleInfoRequest requests the rules of an object, e.g. a quest.
type GetRuleInfoRequest struct {
	ObjectOID types.OID
	Details   bool
	OwnerOID  types.OID
}

func (req *GetRuleInfoRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.ObjectOID)
	req.Details = reader.ReadBool()
	reader.ReadObject(&req.OwnerOID)
}

// GetRuleInfoResponse contains the rules of the object.
type GetRuleInfoResponse struct {
	Property  types.RuleProperty
	ObjectOID types.OID
}

func (res *GetRuleInfoResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.Property)
	writer.WriteObject(&res.ObjectOID)
}
