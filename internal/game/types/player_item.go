package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type PlayerItem struct {
	OID                     OID
	Item                    Item
	SecretCode              string
	Ordinal                 int32
	ParentPioID             OID
	SlotID                  OID
	PlayerAvatarOutfitID    OID
	InventoryPosition       InventoryPosition
	PlayerMazePiecesID      OID
	IsYard                  bool
	PlayerMazeID            OID
	PlayerAvatarID          OID
	PlayerID                OID
	IsItemUsed              bool
	PlayerContainerID       OID
	PlacedPlayerContainerID OID
	SellPrice               int32
	StoreThemeID            OID
	CreateDate              time.Time
	GrowthCompletionDate    time.Time
	GrowthStartDate         time.Time
	MatureEndDate           time.Time
	DecayEndDate            time.Time
	HarvestDate             time.Time
	AttachedItems           []PlayerItem
	SendingID               OID
	Quantity                int32
	UnitsToExpire           int32
	QualityIndex            int32
}

func (pi *PlayerItem) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pi.OID)
	writer.WriteObject(&pi.Item)
	writer.WriteString(pi.SecretCode)
	writer.WriteInt32(pi.Ordinal)
	writer.WriteObject(&pi.ParentPioID)
	writer.WriteObject(&pi.SlotID)
	writer.WriteObject(&pi.PlayerAvatarOutfitID)
	writer.WriteObject(&pi.InventoryPosition)
	writer.WriteObject(&pi.PlayerMazePiecesID)
	writer.WriteBool(pi.IsYard)
	writer.WriteObject(&pi.PlayerMazeID)
	writer.WriteObject(&pi.PlayerAvatarID)
	writer.WriteObject(&pi.PlayerID)
	writer.WriteBool(pi.IsItemUsed)
	writer.WriteObject(&pi.PlayerContainerID)
	writer.WriteObject(&pi.PlacedPlayerContainerID)
	writer.WriteInt32(pi.SellPrice)
	writer.WriteObject(&pi.StoreThemeID)
	writer.WriteUtcDate(pi.CreateDate)
	writer.WriteUtcDate(pi.GrowthCompletionDate)
	writer.WriteUtcDate(pi.GrowthStartDate)
	writer.WriteUtcDate(pi.MatureEndDate)
	writer.WriteUtcDate(pi.DecayEndDate)
	writer.WriteUtcDate(pi.HarvestDate)
	gsf.WriteSlice(writer, pi.AttachedItems, func(value PlayerItem) {
		writer.WriteObject(&value)
	})
	writer.WriteObject(&pi.SendingID)
	writer.WriteInt32(pi.Quantity)
	writer.WriteInt32(pi.UnitsToExpire)
	writer.WriteInt32(pi.QualityIndex)
}

func (pi *PlayerItem) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pi.OID)
	reader.ReadObject(&pi.Item)
	pi.SecretCode = reader.ReadString()
	pi.Ordinal = reader.ReadInt32()
	reader.ReadObject(&pi.ParentPioID)
	reader.ReadObject(&pi.SlotID)
	reader.ReadObject(&pi.PlayerAvatarOutfitID)
	reader.ReadObject(&pi.InventoryPosition)
	reader.ReadObject(&pi.PlayerMazePiecesID)
	pi.IsYard = reader.ReadBool()
	reader.ReadObject(&pi.PlayerMazeID)
	reader.ReadObject(&pi.PlayerAvatarID)
	reader.ReadObject(&pi.PlayerID)
	pi.IsItemUsed = reader.ReadBool()
	reader.ReadObject(&pi.PlayerContainerID)
	reader.ReadObject(&pi.PlacedPlayerContainerID)
	pi.SellPrice = reader.ReadInt32()
	reader.ReadObject(&pi.StoreThemeID)
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
	reader.ReadObject(&pi.SendingID)
	pi.Quantity = reader.ReadInt32()
	pi.UnitsToExpire = reader.ReadInt32()
	pi.QualityIndex = reader.ReadInt32()
}
