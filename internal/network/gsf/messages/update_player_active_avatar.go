package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// UpdatePlayerActiveAvatarRequest contains the avatar to make active.
type UpdatePlayerActiveAvatarRequest struct {
	PlayerAvatarOID types.OID
}

func (req *UpdatePlayerActiveAvatarRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarOID)
}

// UpdatePlayerActiveAvatarResponse contains the new active avatar.
type UpdatePlayerActiveAvatarResponse struct {
	ActivePlayerAvatar types.PlayerAvatar
}

func (res *UpdatePlayerActiveAvatarResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.ActivePlayerAvatar)
}
