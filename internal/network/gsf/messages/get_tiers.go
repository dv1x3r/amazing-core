package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetTiersRequest requests subscription tier definitions.
type GetTiersRequest struct {
}

func (req *GetTiersRequest) Deserialize(reader gsf.ProtocolReader) {
}

// GetTiersResponse contains subscription tier definitions.
type GetTiersResponse struct {
	Tiers []types.Tier
}

func (res *GetTiersResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Tiers, func(value types.Tier) {
		writer.WriteObject(&value)
	})
}
