package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// PurchaseItemsRequest buys the products of a store, RequestedQuantity and RequestedPrice are set by the client.
type PurchaseItemsRequest struct {
	StoreOID types.OID
	Items    []types.StoreItem
}

func (req *PurchaseItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.StoreOID)
	req.Items = gsf.ReadSlice(reader, func() types.StoreItem {
		var value types.StoreItem
		reader.ReadObject(&value)
		return value
	})
}

// PurchaseItemsResponse contains the purchase result, PurchaseResponse 0 is a success.
type PurchaseItemsResponse struct {
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
