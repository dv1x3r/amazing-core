package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// PlayerStats is a player avatar statistic entry.
type PlayerStats struct {
	OID             OID
	PlayerAvatarOID OID
	StatsTypeOID    OID
	Level           int32
	ObjectOID       OID
}

func (ps *PlayerStats) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&ps.OID)
	writer.WriteObject(&ps.PlayerAvatarOID)
	writer.WriteObject(&ps.StatsTypeOID)
	writer.WriteInt32(ps.Level)
	writer.WriteObject(&ps.ObjectOID)
}

func (ps *PlayerStats) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&ps.OID)
	reader.ReadObject(&ps.PlayerAvatarOID)
	reader.ReadObject(&ps.StatsTypeOID)
	ps.Level = reader.ReadInt32()
	reader.ReadObject(&ps.ObjectOID)
}
