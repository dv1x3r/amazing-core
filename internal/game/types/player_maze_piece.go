package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type PlayerMazePiece struct {
	OID          OID
	X            float64
	Y            float64
	Z            float64
	Rotation     string
	Ordinal      int16
	ObjectID     OID
	PlayerMazeID OID
}

func (pmp *PlayerMazePiece) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pmp.OID)
	writer.WriteFloat64(pmp.X)
	writer.WriteFloat64(pmp.Y)
	writer.WriteFloat64(pmp.Z)
	writer.WriteString(pmp.Rotation)
	writer.WriteInt16(pmp.Ordinal)
	writer.WriteObject(&pmp.ObjectID)
	writer.WriteObject(&pmp.PlayerMazeID)
}

func (pmp *PlayerMazePiece) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pmp.OID)
	pmp.X = reader.ReadFloat64()
	pmp.Y = reader.ReadFloat64()
	pmp.Z = reader.ReadFloat64()
	pmp.Rotation = reader.ReadString()
	pmp.Ordinal = reader.ReadInt16()
	reader.ReadObject(&pmp.ObjectID)
	reader.ReadObject(&pmp.PlayerMazeID)
}
