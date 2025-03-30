package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type Position struct {
	X int32
	Y int32
	Z int32
	T int32
}

func (p *Position) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(p.X)
	writer.WriteInt32(p.Y)
	writer.WriteInt32(p.Z)
	writer.WriteInt32(p.T)
}

func (p *Position) Deserialize(reader gsf.ProtocolReader) {
	p.X = reader.ReadInt32()
	p.Y = reader.ReadInt32()
	p.Z = reader.ReadInt32()
	p.T = reader.ReadInt32()
}
