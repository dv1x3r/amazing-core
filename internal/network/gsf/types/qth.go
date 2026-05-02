package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// QTH is a compact orientation tuple used by GSF movement data.
type QTH struct {
	SW bool
	SX bool
	SY bool
	SZ bool
	CX rune
	CY rune
	CZ rune
}

func (qth *QTH) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(qth.SW)
	writer.WriteBool(qth.SX)
	writer.WriteBool(qth.SY)
	writer.WriteBool(qth.SZ)
	writer.WriteChar(qth.CX)
	writer.WriteChar(qth.CY)
	writer.WriteChar(qth.CZ)
}

func (qth *QTH) Deserialize(reader gsf.ProtocolReader) {
	qth.SW = reader.ReadBool()
	qth.SX = reader.ReadBool()
	qth.SY = reader.ReadBool()
	qth.SZ = reader.ReadBool()
	qth.CX = reader.ReadChar()
	qth.CY = reader.ReadChar()
	qth.CZ = reader.ReadChar()
}
