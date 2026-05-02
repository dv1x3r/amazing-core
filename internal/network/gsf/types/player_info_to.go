package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// PlayerInfoTO contains account and policy data returned during login.
type PlayerInfoTO struct {
	TierID                     OID
	PlayerName                 string
	WorldName                  string
	CrispDataTO                CrispDataTO
	Verified                   bool
	VerificationExpiry         gsf.UnixTime
	ScsBlockExpiry             gsf.UnixTime
	EulaID                     OID
	CurrentEulaID              OID
	U13                        bool
	ChatBlockedParent          bool
	ChatAllowed                bool
	ChatBlockedExpiry          gsf.UnixTime
	Findable                   bool
	FindableDate               gsf.UnixTime
	UserSubscriptionExpiryDate gsf.UnixTime
	QA                         bool
	PlayerSettings             []PlayerSetting
}

func (pi *PlayerInfoTO) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pi.TierID)
	writer.WriteString(pi.PlayerName)
	writer.WriteString(pi.WorldName)
	writer.WriteObject(&pi.CrispDataTO)
	writer.WriteBool(pi.Verified)
	writer.WriteUtcDate(pi.VerificationExpiry)
	writer.WriteUtcDate(pi.ScsBlockExpiry)
	writer.WriteObject(&pi.EulaID)
	writer.WriteObject(&pi.CurrentEulaID)
	writer.WriteBool(pi.U13)
	writer.WriteBool(pi.ChatBlockedParent)
	writer.WriteBool(pi.ChatAllowed)
	writer.WriteUtcDate(pi.ChatBlockedExpiry)
	writer.WriteBool(pi.Findable)
	writer.WriteUtcDate(pi.FindableDate)
	writer.WriteUtcDate(pi.UserSubscriptionExpiryDate)
	writer.WriteBool(pi.QA)
	gsf.WriteSlice(writer, pi.PlayerSettings, func(value PlayerSetting) {
		writer.WriteObject(&value)
	})
}

func (pi *PlayerInfoTO) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pi.TierID)
	pi.PlayerName = reader.ReadString()
	pi.WorldName = reader.ReadString()
	reader.ReadObject(&pi.CrispDataTO)
	pi.Verified = reader.ReadBool()
	pi.VerificationExpiry = reader.ReadUtcDate()
	pi.ScsBlockExpiry = reader.ReadUtcDate()
	reader.ReadObject(&pi.EulaID)
	reader.ReadObject(&pi.CurrentEulaID)
	pi.U13 = reader.ReadBool()
	pi.ChatBlockedParent = reader.ReadBool()
	pi.ChatAllowed = reader.ReadBool()
	pi.ChatBlockedExpiry = reader.ReadUtcDate()
	pi.Findable = reader.ReadBool()
	pi.FindableDate = reader.ReadUtcDate()
	pi.UserSubscriptionExpiryDate = reader.ReadUtcDate()
	pi.QA = reader.ReadBool()
	pi.PlayerSettings = gsf.ReadSlice(reader, func() PlayerSetting {
		var value PlayerSetting
		reader.ReadObject(&value)
		return value
	})
}
