package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type ZoneInstance struct {
	OID       OID
	Zone      Zone
	Occupancy int32
	Ordinal   int32
}

func (zi *ZoneInstance) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&zi.OID)
	writer.WriteObject(&zi.Zone)
	writer.WriteInt32(zi.Occupancy)
	writer.WriteInt32(zi.Ordinal)
}

func (zi *ZoneInstance) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&zi.OID)
	reader.ReadObject(&zi.Zone)
	zi.Occupancy = reader.ReadInt32()
	zi.Ordinal = reader.ReadInt32()
}
