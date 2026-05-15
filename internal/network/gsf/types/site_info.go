package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// SiteInfo contains site account identity data returned during login.
type SiteInfo struct {
	SiteOID       OID
	NicknameFirst string
	NicknameLast  string
	SiteUserOID   OID
}

func (si *SiteInfo) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&si.SiteOID)
	writer.WriteString(si.NicknameFirst)
	writer.WriteString(si.NicknameLast)
	writer.WriteObject(&si.SiteUserOID)
}

func (si *SiteInfo) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&si.SiteOID)
	si.NicknameFirst = reader.ReadString()
	si.NicknameLast = reader.ReadString()
	reader.ReadObject(&si.SiteUserOID)
}
