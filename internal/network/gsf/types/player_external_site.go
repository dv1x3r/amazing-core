package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// PlayerExternalSite is an external account link on a player profile.
type PlayerExternalSite struct {
	OID          OID
	CreateDate   gsf.UnixTime
	PlayerOID    OID
	SiteName     string
	SitePlayerID string
	Token        string
}

func (pes *PlayerExternalSite) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&pes.OID)
	writer.WriteUtcDate(pes.CreateDate)
	writer.WriteObject(&pes.PlayerOID)
	writer.WriteString(pes.SiteName)
	writer.WriteString(pes.SitePlayerID)
	writer.WriteString(pes.Token)
}

func (pes *PlayerExternalSite) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&pes.OID)
	pes.CreateDate = reader.ReadUtcDate()
	reader.ReadObject(&pes.PlayerOID)
	pes.SiteName = reader.ReadString()
	pes.SitePlayerID = reader.ReadString()
	pes.Token = reader.ReadString()
}
