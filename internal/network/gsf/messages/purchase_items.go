package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// PurchaseItemsRequest contains the store products to buy or sell.
type PurchaseItemsRequest struct {
	StoreOID types.OID
	// Items with RequestedQuantity -1 are sold.
	Items []types.StoreItem
}

func (req *PurchaseItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.StoreOID)
	req.Items = gsf.ReadSlice(reader, func() types.StoreItem {
		var value types.StoreItem
		reader.ReadObject(&value)
		return value
	})
}

// PurchaseItemsResponse contains the purchase result and the bought or sold items.
type PurchaseItemsResponse struct {
	// PurchaseResponse is 0 on success, 163 INVALID_PURCHASES, 164 INSUFFICIENT_FUNDS.
	PurchaseResponse            int32
	ReturnedPlayerItems         []types.PlayerItem
	PurchasedPlayerItems        []types.PlayerItem
	ReturnedPlayerBuildObjects  []types.PlayerBuildObject
	PurchasedPlayerBuildObjects []types.PlayerBuildObject
	PurchasedPlayerMazes        []types.PlayerMaze
	PurchasedPlayerAvatars      []types.PlayerAvatar
	ReturnedPlayerAvatars       []types.PlayerAvatar
	ProcessedStoreItems         []types.StoreItem
	AvatarAdoptionItems         []types.PlayerItem
}

func (res *PurchaseItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(res.PurchaseResponse)
	gsf.WriteSlice(writer, res.ReturnedPlayerItems, func(value types.PlayerItem) { writer.WriteObject(&value) })
	gsf.WriteSlice(writer, res.PurchasedPlayerItems, func(value types.PlayerItem) { writer.WriteObject(&value) })
	gsf.WriteSlice(writer, res.ReturnedPlayerBuildObjects, func(value types.PlayerBuildObject) { writer.WriteObject(&value) })
	gsf.WriteSlice(writer, res.PurchasedPlayerBuildObjects, func(value types.PlayerBuildObject) { writer.WriteObject(&value) })
	gsf.WriteSlice(writer, res.PurchasedPlayerMazes, func(value types.PlayerMaze) { writer.WriteObject(&value) })
	gsf.WriteSlice(writer, res.PurchasedPlayerAvatars, func(value types.PlayerAvatar) { writer.WriteObject(&value) })
	gsf.WriteSlice(writer, res.ReturnedPlayerAvatars, func(value types.PlayerAvatar) { writer.WriteObject(&value) })
	gsf.WriteSlice(writer, res.ProcessedStoreItems, func(value types.StoreItem) { writer.WriteObject(&value) })
	gsf.WriteSlice(writer, res.AvatarAdoptionItems, func(value types.PlayerItem) { writer.WriteObject(&value) })
}
