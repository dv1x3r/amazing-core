package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// Player is the client-visible player session object.
type Player struct {
	OID                 OID
	CreateDate          gsf.UnixTime
	ActivePlayerAvatar  PlayerAvatar
	HomeThemeOID        OID
	CurrentRaceMode     RaceMode
	WorkshopOptions     string
	IsTutorialCompleted bool
	YardBuildingOID     OID
	LastLogin           gsf.UnixTime
	PlayTime            gsf.Null[int64]
	IsQA                bool
	HomeVillagePlotOID  OID
	StoreVillagePlotOID OID
	PlayerStoreOID      OID
	PlayerMazeOID       OID
	VillageOID          OID
}

func (p *Player) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&p.OID)
	writer.WriteUtcDate(p.CreateDate)
	writer.WriteObject(&p.ActivePlayerAvatar)
	writer.WriteObject(&p.HomeThemeOID)
	writer.WriteObject(&p.CurrentRaceMode)
	writer.WriteString(p.WorkshopOptions)
	writer.WriteBool(p.IsTutorialCompleted)
	writer.WriteObject(&p.YardBuildingOID)
	writer.WriteUtcDate(p.LastLogin)
	gsf.WriteNullable(writer, p.PlayTime, writer.WriteInt64)
	writer.WriteBool(p.IsQA)
	writer.WriteObject(&p.HomeVillagePlotOID)
	writer.WriteObject(&p.StoreVillagePlotOID)
	writer.WriteObject(&p.PlayerStoreOID)
	writer.WriteObject(&p.PlayerMazeOID)
	writer.WriteObject(&p.VillageOID)
}

func (p *Player) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&p.OID)
	p.CreateDate = reader.ReadUtcDate()
	reader.ReadObject(&p.ActivePlayerAvatar)
	reader.ReadObject(&p.HomeThemeOID)
	reader.ReadObject(&p.CurrentRaceMode)
	p.WorkshopOptions = reader.ReadString()
	p.IsTutorialCompleted = reader.ReadBool()
	reader.ReadObject(&p.YardBuildingOID)
	p.LastLogin = reader.ReadUtcDate()
	p.PlayTime = gsf.ReadNullable(reader, reader.ReadInt64)
	p.IsQA = reader.ReadBool()
	reader.ReadObject(&p.HomeVillagePlotOID)
	reader.ReadObject(&p.StoreVillagePlotOID)
	reader.ReadObject(&p.PlayerStoreOID)
	reader.ReadObject(&p.PlayerMazeOID)
	reader.ReadObject(&p.VillageOID)
}
