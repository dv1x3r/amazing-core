package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type OID struct {
	Class  int32
	Type   int32
	Server int32
	Number int64
}

func (oid *OID) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(oid.Class)
	writer.WriteInt32(oid.Type)
	writer.WriteInt32(oid.Server)
	writer.WriteInt64(oid.Number)
}

func (oid *OID) Deserialize(reader gsf.ProtocolReader) {
	oid.Class = reader.ReadInt32()
	oid.Type = reader.ReadInt32()
	oid.Server = reader.ReadInt32()
	oid.Number = reader.ReadInt64()
}
