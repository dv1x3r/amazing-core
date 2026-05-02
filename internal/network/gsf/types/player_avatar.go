package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// PlayerAvatar is a player-owned avatar record with avatar definition and play state.
type PlayerAvatar struct {
	OID                  OID
	Avatar               Avatar
	PlayerID             OID
	Name                 string
	Bio                  string
	SecretCode           string
	CreateTS             gsf.UnixTime
	PlayerAvatarOutfitID OID
	OutfitNo             int16
	PlayTime             gsf.Null[int64]
	LastPlay             gsf.UnixTime
}

func (pa *PlayerAvatar) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pa.OID)
	writer.WriteObject(&pa.Avatar)
	writer.WriteObject(&pa.PlayerID)
	writer.WriteString(pa.Name)
	writer.WriteString(pa.Bio)
	writer.WriteString(pa.SecretCode)
	writer.WriteUtcDate(pa.CreateTS)
	writer.WriteObject(&pa.PlayerAvatarOutfitID)
	writer.WriteInt16(pa.OutfitNo)
	gsf.WriteNullable(writer, pa.PlayTime, writer.WriteInt64)
	writer.WriteUtcDate(pa.LastPlay)
}

func (pa *PlayerAvatar) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pa.OID)
	reader.ReadObject(&pa.Avatar)
	reader.ReadObject(&pa.PlayerID)
	pa.Name = reader.ReadString()
	pa.Bio = reader.ReadString()
	pa.SecretCode = reader.ReadString()
	pa.CreateTS = reader.ReadUtcDate()
	reader.ReadObject(&pa.PlayerAvatarOutfitID)
	pa.OutfitNo = reader.ReadInt16()
	pa.PlayTime = gsf.ReadNullable(reader, reader.ReadInt64)
	pa.LastPlay = reader.ReadUtcDate()
}
