package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// AddOutfitRequest contains the avatar and the number of the new outfit.
type AddOutfitRequest struct {
	PlayerAvatarOID types.OID
	OutfitNo        int16
}

func (req *AddOutfitRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarOID)
	req.OutfitNo = reader.ReadInt16()
}

// AddOutfitResponse contains the new outfit and its number.
type AddOutfitResponse struct {
	PlayerAvatarOutfitOID types.OID
	OutfitNo              int16
}

func (res *AddOutfitResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.PlayerAvatarOutfitOID)
	writer.WriteInt16(res.OutfitNo)
}
