package types

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/chatchanneltypevalue"
)

// ChatChannelType is a chat channel OID and its client value.
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
