package bitprotocol

import "math"

var (
	maskBml    = [8]byte{0, 128, 192, 224, 240, 248, 252, 254}
	maskBmr    = [8]byte{255, 127, 63, 31, 15, 7, 3, 1}
	maskBm     = [8]byte{128, 64, 32, 16, 8, 4, 2, 1}
	widthValid = [9]bool{false, false, true, false, true, false, false, false, true}
	widthMax   = [9]int64{7, 127, 32_767, 8_388_607, 2_147_483_647, 549_755_813_887, 140_737_488_355_327, 140_737_488_355_327, 9_223_372_036_854_775_807}
	widthMin   = [9]int64{-8, -128, -32_768, -8_388_608, -2_147_483_648, -549_755_813_888, -140_737_488_355_328, -140_737_488_355_328, -9_223_372_036_854_775_808}
	widthSBit  = [9]int64{8, 128, 32_768, 8_388_608, 2_147_483_648, 549_755_813_888, 140_737_488_355_328, 36_028_797_018_963_968, -9_223_372_036_854_775_808}
)

type BitStream struct {
	buf []byte
	lim int
	pos int
}

func NewBitStreamWithData(data []byte) *BitStream {
	return &BitStream{
		buf: data,
		lim: len(data) << 3,
	}
}

func NewBitStream() *BitStream {
	data := make([]byte, 1024)
	return NewBitStreamWithData(data)
}

/*
IndexByte returns the byte index corresponding to the given bit index.

The bit index is divided by 8 to determine the byte index:
  - Bit indices 1-7 return 0 (first byte)
  - Bit indices 8-15 return 1 (second byte)
  - Bit indices 16-23 return 2 (third byte)
  - And so on...
*/
func (bs *BitStream) IndexByte(index int) int {
	return index >> 3
}

/*
IndexBit returns the bit position within the byte corresponding to the given bit index.

The bit index is taken modulo 8 to determine the position within the byte:
  - Bit indices 0-7 return 0-7 (positions within the first byte)
  - Bit indices 8-15 return 0-7 (positions within the second byte)
  - Bit indices 16-23 return 0-7 (positions within the third byte)
  - And so on...
*/
func (bs *BitStream) IndexBit(index int) int {
	return index % 8
}

// Length returns the total number of bytes in the bitstream.
func (bs *BitStream) Length() int {
	return bs.IndexByte(bs.pos + 8 - 1)
}

/*
Resize resizes the bitstream buffer to ensure it can accommodate at least the requested number of bits.

The new buffer size is determined by taking the larger value between:
  - Doubling the current buffer size, or
  - The requested number of bytes (min) plus an additional 8192 bytes for overhead.
*/
func (bs *BitStream) Resize(min int) {
	minBytes := bs.IndexByte(min+8-1) + 8192
	newSize := int(math.Max(float64(len(bs.buf)*2), float64(minBytes)))
	newBuf := make([]byte, newSize)
	copy(newBuf, bs.buf)
	bs.buf = newBuf
	bs.lim = newSize << 3
}

// Put writes a single bit to the bitstream buffer.
func (bs *BitStream) Put(bit bool) bool {
	if bs.pos+1 > bs.lim {
		bs.Resize(bs.pos + 1)
	}

	if bit {
		bs.buf[bs.IndexByte(bs.pos)] |= maskBm[bs.IndexBit(bs.pos)]
	} else {
		bs.buf[bs.IndexByte(bs.pos)] &= ^maskBm[bs.IndexBit(bs.pos)]
	}

	bs.pos += 1
	return bit
}

// PutByte writes a single byte (8 bits) to the bitstream buffer.
func (bs *BitStream) PutByte(b byte) {
	if bs.pos+8 > bs.lim {
		bs.Resize(bs.pos + 8)
	}

	indexByte := bs.IndexByte(bs.pos)
	indexBit := bs.IndexBit(bs.pos)

	if indexBit == 0 {
		// beginning of a new byte
		bs.buf[indexByte] = b
	} else {
		// inserting into the middle of an existing byte
		tmp := int(0xFF&b) << (8 - indexBit)   // align the byte b with the correct position in the current byte
		bs.buf[indexByte] &= maskBml[indexBit] // clear bits that will be overwritten
		bs.buf[indexByte] |= byte(tmp >> 8)    // write the shifted byte
		// write the remaining bits into the next byte
		bs.buf[indexByte+1] = maskBml[indexBit] & byte(tmp)
	}

	bs.pos += 8
}

// PutBytesAligned writes a byte array to the bitstream buffer, starting from the current byte boundary.
func (bs *BitStream) PutBytesAligned(a []byte) {
	if len(a) == 0 {
		return
	}

	indexByte := bs.IndexByte(bs.pos)
	indexBit := bs.IndexBit(bs.pos)
	toReserve := len(a) * 8

	// skip the remaining byte
	if indexBit != 0 {
		toReserve += 8 - indexBit
		indexByte += 1
	}

	if bs.pos+toReserve > bs.lim {
		bs.Resize(bs.lim + toReserve)
	}

	copy(bs.buf[indexByte:indexByte+len(a)], a)
	bs.pos += toReserve
}

func (bs *BitStream) Peek() bool {
	if bs.pos+1 > bs.lim {
		panic("buffer underflow")
	}
	return (bs.buf[bs.IndexByte(bs.pos)])&(maskBm[bs.IndexBit(bs.pos)]) != 0
}

func (bs *BitStream) PeekByte() byte {
	if bs.pos+8 > bs.lim {
		panic("buffer underflow")
	}

	indexByte := bs.IndexByte(bs.pos)
	indexBit := bs.IndexBit(bs.pos)

	if indexBit == 0 {
		return bs.buf[indexByte]
	}

	b := (bs.buf[indexByte] & maskBmr[indexBit]) << indexBit
	return b | ((bs.buf[indexByte+1] & maskBml[indexBit]) >> (8 - indexBit))
}

// Get reads a single bit from the bitstream buffer.
func (bs *BitStream) Get() bool {
	res := bs.Peek()
	bs.pos += 1
	return res
}

// GetByte reads a single byte (8 bits) from the bitstream buffer.
func (bs *BitStream) GetByte() byte {
	res := bs.PeekByte()
	bs.pos += 8
	return res
}

// GetBytesAligned reads a byte array from the bitstream buffer, starting from the current byte boundary.
func (bs *BitStream) GetBytesAligned(count int) []byte {
	if count <= 0 {
		return nil
	}

	a := make([]byte, count)

	indexByte := bs.IndexByte(bs.pos)
	indexBit := bs.IndexBit(bs.pos)
	toRead := count * 8

	if indexBit != 0 {
		toRead += 8 - indexBit
		indexByte += 1
	}

	if bs.pos+toRead > bs.lim {
		panic("buffer underflow")
	}

	copy(a, bs.buf[indexByte:indexByte+count])
	bs.pos += toRead
	return a
}

// PutInt writes a non-compressed integer to the bitstream buffer.
func (bs *BitStream) PutInt(val int64, w int) {
	if !widthValid[w] {
		panic("invalid width")
	}

	for num := (w - 1) * 8; num >= 0; num -= 8 {
		bs.PutByte(byte(0xFF & (val >> num)))
	}
}

// GetInt reads a non-compressed integer from the bitstream buffer.
func (bs *BitStream) GetInt(w int) int64 {
	var res int64
	for i := 0; i < w; i++ {
		res = (res << 8) | int64(0xFF&bs.GetByte())
	}
	return res
}

// PutIntCompressed writes a compressed integer to the bitstream buffer.
func (bs *BitStream) PutIntCompressed(val int64, w int) {
	if !widthValid[w] {
		panic("invalid width")
	}

	toWrite := 0 // number of bytes required
	if val > 0 {
		for toWrite < w && val > widthMax[toWrite] {
			toWrite++
		}
	} else {
		for toWrite < w && val < widthMin[toWrite] {
			toWrite++
		}
	}

	// write the bit indicating if the full width is required
	if bs.Put(toWrite < w) {
		// write a sequence of bits showing how many bytes are required
		for j := 0; bs.Put(j < toWrite); j++ {
		}
	}

	if toWrite > 0 {
		// write the value in i chunks of 8 bits
		for num := (toWrite - 1) * 8; num >= 0; num -= 8 {
			bs.PutByte(byte(0xFF & (val >> num)))
		}
	} else {
		// write the integer using 4 bits
		for num := int64(8); num > 0; num >>= 1 {
			bs.Put((val & num) != 0)
		}
	}
}

// GetIntCompressed reads a compressed integer from the bitstream buffer.
func (bs *BitStream) GetIntCompressed(w int) int64 {

	var toRead int // number of bytes required
	if bs.Get() {
		// read a sequence of bits showing how many bytes are required
		for toRead < w && bs.Get() {
			toRead++
		}
	} else {
		// full width is required
		toRead = w
	}

	var res int64
	if toRead > 0 {
		// read the value in i chunks of 8 bits
		for j := 0; j < toRead; j++ {
			res = (res << 8) + int64(0xFF&bs.GetByte())
		}
	} else {
		// read the integer using 4 bits
		for k := 0; k < 4; k++ {
			if bs.Get() {
				res = (res << 1) + 1
			} else {
				res = res << 1
			}
		}
	}

	// switch negative values
	if (res & widthSBit[toRead]) != 0 {
		res |= widthMin[toRead]
	}

	return res
}
