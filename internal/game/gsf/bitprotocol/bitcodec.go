package bitprotocol

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

var epochUnix = time.Date(1, 3, 1, 0, 0, 0, 0, time.UTC).Unix()

type BitCodec struct {
}

func NewBitCodec() *BitCodec {
	return &BitCodec{}
}

func (BitCodec) NewReader(data []byte) gsf.ProtocolReader {
	return NewBitReader(data)
}

func (BitCodec) NewWriter() gsf.ProtocolWriter {
	return NewBitWriter()
}

/*
ReadLength decodes an integer length from a variable-length format.

The function reads bytes from the provided reader and reconstructs the length by using the most significant bit (MSB) to indicate whether there are more bytes.
  - Each byte contributes 7 bits to the length, with the MSB (0x80) being used to chain subsequent bytes.
  - The process stops when a byte with an MSB of 0 is encountered, indicating that the final byte has been read.

The function returns the decoded length and any potential error encountered while reading from the reader.
*/
func (BitCodec) ReadLength(reader io.ByteReader) (int, error) {
	length := 0
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// Shift the current length to make room for the next 7 bits
		// (in case of multiple bytes)
		length <<= 7

		// Contribute 7 bits to the length
		length |= 0x7F & int(b)

		// Stop when MSB of 0 is encountered
		if (0x80 & b) == 0 {
			break
		}
	}
	return length, nil
}

/*
WriteLength encodes an integer length into a variable-length format and writes it to the writer.

The encoding uses the most significant bit (MSB) of each byte to indicate if there are additional bytes.
  - A single byte is used for lengths <= 127.
  - Two bytes for lengths <= 16,383.
  - Three bytes for lengths <= 2,097,151.
  - Four bytes for lengths <= 268,435,455.

The function ensures that the length is within the valid range (0 to 268,435,455) and errors if the length exceeds this range.
*/
func (BitCodec) WriteLength(writer io.ByteWriter, length int) error {
	if length > 268_435_455 {
		return fmt.Errorf("length is out of range: %d bytes", length)
	}
	// Write the most significant 7 bits of the length, with MSB set to 1 if another byte follows
	if length > 2_097_151 {
		writer.WriteByte(byte(0x80 | (0x7F & (length >> 21))))
	}
	if length > 16_383 {
		writer.WriteByte(byte(0x80 | (0x7F & (length >> 14))))
	}
	if length > 127 {
		writer.WriteByte(byte(0x80 | (0x7F & (length >> 7))))
	}
	return writer.WriteByte(byte(0x7F & length))
}

type BitReader struct {
	stream *BitStream
}

func NewBitReader(data []byte) *BitReader {
	return &BitReader{stream: NewBitStreamWithData(data)}
}

func (br *BitReader) ReadObject(value gsf.Deserializable) {
	if !br.stream.Get() {
		value.Deserialize(br)
	}
}

func (br *BitReader) ReadBool() bool {
	return br.stream.Get()
}

func (br *BitReader) ReadInt16() int16 {
	return int16(br.stream.GetIntCompressed(2))
}

func (br *BitReader) ReadInt32() int32 {
	return int32(br.stream.GetIntCompressed(4))
}

func (br *BitReader) ReadInt64() int64 {
	return br.stream.GetIntCompressed(8)
}

func (br *BitReader) ReadFloat32() float32 {
	bits := uint32(br.stream.GetInt(4))
	return math.Float32frombits(bits)
}

func (br *BitReader) ReadFloat64() float64 {
	bits := uint64(br.stream.GetInt(8))
	return math.Float64frombits(bits)
}

func (br *BitReader) ReadChar() rune {
	return rune(br.stream.GetIntCompressed(2))
}

func (br *BitReader) ReadString() string {
	length := br.stream.GetIntCompressed(4)
	array := br.stream.GetBytesAligned(int(length))
	return string(array)
}

func (br *BitReader) ReadBytes() []byte {
	length := br.stream.GetIntCompressed(4)
	array := br.stream.GetBytesAligned(int(length))
	return array
}

func (br *BitReader) ReadUtcDate() time.Time {
	if br.stream.Get() {
		return time.Time{}
	}
	seconds := br.stream.GetInt(8) - 31_622_400
	if seconds < 0 {
		seconds = 0
	}
	return time.Unix(seconds+epochUnix, 0).In(time.UTC)
}

type BitWriter struct {
	stream *BitStream
}

func NewBitWriter() *BitWriter {
	return &BitWriter{stream: NewBitStream()}
}

func (bw *BitWriter) CommitTo(writer io.Writer) {
	size := bw.stream.Length()
	writer.Write(bw.stream.buf[:size])
	bw.stream = NewBitStream()
}

func (bw *BitWriter) WriteObject(value gsf.Serializable) {
	if !bw.stream.Put(value == nil) {
		value.Serialize(bw)
	}
}

func (bw *BitWriter) WriteBool(value bool) {
	bw.stream.Put(value)
}

func (bw *BitWriter) WriteInt16(value int16) {
	bw.stream.PutIntCompressed(int64(value), 2)
}

func (bw *BitWriter) WriteInt32(value int32) {
	bw.stream.PutIntCompressed(int64(value), 4)
}

func (bw *BitWriter) WriteInt64(value int64) {
	bw.stream.PutIntCompressed(value, 8)
}

func (bw *BitWriter) WriteFloat32(value float32) {
	bits := int64(math.Float32bits(value))
	bw.stream.PutInt(bits, 4)
}

func (bw *BitWriter) WriteFloat64(value float64) {
	bits := int64(math.Float64bits(value))
	bw.stream.PutInt(bits, 8)
}

func (bw *BitWriter) WriteChar(value rune) {
	bw.stream.PutIntCompressed(int64(value), 2)
}

func (bw *BitWriter) WriteString(value string) {
	if len(value) > 2_147_483_647 {
		panic("string is too large to encode")
	}
	bw.stream.PutIntCompressed(int64(len(value)), 4)
	bw.stream.PutBytesAligned([]byte(value))
}

func (bw *BitWriter) WriteBytes(array []byte) {
	if len(array) > 2_147_483_647 {
		panic("array is too large to encode")
	}
	bw.stream.PutIntCompressed(int64(len(array)), 4)
	bw.stream.PutBytesAligned(array)
}

func (bw *BitWriter) WriteUtcDate(value time.Time) {
	if bw.stream.Put(value.IsZero()) {
		return
	}
	seconds := value.Unix() - epochUnix + 31_622_400
	bw.stream.PutInt(seconds, 8)
}
