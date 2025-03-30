package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type Dimensions struct {
	CX  int32
	CY  int32
	CZ  int32
	BOX int32
	BOY int32
	BOZ int32
}

func (d *Dimensions) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(d.CX)
	writer.WriteInt32(d.CY)
	writer.WriteInt32(d.CZ)
	writer.WriteInt32(d.BOX)
	writer.WriteInt32(d.BOY)
	writer.WriteInt32(d.BOZ)
}

func (d *Dimensions) Deserialize(reader gsf.ProtocolReader) {
	d.CX = reader.ReadInt32()
	d.CY = reader.ReadInt32()
	d.CZ = reader.ReadInt32()
	d.BOX = reader.ReadInt32()
	d.BOY = reader.ReadInt32()
	d.BOZ = reader.ReadInt32()
}
