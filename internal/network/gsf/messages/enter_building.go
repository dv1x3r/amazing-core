package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// EnterBuildingRequest is sent when the client enters a building location.
type EnterBuildingRequest struct {
	LocOID      types.OID
	BuildingOID types.OID
	Pos         types.Position
	Orientation types.QTH
}

func (req *EnterBuildingRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.LocOID)
	reader.ReadObject(&req.BuildingOID)
	reader.ReadObject(&req.Pos)
	reader.ReadObject(&req.Orientation)
}

// EnterBuildingResponse contains the building ID for the entered building.
type EnterBuildingResponse struct {
	BuildingOID types.OID
}

func (res *EnterBuildingResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.BuildingOID)
}
