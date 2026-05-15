package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// Tier is a subscription tier definition.
type Tier struct {
	AssetContainer
	RotationDays      int16
	RotationRate      int16
	ReportingLevelOID OID
	Paid              bool
	Premium           bool
	Closed            bool
	PricingInfo       string
	ExpiryPeriod      int16
	Ordinal           int32
	ExpiryTierOID     OID
}

func (t *Tier) Serialize(writer gsf.ProtocolWriter) {
	t.AssetContainer.Serialize(writer)
	writer.WriteInt16(t.RotationDays)
	writer.WriteInt16(t.RotationRate)
	writer.WriteObject(&t.ReportingLevelOID)
	writer.WriteBool(t.Paid)
	writer.WriteBool(t.Premium)
	writer.WriteBool(t.Closed)
	writer.WriteString(t.PricingInfo)
	writer.WriteInt16(t.ExpiryPeriod)
	writer.WriteInt32(t.Ordinal)
	writer.WriteObject(&t.ExpiryTierOID)
}

func (t *Tier) Deserialize(reader gsf.ProtocolReader) {
	t.AssetContainer.Deserialize(reader)
	t.RotationDays = reader.ReadInt16()
	t.RotationRate = reader.ReadInt16()
	reader.ReadObject(&t.ReportingLevelOID)
	t.Paid = reader.ReadBool()
	t.Premium = reader.ReadBool()
	t.Closed = reader.ReadBool()
	t.PricingInfo = reader.ReadString()
	t.ExpiryPeriod = reader.ReadInt16()
	t.Ordinal = reader.ReadInt32()
	reader.ReadObject(&t.ExpiryTierOID)
}
