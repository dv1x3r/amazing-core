package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

type GetTiersRequest struct {
}

func (req *GetTiersRequest) Deserialize(reader gsf.ProtocolReader) {
}

type GetTiersResponse struct {
	Tiers []types.Tier
}

func (res *GetTiersResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Tiers, func(value types.Tier) {
		writer.WriteObject(&value)
	})
}
