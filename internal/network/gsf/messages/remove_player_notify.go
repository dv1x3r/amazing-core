package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// RemovePlayerNotify is a SYNC notification payload for a player leaving visibility.
type RemovePlayerNotify struct {
	POID              types.OID
	PlayerVillagerOID types.OID
	LCP               bool
	Type              int32
}

func (n *RemovePlayerNotify) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&n.POID)
	writer.WriteObject(&n.PlayerVillagerOID)
	writer.WriteBool(n.LCP)
	writer.WriteInt32(n.Type)
}

func (n *RemovePlayerNotify) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&n.POID)
	reader.ReadObject(&n.PlayerVillagerOID)
	n.LCP = reader.ReadBool()
	n.Type = reader.ReadInt32()
}
