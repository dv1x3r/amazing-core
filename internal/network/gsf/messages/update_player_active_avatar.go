package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// UpdatePlayerActiveAvatarRequest selects the player's active avatar.
type UpdatePlayerActiveAvatarRequest struct {
	PlayerAvatarID types.OID
}

func (req *UpdatePlayerActiveAvatarRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarID)
}

// UpdatePlayerActiveAvatarResponse contains the newly active player avatar.
type UpdatePlayerActiveAvatarResponse struct {
	ActivePlayerAvatar types.PlayerAvatar
}

func (res *UpdatePlayerActiveAvatarResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.ActivePlayerAvatar)
}
