package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// AddOutfitRequest requests a new player outfit creation.
type AddOutfitRequest struct {
	PlayerAvatarOID types.OID
	OutfitNo        int16
}

func (req *AddOutfitRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerAvatarOID)
	req.OutfitNo = reader.ReadInt16()
}

// AddOutfitResponse contains OID for the new outfit instance.
type AddOutfitResponse struct {
	PlayerAvatarOutfitOID types.OID
	OutfitNo              int16
}

func (res *AddOutfitResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.PlayerAvatarOutfitOID)
	writer.WriteInt16(res.OutfitNo)
}
