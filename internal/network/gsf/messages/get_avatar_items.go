package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetAvatarItemsRequest requests item instances owned by a player avatar.
type GetAvatarItemsRequest struct {
	PlayerAvatarID types.OID
	PlayerID       types.OID
}

func (req *GetAvatarItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarID)
	reader.ReadObject(&req.PlayerID)
}

// GetAvatarItemsResponse contains item instances owned by a player avatar.
type GetAvatarItemsResponse struct {
	AvatarItems []types.PlayerItem
}

func (res *GetAvatarItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.AvatarItems, func(value types.PlayerItem) {
		writer.WriteObject(&value)
	})
}
