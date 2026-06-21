package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
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

// OIDFromValues creates an OID out of values.
func OIDFromValues(class byte, oidtype byte, server byte, number int64) OID {
	return OID{
		Class:  class,
		Type:   oidtype,
		Server: server,
		Number: number,
	}
}

// OIDFromInt64 decodes a packed int64 into an OID.
func OIDFromInt64(v int64) OID {
	return OID{
		Class:  byte((v >> 56) & 0xFF),
		Type:   byte((v >> 48) & 0xFF),
		Server: byte((v >> 40) & 0xFF),
		Number: v & 0xFFFFFFFFFF,
	}
}

// OIDFromString converts string into int64, and decodes a packed int64 into an OID.
func OIDFromString(sv string) OID {
	v, _ := strconv.ParseInt(sv, 10, 64)
	return OIDFromInt64(v)
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
	value |= int64(oid.Number) & 0xFFFFFFFFFF
	return value
}

// String returns the dashed class-type-server-number representation of oid.
func (oid OID) String() string {
	return fmt.Sprintf("%d-%d-%d-%d", oid.Class, oid.Type, oid.Server, oid.Number)
}

// CDNID encodes OID integer into a base64 CDN ID.
func (oid OID) CDNID() string {
	return base64.RawStdEncoding.EncodeToString([]byte(strconv.Itoa(int(oid.Int64()))))
}

// UnmarshalJSON implements json.Unmarshaler for OID.
func (oid *OID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		*oid = OID{}
		return nil
	}

	// Handle quoted string: "1234567890"
	if len(data) > 1 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*oid = OIDFromInt64(v)
		return nil
	}

	// Handle raw number (legacy / internal use)
	var value int64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*oid = OIDFromInt64(value)
	return nil
}

// MarshalJSON implements json.Marshaler for OID.
func (oid OID) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(oid.Int64(), 10))
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

// Serialize implements gsf.Serializable for OID.
func (oid *OID) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(int32(oid.Class))
	writer.WriteInt32(int32(oid.Type))
	writer.WriteInt32(int32(oid.Server))
	writer.WriteInt64(oid.Number)
}

// Deserialize implements gsf.Deserializable for OID.
func (oid *OID) Deserialize(reader gsf.ProtocolReader) {
	oid.Class = byte(reader.ReadInt32())
	oid.Type = byte(reader.ReadInt32())
	oid.Server = byte(reader.ReadInt32())
	oid.Number = reader.ReadInt64()
}
