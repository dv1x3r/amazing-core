package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type PlayerHome struct {
	PlayerMaze   PlayerMaze
	PlayerName   string
	Findable     bool
	FindableDate time.Time
	HomeTheme    AssetContainer
	PlayerID     OID
	PlayerMazes  []PlayerMaze
}

func (ph *PlayerHome) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&ph.PlayerMaze)
	writer.WriteString(ph.PlayerName)
	writer.WriteBool(ph.Findable)
	writer.WriteUtcDate(ph.FindableDate)
	writer.WriteObject(&ph.HomeTheme)
	writer.WriteObject(&ph.PlayerID)
	gsf.WriteSlice(writer, ph.PlayerMazes, func(value PlayerMaze) {
		writer.WriteObject(&value)
	})
}

func (ph *PlayerHome) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&ph.PlayerMaze)
	ph.PlayerName = reader.ReadString()
	ph.Findable = reader.ReadBool()
	ph.FindableDate = reader.ReadUtcDate()
	reader.ReadObject(&ph.HomeTheme)
	reader.ReadObject(&ph.PlayerID)
	ph.PlayerMazes = gsf.ReadSlice(reader, func() PlayerMaze {
		var value PlayerMaze
		reader.ReadObject(&value)
		return value
	})
}
