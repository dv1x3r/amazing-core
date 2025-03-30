package messages

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types"
)

type GetOutfitItemsRequest struct {
	PlayerAvatarOutfitID types.OID
	PlayerID             types.OID
}

func (req *GetOutfitItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarOutfitID)
	reader.ReadObject(&req.PlayerID)
}

type GetOutfitItemsResponse struct {
	OutfitItems []types.PlayerItem
}

func (res *GetOutfitItemsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.OutfitItems, func(value types.PlayerItem) {
		writer.WriteObject(&value)
	})
}
