package types

import (
	"encoding/base64"
	"strconv"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type OID struct {
	Class  byte
	Type   byte
	Server byte
	Number int64
}

func (oid *OID) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(int32(oid.Class))
	writer.WriteInt32(int32(oid.Type))
	writer.WriteInt32(int32(oid.Server))
	writer.WriteInt64(oid.Number)
}

func (oid *OID) Deserialize(reader gsf.ProtocolReader) {
	oid.Class = byte(reader.ReadInt32())
	oid.Type = byte(reader.ReadInt32())
	oid.Server = byte(reader.ReadInt32())
	oid.Number = reader.ReadInt64()
}

func (oid *OID) ToLong() int {
	var value int
	value |= int(oid.Class) << 56
	value |= int(oid.Type) << 48
	value |= int(oid.Server) << 40
	value |= int(oid.Number)
	return value
}

func (oid *OID) FromLong(value int) {
	oid.Class = byte((value >> 56) & 0xFF)
	oid.Type = byte((value >> 48) & 0xFF)
	oid.Server = byte((value >> 40) & 0xFF)
	oid.Number = int64(value & 0xFFFFFFFFFF)
}

func (oid *OID) FromCDNID(cdnid string) {
	str, _ := base64.RawStdEncoding.DecodeString(cdnid)
	value, _ := strconv.Atoi(string(str))
	oid.FromLong(value)
}
