package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type Zone struct {
	RuleContainer
	Dimensions Dimensions
	Buildings  []Building
	PTag       string
	Capacity   int32
}

func (z *Zone) Serialize(writer gsf.ProtocolWriter) {
	z.RuleContainer.Serialize(writer)
	writer.WriteObject(&z.Dimensions)
	gsf.WriteSlice(writer, z.Buildings, func(value Building) {
		writer.WriteObject(&value)
	})
	writer.WriteString(z.PTag)
	writer.WriteInt32(z.Capacity)
}

func (z *Zone) Deserialize(reader gsf.ProtocolReader) {
	z.RuleContainer.Deserialize(reader)
	reader.ReadObject(&z.Dimensions)
	z.Buildings = gsf.ReadSlice(reader, func() Building {
		var value Building
		reader.ReadObject(&value)
		return value
	})
	z.PTag = reader.ReadString()
	z.Capacity = reader.ReadInt32()
}
