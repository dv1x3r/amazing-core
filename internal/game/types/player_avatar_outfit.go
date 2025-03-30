package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type PlayerAvatarOutfit struct {
	OID            OID
	PlayerID       OID
	PlayerAvatarID OID
	OutfitNo       int16
}

func (ps *PlayerAvatarOutfit) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&ps.OID)
	writer.WriteObject(&ps.PlayerID)
	writer.WriteObject(&ps.PlayerAvatarID)
	writer.WriteInt16(ps.OutfitNo)
}

func (ps *PlayerAvatarOutfit) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&ps.OID)
	reader.ReadObject(&ps.PlayerID)
	reader.ReadObject(&ps.PlayerAvatarID)
	ps.OutfitNo = reader.ReadInt16()
}
