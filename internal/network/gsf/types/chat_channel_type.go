package types

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/chatchanneltypevalue"
)

// ChatChannelType maps a chat channel object ID to a client channel value.
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
