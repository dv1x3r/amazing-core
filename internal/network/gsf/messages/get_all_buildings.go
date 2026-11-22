package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetAllBuildingsRequest contains the zone of the buildings.
type GetAllBuildingsRequest struct {
	ZoneOID types.OID
}

func (req *GetAllBuildingsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.ZoneOID)
}

// GetAllBuildingsResponse contains the buildings of the zone.
type GetAllBuildingsResponse struct {
	Buildings []types.Building
}

func (res *GetAllBuildingsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Buildings, func(value types.Building) {
		writer.WriteObject(&value)
	})
}
