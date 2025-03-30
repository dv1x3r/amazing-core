package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type Player struct {
	OID                 OID
	CreateDate          time.Time
	ActivePlayerAvatar  PlayerAvatar
	HomeThemeID         OID
	CurrentRaceMode     RaceMode
	WorkshopOptions     string
	IsTutorialCompleted bool
	YardBuildingID      OID
	LastLogin           time.Time
	PlayTime            gsf.Null[int64]
	IsQA                bool
	HomeVillagePlotID   OID
	StoreVillagePlotID  OID
	PlayerStoreID       OID
	PlayerMazeID        OID
	VillageID           OID
}

func (p *Player) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&p.OID)
	writer.WriteUtcDate(p.CreateDate)
	writer.WriteObject(&p.ActivePlayerAvatar)
	writer.WriteObject(&p.HomeThemeID)
	writer.WriteObject(&p.CurrentRaceMode)
	writer.WriteString(p.WorkshopOptions)
	writer.WriteBool(p.IsTutorialCompleted)
	writer.WriteObject(&p.YardBuildingID)
	writer.WriteUtcDate(p.LastLogin)
	gsf.WriteNullable(writer, p.PlayTime, writer.WriteInt64)
	writer.WriteBool(p.IsQA)
	writer.WriteObject(&p.HomeVillagePlotID)
	writer.WriteObject(&p.StoreVillagePlotID)
	writer.WriteObject(&p.PlayerStoreID)
	writer.WriteObject(&p.PlayerMazeID)
	writer.WriteObject(&p.VillageID)
}

func (p *Player) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&p.OID)
	p.CreateDate = reader.ReadUtcDate()
	reader.ReadObject(&p.ActivePlayerAvatar)
	reader.ReadObject(&p.HomeThemeID)
	reader.ReadObject(&p.CurrentRaceMode)
	p.WorkshopOptions = reader.ReadString()
	p.IsTutorialCompleted = reader.ReadBool()
	reader.ReadObject(&p.YardBuildingID)
	p.LastLogin = reader.ReadUtcDate()
	p.PlayTime = gsf.ReadNullable(reader, reader.ReadInt64)
	p.IsQA = reader.ReadBool()
	reader.ReadObject(&p.HomeVillagePlotID)
	reader.ReadObject(&p.StoreVillagePlotID)
	reader.ReadObject(&p.PlayerStoreID)
	reader.ReadObject(&p.PlayerMazeID)
	reader.ReadObject(&p.VillageID)
}
