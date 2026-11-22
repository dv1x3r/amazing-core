package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetAvatarItemsRequest contains the player avatar of the items.
type GetAvatarItemsRequest struct {
	PlayerAvatarOID types.OID
	PlayerOID       types.OID
}

func (req *GetAvatarItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarOID)
	reader.ReadObject(&req.PlayerOID)
}

// GetAvatarItemsResponse contains the items of the player avatar.
type GetAvatarItemsResponse struct {
	AvatarItems []types.PlayerItem
}

func (res *GetAvatarItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.AvatarItems, func(value types.PlayerItem) {
		writer.WriteObject(&value)
	})
}
