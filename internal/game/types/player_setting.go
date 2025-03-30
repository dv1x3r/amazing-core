package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type PlayerSetting struct {
	OID         OID
	PlayerID    OID
	SettingName string
	Value       string
}

func (ps *PlayerSetting) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&ps.OID)
	writer.WriteObject(&ps.PlayerID)
	writer.WriteString(ps.SettingName)
	writer.WriteString(ps.Value)
}

func (ps *PlayerSetting) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&ps.OID)
	reader.ReadObject(&ps.PlayerID)
	ps.SettingName = reader.ReadString()
	ps.Value = reader.ReadString()
}
