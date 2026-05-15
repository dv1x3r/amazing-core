package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// VillageRolePlayer assigns a player to a village role.
type VillageRolePlayer struct {
	OID            OID
	VillageRoleOID OID
	PlayerOID      OID
	VillageOID     OID
}

func (vrp *VillageRolePlayer) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&vrp.OID)
	writer.WriteObject(&vrp.VillageRoleOID)
	writer.WriteObject(&vrp.PlayerOID)
	writer.WriteObject(&vrp.VillageOID)
}

func (vrp *VillageRolePlayer) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&vrp.OID)
	reader.ReadObject(&vrp.VillageRoleOID)
	reader.ReadObject(&vrp.PlayerOID)
	reader.ReadObject(&vrp.VillageOID)
}
