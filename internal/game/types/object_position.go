package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type ObjectPosition struct {
	X int32
	Y int32
	Z int32
}

func (op *ObjectPosition) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(op.X)
	writer.WriteInt32(op.Y)
	writer.WriteInt32(op.Z)
}

func (op *ObjectPosition) Deserialize(reader gsf.ProtocolReader) {
	op.X = reader.ReadInt32()
	op.Y = reader.ReadInt32()
	op.Z = reader.ReadInt32()
}
