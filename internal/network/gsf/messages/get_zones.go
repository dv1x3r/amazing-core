package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetZonesRequest is empty.
type GetZonesRequest struct {
}

func (req *GetZonesRequest) Deserialize(reader gsf.ProtocolReader) {
}

// GetZonesResponse contains the zones.
type GetZonesResponse struct {
	Zones []types.Zone
}

func (res *GetZonesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Zones, func(value types.Zone) {
		writer.WriteObject(&value)
	})
}
