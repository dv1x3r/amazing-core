package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type VillageRolePlayer struct {
	OID           OID
	VillageRoleID OID
	PlayerID      OID
	VillageID     OID
}

func (vrp *VillageRolePlayer) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&vrp.OID)
	writer.WriteObject(&vrp.VillageRoleID)
	writer.WriteObject(&vrp.PlayerID)
	writer.WriteObject(&vrp.VillageID)
}

func (vrp *VillageRolePlayer) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&vrp.OID)
	reader.ReadObject(&vrp.VillageRoleID)
	reader.ReadObject(&vrp.PlayerID)
	reader.ReadObject(&vrp.VillageID)
}
