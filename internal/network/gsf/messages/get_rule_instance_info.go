package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetRuleInstanceInfoRequest requests the rules of an instance, e.g. a player quest.
type GetRuleInstanceInfoRequest struct {
	InstanceOID types.OID
	TemplateOID types.OID
	LocationOID types.OID
	Details     bool
	OwnerOID    types.OID
}

func (req *GetRuleInstanceInfoRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.InstanceOID)
	reader.ReadObject(&req.TemplateOID)
	reader.ReadObject(&req.LocationOID)
	req.Details = reader.ReadBool()
	reader.ReadObject(&req.OwnerOID)
}

// GetRuleInstanceInfoResponse contains the rules of the instance.
type GetRuleInstanceInfoResponse struct {
	Property    types.RuleProperty
	InstanceOID types.OID
}

func (res *GetRuleInstanceInfoResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.Property)
	writer.WriteObject(&res.InstanceOID)
}
