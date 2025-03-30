package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type SiteInfo struct {
	SiteID        OID
	NicknameFirst string
	NicknameLast  string
	SiteUserID    OID
}

func (si *SiteInfo) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&si.SiteID)
	writer.WriteString(si.NicknameFirst)
	writer.WriteString(si.NicknameLast)
	writer.WriteObject(&si.SiteUserID)
}

func (si *SiteInfo) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&si.SiteID)
	si.NicknameFirst = reader.ReadString()
	si.NicknameLast = reader.ReadString()
	reader.ReadObject(&si.SiteUserID)
}
