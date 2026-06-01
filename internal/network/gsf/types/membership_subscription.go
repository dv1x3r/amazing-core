package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

type MembershipSubscription struct {
	AssetContainer
	SKU      string
	Duration int32
}

func (ms *MembershipSubscription) Serialize(writer gsf.ProtocolWriter) {
	ms.AssetContainer.Serialize(writer)
	writer.WriteString(ms.SKU)
	writer.WriteInt32(ms.Duration)
}

func (ms *MembershipSubscription) Deserialize(reader gsf.ProtocolReader) {
	ms.AssetContainer.Deserialize(reader)
	ms.SKU = reader.ReadString()
	ms.Duration = reader.ReadInt32()
}
