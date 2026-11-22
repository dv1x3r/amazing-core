package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// OIDLongPair is an OID with an optional number.
type OIDLongPair struct {
	A OID
	B gsf.Null[int64]
}

func (p *OIDLongPair) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&p.A)
	gsf.WriteNullable(writer, p.B, writer.WriteInt64)
}

func (p *OIDLongPair) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&p.A)
	p.B = gsf.ReadNullable(reader, reader.ReadInt64)
}
