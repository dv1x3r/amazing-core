package notify

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// Heartbeat is the player of the SYNC keepalive.
type Heartbeat struct {
	// POID is the player OID carried by the sync heartbeat.
	POID types.OID
}

func (n *Heartbeat) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&n.POID)
}
