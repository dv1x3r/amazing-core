package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

type Currency struct {
	AssetContainer
	StatsTypeOID OID
	IsDefault    gsf.Null[bool]
}

func (c *Currency) Serialize(writer gsf.ProtocolWriter) {
	c.AssetContainer.Serialize(writer)
	writer.WriteObject(&c.StatsTypeOID)
	gsf.WriteNullable(writer, c.IsDefault, writer.WriteBool)
}

func (c *Currency) Deserialize(reader gsf.ProtocolReader) {
	c.AssetContainer.Deserialize(reader)
	reader.ReadObject(&c.StatsTypeOID)
	c.IsDefault = gsf.ReadNullable(reader, reader.ReadBool)
}
