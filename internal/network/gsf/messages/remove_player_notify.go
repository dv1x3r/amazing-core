package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

type RemovePlayerNotify struct {
	PID              types.OID
	PlayerVillagerID types.OID
	LCP              bool
	Type             int32
}

func (n *RemovePlayerNotify) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&n.PID)
	writer.WriteObject(&n.PlayerVillagerID)
	writer.WriteBool(n.LCP)
	writer.WriteInt32(n.Type)
}

func (n *RemovePlayerNotify) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&n.PID)
	reader.ReadObject(&n.PlayerVillagerID)
	n.LCP = reader.ReadBool()
	n.Type = reader.ReadInt32()
}
