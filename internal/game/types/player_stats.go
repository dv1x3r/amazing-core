package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type PlayerStats struct {
	OID            OID
	PlayerAvatarID OID
	StatsTypeID    OID
	Level          int32
	ObjectID       OID
}

func (ps *PlayerStats) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&ps.OID)
	writer.WriteObject(&ps.PlayerAvatarID)
	writer.WriteObject(&ps.StatsTypeID)
	writer.WriteInt32(ps.Level)
	writer.WriteObject(&ps.ObjectID)
}

func (ps *PlayerStats) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&ps.OID)
	reader.ReadObject(&ps.PlayerAvatarID)
	reader.ReadObject(&ps.StatsTypeID)
	ps.Level = reader.ReadInt32()
	reader.ReadObject(&ps.ObjectID)
}
