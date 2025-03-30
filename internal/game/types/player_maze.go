package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type PlayerMaze struct {
	OID              OID
	Name             string
	Size             int64
	Thumbnail        []byte
	PublishTimestamp time.Time
	NumRooms         int16
	NumTubes         int16
	Rating           gsf.Null[int16]
	IsLocked         bool
	IsHomeMaze       bool
	IsPublished      bool
	IsPublishExpired bool
	PlayerID         OID
	MazePieces       []PlayerMazePiece
	HomeTheme        AssetContainer
	ParentID         OID
	SourceID         OID
}

func (pm *PlayerMaze) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pm.OID)
	writer.WriteString(pm.Name)
	writer.WriteInt64(pm.Size)
	writer.WriteBytes(pm.Thumbnail)
	writer.WriteUtcDate(pm.PublishTimestamp)
	writer.WriteInt16(pm.NumRooms)
	writer.WriteInt16(pm.NumTubes)
	gsf.WriteNullable(writer, pm.Rating, writer.WriteInt16)
	writer.WriteBool(pm.IsLocked)
	writer.WriteBool(pm.IsHomeMaze)
	writer.WriteBool(pm.IsPublished)
	writer.WriteBool(pm.IsPublishExpired)
	writer.WriteObject(&pm.PlayerID)
	gsf.WriteSlice(writer, pm.MazePieces, func(value PlayerMazePiece) {
		writer.WriteObject(&value)
	})
	writer.WriteObject(&pm.HomeTheme)
	writer.WriteObject(&pm.ParentID)
	writer.WriteObject(&pm.SourceID)
}

func (pm *PlayerMaze) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pm.OID)
	pm.Name = reader.ReadString()
	pm.Size = reader.ReadInt64()
	pm.Thumbnail = reader.ReadBytes()
	pm.PublishTimestamp = reader.ReadUtcDate()
	pm.NumRooms = reader.ReadInt16()
	pm.NumTubes = reader.ReadInt16()
	pm.Rating = gsf.ReadNullable(reader, reader.ReadInt16)
	pm.IsLocked = reader.ReadBool()
	pm.IsHomeMaze = reader.ReadBool()
	pm.IsPublished = reader.ReadBool()
	pm.IsPublishExpired = reader.ReadBool()
	reader.ReadObject(&pm.PlayerID)
	pm.MazePieces = gsf.ReadSlice(reader, func() PlayerMazePiece {
		var value PlayerMazePiece
		reader.ReadObject(&value)
		return value
	})
	reader.ReadObject(&pm.HomeTheme)
	reader.ReadObject(&pm.ParentID)
	reader.ReadObject(&pm.SourceID)
}
