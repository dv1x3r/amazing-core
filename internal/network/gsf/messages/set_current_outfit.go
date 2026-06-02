package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// SetCurrentOutfitRequest selects the player's active outfit.
type SetCurrentOutfitRequest struct {
	PlayerAvatarOID       types.OID
	PlayerAvatarOutfitOID types.OID
	OutfitNo              int16
}

func (req *SetCurrentOutfitRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarOID)
	reader.ReadObject(&req.PlayerAvatarOutfitOID)
	req.OutfitNo = reader.ReadInt16()
}

// SetCurrentOutfitResponse contains the update status.
type SetCurrentOutfitResponse struct {
	IsUpdated bool
}

func (res *SetCurrentOutfitResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteBool(res.IsUpdated)
}
