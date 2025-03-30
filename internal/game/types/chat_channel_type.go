package types

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types/chatchanneltypevalue"
)

type ChatChannelType struct {
	OID   OID
	Value chatchanneltypevalue.ChatChannelTypeValue
}

func (cct *ChatChannelType) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&cct.OID)
	writer.WriteString(cct.Value.String())
}

func (cct *ChatChannelType) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&cct.OID)
	cct.Value = chatchanneltypevalue.Parse(reader.ReadString())
}
