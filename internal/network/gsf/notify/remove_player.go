package notify

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// RemovePlayer is the payload for GSFRemovePlayerNotify.
type RemovePlayer struct {
	// POID is the player OID.
	POID types.OID

	PlayerVillagerOID types.OID

	// LCP is an unused client flag.
	LCP  bool
	Type int32
}

func (n *RemovePlayer) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&n.POID)
	writer.WriteObject(&n.PlayerVillagerOID)
	writer.WriteBool(n.LCP)
	writer.WriteInt32(n.Type)
}
