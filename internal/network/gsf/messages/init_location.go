package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// InitLocationRequest contains the location to load.
type InitLocationRequest struct {
	// Player OID for home zone, and 292733975781503755 for Spring Bay zone.
	LocOID types.OID
}

func (req *InitLocationRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.LocOID)
}

// InitLocationResponse contains the location data and the SYNC server address.
type InitLocationResponse struct {
	ZoneInstance *types.ZoneInstance
	Village      *types.Village
	Home         *types.PlayerHome

	// Token used to authenticate with the SYNC server.
	SyncServerToken string

	// IP address of the SYNC server.
	SyncServerIP string

	// Port of the SYNC server.
	SyncServerPort int32
}

func (res *InitLocationResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(res.ZoneInstance)
	writer.WriteObject(res.Village)
	writer.WriteObject(res.Home)
	writer.WriteString(res.SyncServerToken)
	writer.WriteString(res.SyncServerIP)
	writer.WriteInt32(res.SyncServerPort)
}
