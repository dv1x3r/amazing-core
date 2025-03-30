package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type Item struct {
	AssetContainer
	CreateDate         time.Time
	Depth              int32
	Height             int32
	Width              int32
	IsConsumable       bool
	IsPresentable      bool
	IsTradeable        bool
	IsLighting         bool
	IsAnimated         bool
	AcceptsPresentable bool
	InContainer        bool
	PresentableSlots   string
	QualityIndex       int32
	SellPrice          int64
	BuyPrice           int64
	Name               string
	IsUserSellable     bool
	GrowthRate         int32
	SpawnPoint         string
	MatureDuration     int32
	DecayDuration      int32
	Quantity           int32
	ItemCategories     []ItemCategory
	AcceptableSlotIds  []OID
}

func (i *Item) Serialize(writer gsf.ProtocolWriter) {
	i.AssetContainer.Serialize(writer)
	writer.WriteUtcDate(i.CreateDate)
	writer.WriteInt32(i.Depth)
	writer.WriteInt32(i.Height)
	writer.WriteInt32(i.Width)
	writer.WriteBool(i.IsConsumable)
	writer.WriteBool(i.IsPresentable)
	writer.WriteBool(i.IsTradeable)
	writer.WriteBool(i.IsLighting)
	writer.WriteBool(i.IsAnimated)
	writer.WriteBool(i.AcceptsPresentable)
	writer.WriteBool(i.InContainer)
	writer.WriteString(i.PresentableSlots)
	writer.WriteInt32(i.QualityIndex)
	writer.WriteInt64(i.SellPrice)
	writer.WriteInt64(i.BuyPrice)
	writer.WriteString(i.Name)
	writer.WriteBool(i.IsUserSellable)
	writer.WriteInt32(i.GrowthRate)
	writer.WriteString(i.SpawnPoint)
	writer.WriteInt32(i.MatureDuration)
	writer.WriteInt32(i.DecayDuration)
	writer.WriteInt32(i.Quantity)
	gsf.WriteSlice(writer, i.ItemCategories, func(value ItemCategory) {
		writer.WriteObject(&value)
	})
	gsf.WriteSlice(writer, i.AcceptableSlotIds, func(value OID) {
		writer.WriteObject(&value)
	})
}

func (i *Item) Deserialize(reader gsf.ProtocolReader) {
	i.AssetContainer.Deserialize(reader)
	i.CreateDate = reader.ReadUtcDate()
	i.Depth = reader.ReadInt32()
	i.Height = reader.ReadInt32()
	i.Width = reader.ReadInt32()
	i.IsConsumable = reader.ReadBool()
	i.IsPresentable = reader.ReadBool()
	i.IsTradeable = reader.ReadBool()
	i.IsLighting = reader.ReadBool()
	i.IsAnimated = reader.ReadBool()
	i.AcceptsPresentable = reader.ReadBool()
	i.InContainer = reader.ReadBool()
	i.PresentableSlots = reader.ReadString()
	i.QualityIndex = reader.ReadInt32()
	i.SellPrice = reader.ReadInt64()
	i.BuyPrice = reader.ReadInt64()
	i.Name = reader.ReadString()
	i.IsUserSellable = reader.ReadBool()
	i.GrowthRate = reader.ReadInt32()
	i.SpawnPoint = reader.ReadString()
	i.MatureDuration = reader.ReadInt32()
	i.DecayDuration = reader.ReadInt32()
	i.Quantity = reader.ReadInt32()
	i.ItemCategories = gsf.ReadSlice(reader, func() ItemCategory {
		var value ItemCategory
		reader.ReadObject(&value)
		return value
	})
	i.AcceptableSlotIds = gsf.ReadSlice(reader, func() OID {
		var value OID
		reader.ReadObject(&value)
		return value
	})
}
