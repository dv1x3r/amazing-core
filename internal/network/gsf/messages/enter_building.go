package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

type EnterBuildingRequest struct {
	LocID       types.OID
	BuildingID  types.OID
	Pos         types.Position
	Orientation types.QTH
}

func (req *EnterBuildingRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.LocID)
	reader.ReadObject(&req.BuildingID)
	reader.ReadObject(&req.Pos)
	reader.ReadObject(&req.Orientation)
}

type EnterBuildingResponse struct {
	BuildingID types.OID
}

func (res *EnterBuildingResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.BuildingID)
}
