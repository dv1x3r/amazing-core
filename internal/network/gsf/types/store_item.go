package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// StoreItem is a store item instance.
type StoreItem struct {
	RuleContainer
	ObjectOID         OID
	StoreOID          OID
	SKU               string
	Quantity          int64
	PurchaseLimit     int64
	Price             int64
	SalePrice         gsf.Null[int64]
	CurrencyOID       OID
	IsBuyBack         bool
	CreateDate        gsf.UnixTime
	UpdateCount       int32
	RequestedQuantity int64
	RequestedPrice    int64
	RequestedName     string
	PurchaseStatus    int32
	Item              Item
	StoreTheme        StoreTheme
	MazePiece         MazePiece
	Avatar            Avatar
	Currency          Currency
	Subscription      MembershipSubscription
	FeaturedStart     gsf.UnixTime
	FeaturedEnd       gsf.UnixTime
}

func (si *StoreItem) Serialize(writer gsf.ProtocolWriter) {
	si.RuleContainer.Serialize(writer)
	writer.WriteObject(&si.ObjectOID)
	writer.WriteObject(&si.StoreOID)
	writer.WriteString(si.SKU)
	writer.WriteInt64(si.Quantity)
	writer.WriteInt64(si.PurchaseLimit)
	writer.WriteInt64(si.Price)
	gsf.WriteNullable(writer, si.SalePrice, writer.WriteInt64)
	writer.WriteObject(&si.CurrencyOID)
	writer.WriteBool(si.IsBuyBack)
	writer.WriteUtcDate(si.CreateDate)
	writer.WriteInt32(si.UpdateCount)
	writer.WriteInt64(si.RequestedQuantity)
	writer.WriteInt64(si.RequestedPrice)
	writer.WriteString(si.RequestedName)
	writer.WriteInt32(si.PurchaseStatus)
	writer.WriteObject(&si.Item)
	writer.WriteObject(&si.StoreTheme)
	writer.WriteObject(&si.MazePiece)
	writer.WriteObject(&si.Avatar)
	writer.WriteObject(&si.Currency)
	writer.WriteObject(&si.Subscription)
	writer.WriteUtcDate(si.FeaturedStart)
	writer.WriteUtcDate(si.FeaturedEnd)
}

func (si *StoreItem) Deserialize(reader gsf.ProtocolReader) {
	si.RuleContainer.Deserialize(reader)
	reader.ReadObject(&si.ObjectOID)
	reader.ReadObject(&si.StoreOID)
	si.SKU = reader.ReadString()
	si.Quantity = reader.ReadInt64()
	si.PurchaseLimit = reader.ReadInt64()
	si.Price = reader.ReadInt64()
	si.SalePrice = gsf.ReadNullable(reader, reader.ReadInt64)
	reader.ReadObject(&si.CurrencyOID)
	si.IsBuyBack = reader.ReadBool()
	si.CreateDate = reader.ReadUtcDate()
	si.UpdateCount = reader.ReadInt32()
	si.RequestedQuantity = reader.ReadInt64()
	si.RequestedPrice = reader.ReadInt64()
	si.RequestedName = reader.ReadString()
	si.PurchaseStatus = reader.ReadInt32()
	reader.ReadObject(&si.Item)
	reader.ReadObject(&si.StoreTheme)
	reader.ReadObject(&si.MazePiece)
	reader.ReadObject(&si.Avatar)
	reader.ReadObject(&si.Currency)
	reader.ReadObject(&si.Subscription)
	si.FeaturedStart = reader.ReadUtcDate()
	si.FeaturedEnd = reader.ReadUtcDate()
}
