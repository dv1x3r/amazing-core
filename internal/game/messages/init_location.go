package messages

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types"
)

type InitLocationRequest struct {
	LocID types.OID
}

func (req *InitLocationRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.LocID)
}

type InitLocationResponse struct {
	ZoneInstance    types.ZoneInstance
	Village         types.Village
	Home            types.PlayerHome
	SyncServerToken string
	SyncServerIP    string
	SyncServerPort  int32
}

func (res *InitLocationResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.ZoneInstance)
	writer.WriteObject(&res.Village)
	writer.WriteObject(&res.Home)
	writer.WriteString(res.SyncServerToken)
	writer.WriteString(res.SyncServerIP)
	writer.WriteInt32(res.SyncServerPort)
}
