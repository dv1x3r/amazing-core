package messages

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types"
)

type GetOutfitsRequest struct {
	PlayerAvatarID types.OID
	PlayerID       types.OID
}

func (req *GetOutfitsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarID)
	reader.ReadObject(&req.PlayerID)
}

type GetOutfitsResponse struct {
	PlayerAvatarOutfits []types.PlayerAvatarOutfit
}

func (res *GetOutfitsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.PlayerAvatarOutfits, func(value types.PlayerAvatarOutfit) {
		writer.WriteObject(&value)
	})
}
