package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/dv1x3r/amazing-core/internal/network/gsf"
)

// OID is the GSF object identifier split into class, type, server, and number.
type OID struct {
	Class  byte
	Type   byte
	Server byte
	Number int64
}

// OIDFromInt64 decodes a packed int64 into an OID.
func OIDFromInt64(v int64) OID {
	var oid OID
	oid.Class = byte((v >> 56) & 0xFF)
	oid.Type = byte((v >> 48) & 0xFF)
	oid.Server = byte((v >> 40) & 0xFF)
	oid.Number = v & 0xFFFFFFFFFF
	return oid
}

// OIDFromCDNID decodes a CDN ID string into an OID.
func OIDFromCDNID(cdnid string) (OID, error) {
	str, err := base64.RawStdEncoding.DecodeString(cdnid)
	if err != nil {
		return OID{}, err
	}
	v, err := strconv.ParseInt(string(str), 10, 64)
	if err != nil {
		return OID{}, err
	}
	return OIDFromInt64(v), nil
}

// Int64 returns the packed int64 representation of oid.
func (oid OID) Int64() int64 {
	var value int64
	value |= int64(oid.Class) << 56
	value |= int64(oid.Type) << 48
	value |= int64(oid.Server) << 40
	value |= int64(oid.Number)
	return value
}

// String returns the dashed class-type-server-number representation of oid.
func (oid OID) String() string {
	return fmt.Sprintf("%d-%d-%d-%d", oid.Class, oid.Type, oid.Server, oid.Number)
}

// Scan implements sql.Scanner for OID.
func (oid *OID) Scan(value any) error {
	var n sql.NullInt64
	if err := n.Scan(value); err != nil {
		return err
	}
	if n.Valid {
		*oid = OIDFromInt64(n.Int64)
	}
	return nil
}

// Value implements driver.Valuer for OID.
func (oid OID) Value() (driver.Value, error) {
	return oid.Int64(), nil
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
