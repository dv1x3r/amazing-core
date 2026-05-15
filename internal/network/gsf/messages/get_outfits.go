package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetOutfitsRequest requests saved outfit presets for a player avatar.
type GetOutfitsRequest struct {
	PlayerAvatarOID types.OID
	PlayerOID       types.OID
}

func (req *GetOutfitsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarOID)
	reader.ReadObject(&req.PlayerOID)
}

// GetOutfitsResponse contains saved outfit presets for a player avatar.
type GetOutfitsResponse struct {
	PlayerAvatarOutfits []types.PlayerAvatarOutfit
}

func (res *GetOutfitsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.PlayerAvatarOutfits, func(value types.PlayerAvatarOutfit) {
		writer.WriteObject(&value)
	})
}
