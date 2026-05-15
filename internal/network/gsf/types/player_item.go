package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// PlayerItem is a player-owned item instance.
type PlayerItem struct {
	OID                      OID
	Item                     Item
	SecretCode               string
	Ordinal                  int32
	ParentPioOID             OID
	SlotOID                  OID
	PlayerAvatarOutfitOID    OID
	InventoryPosition        InventoryPosition
	PlayerMazePiecesOID      OID
	IsYard                   bool
	PlayerMazeOID            OID
	PlayerAvatarOID          OID
	PlayerOID                OID
	IsItemUsed               bool
	PlayerContainerOID       OID
	PlacedPlayerContainerOID OID
	SellPrice                int32
	StoreThemeOID            OID
	CreateDate               gsf.UnixTime
	GrowthCompletionDate     gsf.UnixTime
	GrowthStartDate          gsf.UnixTime
	MatureEndDate            gsf.UnixTime
	DecayEndDate             gsf.UnixTime
	HarvestDate              gsf.UnixTime
	AttachedItems            []PlayerItem
	SendingOID               OID
	Quantity                 int32
	UnitsToExpire            int32
	QualityIndex             int32
}

func (pi *PlayerItem) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pi.OID)
	writer.WriteObject(&pi.Item)
	writer.WriteString(pi.SecretCode)
	writer.WriteInt32(pi.Ordinal)
	writer.WriteObject(&pi.ParentPioOID)
	writer.WriteObject(&pi.SlotOID)
	writer.WriteObject(&pi.PlayerAvatarOutfitOID)
	writer.WriteObject(&pi.InventoryPosition)
	writer.WriteObject(&pi.PlayerMazePiecesOID)
	writer.WriteBool(pi.IsYard)
	writer.WriteObject(&pi.PlayerMazeOID)
	writer.WriteObject(&pi.PlayerAvatarOID)
	writer.WriteObject(&pi.PlayerOID)
	writer.WriteBool(pi.IsItemUsed)
	writer.WriteObject(&pi.PlayerContainerOID)
	writer.WriteObject(&pi.PlacedPlayerContainerOID)
	writer.WriteInt32(pi.SellPrice)
	writer.WriteObject(&pi.StoreThemeOID)
	writer.WriteUtcDate(pi.CreateDate)
	writer.WriteUtcDate(pi.GrowthCompletionDate)
	writer.WriteUtcDate(pi.GrowthStartDate)
	writer.WriteUtcDate(pi.MatureEndDate)
	writer.WriteUtcDate(pi.DecayEndDate)
	writer.WriteUtcDate(pi.HarvestDate)
	gsf.WriteSlice(writer, pi.AttachedItems, func(value PlayerItem) {
		writer.WriteObject(&value)
	})
	writer.WriteObject(&pi.SendingOID)
	writer.WriteInt32(pi.Quantity)
	writer.WriteInt32(pi.UnitsToExpire)
	writer.WriteInt32(pi.QualityIndex)
}

func (pi *PlayerItem) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pi.OID)
	reader.ReadObject(&pi.Item)
	pi.SecretCode = reader.ReadString()
	pi.Ordinal = reader.ReadInt32()
	reader.ReadObject(&pi.ParentPioOID)
	reader.ReadObject(&pi.SlotOID)
	reader.ReadObject(&pi.PlayerAvatarOutfitOID)
	reader.ReadObject(&pi.InventoryPosition)
	reader.ReadObject(&pi.PlayerMazePiecesOID)
	pi.IsYard = reader.ReadBool()
	reader.ReadObject(&pi.PlayerMazeOID)
	reader.ReadObject(&pi.PlayerAvatarOID)
	reader.ReadObject(&pi.PlayerOID)
	pi.IsItemUsed = reader.ReadBool()
	reader.ReadObject(&pi.PlayerContainerOID)
	reader.ReadObject(&pi.PlacedPlayerContainerOID)
	pi.SellPrice = reader.ReadInt32()
	reader.ReadObject(&pi.StoreThemeOID)
	pi.CreateDate = reader.ReadUtcDate()
	pi.GrowthCompletionDate = reader.ReadUtcDate()
	pi.GrowthStartDate = reader.ReadUtcDate()
	pi.MatureEndDate = reader.ReadUtcDate()
	pi.DecayEndDate = reader.ReadUtcDate()
	pi.HarvestDate = reader.ReadUtcDate()
	pi.AttachedItems = gsf.ReadSlice(reader, func() PlayerItem {
		var value PlayerItem
		reader.ReadObject(&value)
		return value
	})
	reader.ReadObject(&pi.SendingOID)
	pi.Quantity = reader.ReadInt32()
	pi.UnitsToExpire = reader.ReadInt32()
	pi.QualityIndex = reader.ReadInt32()
}
