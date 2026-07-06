package notify

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// ChangeObject is the payload for GSFChangeObjectNotify.
type ChangeObject struct {
	InstanceOID types.OID
	TemplateOID types.OID

	// PlayerOID is the owner/player OID carried as pid.
	PlayerOID types.OID

	// VillageOID is the village/context OID carried as vid.
	VillageOID  types.OID
	LocationOID types.OID

	// Ver is the player-details cache version used by clients to refetch appearance data.
	Ver int32

	Pos   types.Position
	QTH   types.QTH
	State map[string]string
}

func (n *ChangeObject) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&n.InstanceOID)
	writer.WriteObject(&n.TemplateOID)
	writer.WriteObject(&n.PlayerOID)
	writer.WriteObject(&n.VillageOID)
	writer.WriteObject(&n.LocationOID)
	writer.WriteInt32(n.Ver)
	writer.WriteObject(&n.Pos)
	writer.WriteObject(&n.QTH)
	gsf.WriteMap(writer, n.State, writer.WriteString)
}
