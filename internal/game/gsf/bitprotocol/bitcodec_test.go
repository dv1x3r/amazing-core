package bitprotocol_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf/bitprotocol"
)

func TestLength(t *testing.T) {
	t.Run("BitCodec.ReadLength", func(t *testing.T) {
		tests := []struct {
			Input    []byte
			Expected int
		}{
			{Input: []byte{0x00}, Expected: 0},
			{Input: []byte{0x15}, Expected: 21},
			{Input: []byte{0x81, 0x00}, Expected: 128},
			{Input: []byte{0x81, 0x80, 0x00}, Expected: 16_384},
			{Input: []byte{0x81, 0x80, 0x80, 0x00}, Expected: 2_097_152},
		}
		for _, test := range tests {
			codec := bitprotocol.NewBitCodec()
			result, err := codec.ReadLength(bytes.NewBuffer(test.Input))
			if err != nil {
				t.Error(err)
			}
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitCodec.WriteLength", func(t *testing.T) {
		tests := []struct {
			Input    int
			Expected []byte
		}{
			{Input: 0, Expected: []byte{0x00}},
			{Input: 21, Expected: []byte{0x15}},
			{Input: 128, Expected: []byte{0x81, 0x00}},
			{Input: 16_384, Expected: []byte{0x81, 0x80, 0x00}},
			{Input: 2_097_152, Expected: []byte{0x81, 0x80, 0x80, 0x00}},
		}
		for _, test := range tests {
			buf := &bytes.Buffer{}
			codec := bitprotocol.NewBitCodec()
			codec.WriteLength(buf, test.Input)
			result := buf.Bytes()
			if !reflect.DeepEqual(result, test.Expected) {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})
}

func TestBool(t *testing.T) {
	t.Run("BitReader.ReadBool", func(t *testing.T) {
		tests := []struct {
			Input    string
			Expected bool
		}{
			{Input: "00", Expected: false},
			{Input: "80", Expected: true},
		}
		for _, test := range tests {
			data, _ := hex.DecodeString(test.Input)
			reader := bitprotocol.NewBitReader(data)
			result := reader.ReadBool()
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitWriter.WriteBool", func(t *testing.T) {
		tests := []struct {
			Input    bool
			Expected string
		}{
			{Input: false, Expected: "00"},
			{Input: true, Expected: "80"},
		}
		for _, test := range tests {
			buf := &bytes.Buffer{}
			writer := bitprotocol.NewBitWriter()
			writer.WriteBool(test.Input)
			writer.CommitTo(buf)
			result := fmt.Sprintf("%x", buf.Bytes())
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})
}

func TestInt(t *testing.T) {
	t.Run("BitReader.ReadInt16", func(t *testing.T) {
		tests := []struct {
			Input    string
			Expected int16
		}{
			{Input: "80", Expected: 0},
			{Input: "90", Expected: 4},
			{Input: "9c", Expected: 7},
			{Input: "c100", Expected: 8},
			{Input: "c800", Expected: 64},
			{Input: "cfe0", Expected: 127},
			{Input: "004000", Expected: 128},
			{Input: "3fff80", Expected: 32_767},
			{Input: "b0", Expected: -4},
			{Input: "a0", Expected: -8},
			{Input: "dee0", Expected: -9},
			{Input: "d800", Expected: -64},
			{Input: "d000", Expected: -128},
			{Input: "7fbf80", Expected: -129},
			{Input: "400000", Expected: -32_768},
		}
		for _, test := range tests {
			data, _ := hex.DecodeString(test.Input)
			reader := bitprotocol.NewBitReader(data)
			result := reader.ReadInt16()
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitWriter.WriteInt16", func(t *testing.T) {
		tests := []struct {
			Input    int16
			Expected string
		}{
			{Input: 0, Expected: "80"},
			{Input: 4, Expected: "90"},
			{Input: 7, Expected: "9c"},
			{Input: 8, Expected: "c100"},
			{Input: 64, Expected: "c800"},
			{Input: 127, Expected: "cfe0"},
			{Input: 128, Expected: "004000"},
			{Input: 32_767, Expected: "3fff80"},
			{Input: -4, Expected: "b0"},
			{Input: -8, Expected: "a0"},
			{Input: -9, Expected: "dee0"},
			{Input: -64, Expected: "d800"},
			{Input: -128, Expected: "d000"},
			{Input: -129, Expected: "7fbf80"},
			{Input: -32_768, Expected: "400000"},
		}
		for _, test := range tests {
			buf := &bytes.Buffer{}
			writer := bitprotocol.NewBitWriter()
			writer.WriteInt16(test.Input)
			writer.CommitTo(buf)
			result := fmt.Sprintf("%x", buf.Bytes())
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitReader.ReadInt32", func(t *testing.T) {
		tests := []struct {
			Input    string
			Expected int32
		}{
			{Input: "80", Expected: 0},
			{Input: "90", Expected: 4},
			{Input: "9c", Expected: 7},
			{Input: "c100", Expected: 8},
			{Input: "c800", Expected: 64},
			{Input: "cfe0", Expected: 127},
			{Input: "e00800", Expected: 128},
			{Input: "e7fff0", Expected: 32_767},
			{Input: "3fffffff80", Expected: 2_147_483_647},
			{Input: "b0", Expected: -4},
			{Input: "a0", Expected: -8},
			{Input: "dee0", Expected: -9},
			{Input: "d800", Expected: -64},
			{Input: "d000", Expected: -128},
			{Input: "eff7f0", Expected: -129},
			{Input: "e80000", Expected: -32_768},
			{Input: "4000000000", Expected: -2_147_483_648},
		}
		for _, test := range tests {
			data, _ := hex.DecodeString(test.Input)
			reader := bitprotocol.NewBitReader(data)
			result := reader.ReadInt32()
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitWriter.WriteInt32", func(t *testing.T) {
		tests := []struct {
			Input    int32
			Expected string
		}{
			{Input: 0, Expected: "80"},
			{Input: 4, Expected: "90"},
			{Input: 7, Expected: "9c"},
			{Input: 8, Expected: "c100"},
			{Input: 64, Expected: "c800"},
			{Input: 127, Expected: "cfe0"},
			{Input: 128, Expected: "e00800"},
			{Input: 32_767, Expected: "e7fff0"},
			{Input: 2_147_483_647, Expected: "3fffffff80"},
			{Input: -4, Expected: "b0"},
			{Input: -8, Expected: "a0"},
			{Input: -9, Expected: "dee0"},
			{Input: -64, Expected: "d800"},
			{Input: -128, Expected: "d000"},
			{Input: -129, Expected: "eff7f0"},
			{Input: -32_768, Expected: "e80000"},
			{Input: -2_147_483_648, Expected: "4000000000"},
		}
		for _, test := range tests {
			buf := &bytes.Buffer{}
			writer := bitprotocol.NewBitWriter()
			writer.WriteInt32(test.Input)
			writer.CommitTo(buf)
			result := fmt.Sprintf("%x", buf.Bytes())
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitReader.ReadInt64", func(t *testing.T) {
		tests := []struct {
			Input    string
			Expected int64
		}{
			{Input: "80", Expected: 0},
			{Input: "90", Expected: 4},
			{Input: "9c", Expected: 7},
			{Input: "c100", Expected: 8},
			{Input: "c800", Expected: 64},
			{Input: "cfe0", Expected: 127},
			{Input: "e00800", Expected: 128},
			{Input: "e7fff0", Expected: 32_767},
			{Input: "f9fffffffc", Expected: 2_147_483_647},
			{Input: "3fffffffffffffff80", Expected: 9_223_372_036_854_775_807},
			{Input: "b0", Expected: -4},
			{Input: "a0", Expected: -8},
			{Input: "dee0", Expected: -9},
			{Input: "d800", Expected: -64},
			{Input: "d000", Expected: -128},
			{Input: "eff7f0", Expected: -129},
			{Input: "e80000", Expected: -32_768},
			{Input: "fa00000000", Expected: -2_147_483_648},
			{Input: "400000000000000000", Expected: -9_223_372_036_854_775_808},
		}
		for _, test := range tests {
			data, _ := hex.DecodeString(test.Input)
			reader := bitprotocol.NewBitReader(data)
			result := reader.ReadInt64()
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitWriter.WriteInt64", func(t *testing.T) {
		tests := []struct {
			Input    int64
			Expected string
		}{
			{Input: 0, Expected: "80"},
			{Input: 4, Expected: "90"},
			{Input: 7, Expected: "9c"},
			{Input: 8, Expected: "c100"},
			{Input: 64, Expected: "c800"},
			{Input: 127, Expected: "cfe0"},
			{Input: 128, Expected: "e00800"},
			{Input: 32_767, Expected: "e7fff0"},
			{Input: 2_147_483_647, Expected: "f9fffffffc"},
			{Input: 9_223_372_036_854_775_807, Expected: "3fffffffffffffff80"},
			{Input: -4, Expected: "b0"},
			{Input: -8, Expected: "a0"},
			{Input: -9, Expected: "dee0"},
			{Input: -64, Expected: "d800"},
			{Input: -128, Expected: "d000"},
			{Input: -129, Expected: "eff7f0"},
			{Input: -32_768, Expected: "e80000"},
			{Input: -2_147_483_648, Expected: "fa00000000"},
			{Input: -9_223_372_036_854_775_808, Expected: "400000000000000000"},
		}
		for _, test := range tests {
			buf := &bytes.Buffer{}
			writer := bitprotocol.NewBitWriter()
			writer.WriteInt64(test.Input)
			writer.CommitTo(buf)
			result := fmt.Sprintf("%x", buf.Bytes())
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})
}

func TestFloat(t *testing.T) {
	t.Run("BitReader.ReadFloat32", func(t *testing.T) {
		tests := []struct {
			Input    string
			Expected float32
		}{
			{Input: "00000000", Expected: 0},
			{Input: "3f800000", Expected: 1},
			{Input: "3f99999a", Expected: 1.2},
			{Input: "3f9d70a4", Expected: 1.23},
			{Input: "3f9df3b6", Expected: 1.234},
			{Input: "3f9e0419", Expected: 1.2345},
			{Input: "4640e47e", Expected: 12345.12345},
			{Input: "c640e47e", Expected: -12345.12345},
		}
		for _, test := range tests {
			data, _ := hex.DecodeString(test.Input)
			reader := bitprotocol.NewBitReader(data)
			result := reader.ReadFloat32()
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitReader.WriteFloat32", func(t *testing.T) {
		tests := []struct {
			Input    float32
			Expected string
		}{
			{Input: 0, Expected: "00000000"},
			{Input: 1, Expected: "3f800000"},
			{Input: 1.2, Expected: "3f99999a"},
			{Input: 1.23, Expected: "3f9d70a4"},
			{Input: 1.234, Expected: "3f9df3b6"},
			{Input: 1.2345, Expected: "3f9e0419"},
			{Input: 12345.12345, Expected: "4640e47e"},
			{Input: -12345.12345, Expected: "c640e47e"},
		}
		for _, test := range tests {
			buf := &bytes.Buffer{}
			writer := bitprotocol.NewBitWriter()
			writer.WriteFloat32(test.Input)
			writer.CommitTo(buf)
			result := fmt.Sprintf("%x", buf.Bytes())
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitReader.ReadFloat64", func(t *testing.T) {
		tests := []struct {
			Input    string
			Expected float64
		}{
			{Input: "0000000000000000", Expected: 0},
			{Input: "3ff0000000000000", Expected: 1},
			{Input: "3ff3333333333333", Expected: 1.2},
			{Input: "3ff3ae147ae147ae", Expected: 1.23},
			{Input: "3ff3be76c8b43958", Expected: 1.234},
			{Input: "3ff3c083126e978d", Expected: 1.2345},
			{Input: "40c81c8fcd35a858", Expected: 12345.12345},
			{Input: "c0c81c8fcd35a858", Expected: -12345.12345},
		}
		for _, test := range tests {
			data, _ := hex.DecodeString(test.Input)
			reader := bitprotocol.NewBitReader(data)
			result := reader.ReadFloat64()
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitReader.WriteFloat64", func(t *testing.T) {
		tests := []struct {
			Input    float64
			Expected string
		}{
			{Input: 0, Expected: "0000000000000000"},
			{Input: 1, Expected: "3ff0000000000000"},
			{Input: 1.2, Expected: "3ff3333333333333"},
			{Input: 1.23, Expected: "3ff3ae147ae147ae"},
			{Input: 1.234, Expected: "3ff3be76c8b43958"},
			{Input: 1.2345, Expected: "3ff3c083126e978d"},
			{Input: 12345.12345, Expected: "40c81c8fcd35a858"},
			{Input: -12345.12345, Expected: "c0c81c8fcd35a858"},
		}
		for _, test := range tests {
			buf := &bytes.Buffer{}
			writer := bitprotocol.NewBitWriter()
			writer.WriteFloat64(test.Input)
			writer.CommitTo(buf)
			result := fmt.Sprintf("%x", buf.Bytes())
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})
}

func TestChar(t *testing.T) {
	t.Run("BitReader.ReadChar", func(t *testing.T) {
		tests := []struct {
			Input    string
			Expected rune
		}{
			{Input: "c820", Expected: 'A'},
			{Input: "c840", Expected: 'B'},
			{Input: "c860", Expected: 'C'},
			{Input: "020800", Expected: 'А'},
			{Input: "020880", Expected: 'Б'},
			{Input: "020900", Expected: 'В'},
		}
		for _, test := range tests {
			data, _ := hex.DecodeString(test.Input)
			reader := bitprotocol.NewBitReader(data)
			result := reader.ReadChar()
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitWriter.WriteChar", func(t *testing.T) {
		tests := []struct {
			Input    rune
			Expected string
		}{
			{Input: 'A', Expected: "c820"},
			{Input: 'B', Expected: "c840"},
			{Input: 'C', Expected: "c860"},
			{Input: 'А', Expected: "020800"},
			{Input: 'Б', Expected: "020880"},
			{Input: 'В', Expected: "020900"},
		}
		for _, test := range tests {
			buf := &bytes.Buffer{}
			writer := bitprotocol.NewBitWriter()
			writer.WriteChar(test.Input)
			writer.CommitTo(buf)
			result := fmt.Sprintf("%x", buf.Bytes())
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})
}

func TestString(t *testing.T) {
	t.Run("BitReader.ReadString", func(t *testing.T) {
		tests := []struct {
			Input    string
			Expected string
		}{
			{Input: "c180416d617a696e67576f726c64", Expected: "AmazingWorld"},
		}
		for _, test := range tests {
			data, _ := hex.DecodeString(test.Input)
			reader := bitprotocol.NewBitReader(data)
			result := reader.ReadString()
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitWriter.WriteString", func(t *testing.T) {
		tests := []struct {
			Input    string
			Expected string
		}{
			{Input: "AmazingWorld", Expected: "c180416d617a696e67576f726c64"},
		}
		for _, test := range tests {
			buf := &bytes.Buffer{}
			writer := bitprotocol.NewBitWriter()
			writer.WriteString(test.Input)
			writer.CommitTo(buf)
			result := fmt.Sprintf("%x", buf.Bytes())
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})
}

func TestBytes(t *testing.T) {
	t.Run("BitReader.ReadBytes", func(t *testing.T) {
		tests := []struct {
			Input    string
			Expected []byte
		}{
			{Input: "c180416d617a696e67576f726c64", Expected: []byte{0x41, 0x6D, 0x61, 0x7A, 0x69, 0x6E, 0x67, 0x57, 0x6F, 0x72, 0x6C, 0x64}},
		}
		for _, test := range tests {
			data, _ := hex.DecodeString(test.Input)
			reader := bitprotocol.NewBitReader(data)
			result := fmt.Sprintf("%x", reader.ReadBytes())
			if result != fmt.Sprintf("%x", test.Expected) {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitWriter.WriteBytes", func(t *testing.T) {
		tests := []struct {
			Input    []byte
			Expected string
		}{
			{Input: []byte{0x41, 0x6D, 0x61, 0x7A, 0x69, 0x6E, 0x67, 0x57, 0x6F, 0x72, 0x6C, 0x64}, Expected: "c180416d617a696e67576f726c64"},
		}
		for _, test := range tests {
			buf := &bytes.Buffer{}
			writer := bitprotocol.NewBitWriter()
			writer.WriteBytes(test.Input)
			writer.CommitTo(buf)
			result := fmt.Sprintf("%x", buf.Bytes())
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})
}

func TestUtcDate(t *testing.T) {
	t.Run("BitReader.ReadUtcDate", func(t *testing.T) {
		tests := []struct {
			Input    string
			Expected time.Time
		}{
			{Input: "80", Expected: time.Time{}},
			{Input: "000000076d1105c000", Expected: time.Date(2021, 7, 24, 0, 0, 0, 0, time.UTC)},
			{Input: "00000007704d9c8000", Expected: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
		}
		for _, test := range tests {
			data, _ := hex.DecodeString(test.Input)
			reader := bitprotocol.NewBitReader(data)
			result := reader.ReadUtcDate()
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

	t.Run("BitWriter.WriteUtcDate", func(t *testing.T) {
		tests := []struct {
			Input    time.Time
			Expected string
		}{
			{Input: time.Time{}, Expected: "80"},
			{Input: time.Date(2021, 7, 24, 0, 0, 0, 0, time.UTC), Expected: "000000076d1105c000"},
			{Input: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), Expected: "00000007704d9c8000"},
		}
		for _, test := range tests {
			buf := &bytes.Buffer{}
			writer := bitprotocol.NewBitWriter()
			writer.WriteUtcDate(test.Input)
			writer.CommitTo(buf)
			result := fmt.Sprintf("%x", buf.Bytes())
			if result != test.Expected {
				t.Errorf("❌ test=%+v result=%v\n", test, result)
			}
		}
	})

}
