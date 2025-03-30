package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type PlayerAvatar struct {
	OID                  OID
	Avatar               Avatar
	PlayerID             OID
	Name                 string
	Bio                  string
	SecretCode           string
	CreateTS             time.Time
	PlayerAvatarOutfitID OID
	OutfitNo             int16
	PlayTime             gsf.Null[int64]
	LastPlay             time.Time
}

func (ap *PlayerAvatar) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&ap.OID)
	writer.WriteObject(&ap.Avatar)
	writer.WriteObject(&ap.PlayerID)
	writer.WriteString(ap.Name)
	writer.WriteString(ap.Bio)
	writer.WriteString(ap.SecretCode)
	writer.WriteUtcDate(ap.CreateTS)
	writer.WriteObject(&ap.PlayerAvatarOutfitID)
	writer.WriteInt16(ap.OutfitNo)
	gsf.WriteNullable(writer, ap.PlayTime, writer.WriteInt64)
	writer.WriteUtcDate(ap.LastPlay)
}

func (ap *PlayerAvatar) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&ap.OID)
	reader.ReadObject(&ap.Avatar)
	reader.ReadObject(&ap.PlayerID)
	ap.Name = reader.ReadString()
	ap.Bio = reader.ReadString()
	ap.SecretCode = reader.ReadString()
	ap.CreateTS = reader.ReadUtcDate()
	reader.ReadObject(&ap.PlayerAvatarOutfitID)
	ap.OutfitNo = reader.ReadInt16()
	ap.PlayTime = gsf.ReadNullable(reader, reader.ReadInt64)
	ap.LastPlay = reader.ReadUtcDate()
}
