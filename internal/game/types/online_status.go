package types

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types/onlinestatusvalue"
)

type OnlineStatus struct {
	OID   OID
	Value onlinestatusvalue.OnlineStatusValue
}

func (os *OnlineStatus) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&os.OID)
	writer.WriteString(os.Value.String())
}

func (os *OnlineStatus) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&os.OID)
	os.Value = onlinestatusvalue.Parse(reader.ReadString())
}
