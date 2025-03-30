package messages

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types"
)

type GetZonesRequest struct {
}

func (req *GetZonesRequest) Deserialize(reader gsf.ProtocolReader) {
}

type GetZonesResponse struct {
	Zones []types.Zone
}

func (res *GetZonesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Zones, func(value types.Zone) {
		writer.WriteObject(&value)
	})
}
