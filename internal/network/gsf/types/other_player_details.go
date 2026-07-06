package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// OtherPlayerDetails is the profile payload used to spawn a remote sync player.
type OtherPlayerDetails struct {
	PlayerAvatar      PlayerAvatar
	Clothing          []PlayerItem
	PlayerName        string
	WorldName         string
	TierOID           OID
	PlayerAvatarCount int32
	Level             int32
	XP                int32
	Token             int32
	Energy            int32
	PlayerFriendCount int32
	Findable          bool
	FindableDuration  int32
	ExternalSites     []PlayerExternalSite
}

func (opd *OtherPlayerDetails) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&opd.PlayerAvatar)
	gsf.WriteSlice(writer, opd.Clothing, func(value PlayerItem) {
		writer.WriteObject(&value)
	})
	writer.WriteString(opd.PlayerName)
	writer.WriteString(opd.WorldName)
	writer.WriteObject(&opd.TierOID)
	writer.WriteInt32(opd.PlayerAvatarCount)
	writer.WriteInt32(opd.Level)
	writer.WriteInt32(opd.XP)
	writer.WriteInt32(opd.Token)
	writer.WriteInt32(opd.Energy)
	writer.WriteInt32(opd.PlayerFriendCount)
	writer.WriteBool(opd.Findable)
	writer.WriteInt32(opd.FindableDuration)
	gsf.WriteSlice(writer, opd.ExternalSites, func(value PlayerExternalSite) {
		writer.WriteObject(&value)
	})
}

func (opd *OtherPlayerDetails) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&opd.PlayerAvatar)
	opd.Clothing = gsf.ReadSlice(reader, func() PlayerItem {
		var value PlayerItem
		reader.ReadObject(&value)
		return value
	})
	opd.PlayerName = reader.ReadString()
	opd.WorldName = reader.ReadString()
	reader.ReadObject(&opd.TierOID)
	opd.PlayerAvatarCount = reader.ReadInt32()
	opd.Level = reader.ReadInt32()
	opd.XP = reader.ReadInt32()
	opd.Token = reader.ReadInt32()
	opd.Energy = reader.ReadInt32()
	opd.PlayerFriendCount = reader.ReadInt32()
	opd.Findable = reader.ReadBool()
	opd.FindableDuration = reader.ReadInt32()
	opd.ExternalSites = gsf.ReadSlice(reader, func() PlayerExternalSite {
		var value PlayerExternalSite
		reader.ReadObject(&value)
		return value
	})
}
