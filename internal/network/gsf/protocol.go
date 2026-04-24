package gsf

import "io"

type Serializable interface {
	Serialize(writer ProtocolWriter)
}

type Deserializable interface {
	Deserialize(reader ProtocolReader)
}

type ProtocolCodec interface {
	NewReader(data []byte) ProtocolReader
	NewWriter() ProtocolWriter
	ReadLength(stream io.ByteReader) (int, error)
	WriteLength(stream io.ByteWriter, len int) error
}

type ProtocolReader interface {
	ReadObject(value Deserializable)
	GetByte() byte
	ReadBool() bool
	ReadInt16() int16
	ReadInt32() int32
	ReadInt64() int64
	ReadFloat32() float32
	ReadFloat64() float64
	ReadChar() rune
	ReadString() string
	ReadBytes() []byte
	ReadUtcDate() UnixTime
}

type ProtocolWriter interface {
	CommitTo(stream io.Writer)
	WriteObject(value Serializable)
	PutByte(value byte)
	WriteBool(value bool)
	WriteInt16(value int16)
	WriteInt32(value int32)
	WriteInt64(value int64)
	WriteFloat32(value float32)
	WriteFloat64(value float64)
	WriteChar(value rune)
	WriteString(value string)
	WriteBytes(value []byte)
	WriteUtcDate(value UnixTime)
}

func ReadNullable[T any](reader ProtocolReader, readFn func() T) Null[T] {
	var value Null[T]
	value.Valid = !reader.ReadBool()
	if value.Valid {
		value.V = readFn()
	}
	return value
}

func WriteNullable[T any](writer ProtocolWriter, value Null[T], writeFn func(value T)) {
	writer.WriteBool(!value.Valid)
	if value.Valid {
		writeFn(value.V)
	}
}

func ReadSlice[T any](reader ProtocolReader, readFn func() T) []T {
	length := reader.ReadInt32()
	if length < 0 {
		return nil
	}

	slice := make([]T, length)
	for i := range slice {
		slice[i] = readFn()
	}

	return slice
}

func WriteSlice[T any](writer ProtocolWriter, slice []T, writeFn func(value T)) {
	if slice == nil {
		writer.WriteInt32(-1)
		return
	}

	writer.WriteInt32(int32(len(slice)))
	for _, value := range slice {
		writeFn(value)
	}
}

func ReadMap[T any](reader ProtocolReader, readFn func() T) map[string]T {
	length := reader.ReadInt32()
	if length < 0 {
		return nil
	}

	dict := map[string]T{}
	for range length {
		k := reader.ReadString()
		dict[k] = readFn()
	}

	return dict
}

func WriteMap[T any](writer ProtocolWriter, dict map[string]T, writeFn func(value T)) {
	if dict == nil {
		writer.WriteInt32(-1)
		return
	}

	writer.WriteInt32(int32(len(dict)))
	for k, v := range dict {
		writer.WriteString(k)
		writeFn(v)
	}
}
