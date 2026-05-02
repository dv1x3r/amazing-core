# Networking

## Overview

The game client communicates with the server over a persistent **TCP connection** using a custom binary protocol.
All data is serialized with the **BitProtocol codec** - a bit-level encoding scheme that compresses integers by encoding only the bytes that are actually needed.

The implementation lives in `internal/network/`

## Wire Format

Every message follows this structure:

```
┌─────────────────────────────────────────────────────────┐
│  length prefix  │        payload bytes          │ 0x00  │
│  (varint)       │     (BitStream encoding)      │ (EOF) │
└─────────────────────────────────────────────────────────┘
```

### Length prefix

The length prefix is a **variable-length big-endian integer** where the MSB of each byte signals whether more bytes follow:

Each group of 7 bits contributes to the value, MSB = continuation flag:

```
[1xxxxxxx] [1xxxxxxx] [1xxxxxxx] [0xxxxxxx]
 more bytes more bytes more bytes last byte
```

### Payload

The payload is a **BitStream** - a sequence of bits packed into bytes.
The stream is read and written bit-by-bit; byte boundaries are only respected for raw string/byte-array data (aligned reads/writes).

After the last payload byte comes a **null terminator** `0x00`.

## BitStream Encoding

### Objects (nullable)

Objects begin with a **null flag bit**:

```
[0] -> object is not null
[1] -> null / not present
```

### Booleans

A single **bit**: `1` = true, `0` = false.

### Compressed integers (`int16`, `int32`, `int64`)

Integers are compressed using a variable-width encoding that minimizes bit usage for small values.

The encoding starts with a **compressed flag** bit:

```
[1] -> compressed (fewer than max bytes needed)
      followed by a count of how many bytes are needed
      [0]   -> 0 bytes -> value fits in 4 bits
      [10]  -> 1 byte  -> value is 1 byte  (8 bits)
      [110] -> 2 bytes -> value is 2 bytes (16 bits)
      ...
[0] -> full width (all bytes written)
```

### Floating-point numbers (`float32`, `float64`)

Written as raw not compressed IEEE 754 bytes.

### Strings

Strings are encoded as:

1. A compressed `int32` for the **UTF-8 byte length**.
2. The raw UTF-8 bytes, **byte-aligned** (any remaining bits in the current byte are skipped before writing).

### Byte arrays

Same as strings: compressed length prefix followed by byte-aligned data.

### Dates

Dates are encoded relative to a fixed epoch (`0001-03-01`).
A null flag bit precedes the value:

```
[1] -> zero time (null)
[0] -> seconds since epoch, written as a raw 8-byte integer
```

## Message Structure

Every message has a **header** followed by an optional **body**.
The outer wrapper adds a null flag for the entire message object.

```
[0]          -> message is not null
[header]     -> GSF header fields
[0]          -> body is not null
[body fields]
```

### GSF Header

| Field           | Type     | Condition          | Description                         |
| --------------- | -------- | ------------------ | ----------------------------------- |
| `flags`         | `int32`  | always             | Bit flags                           |
| `svcClass`      | `int32`  | always             | Message service class               |
| `msgType`       | `int32`  | always             | Message type for the given svcClass |
| `requestId`     | `int32`  | if `IsService()`   | Correlates request/response pairs   |
| `logCorrelator` | `string` | if `IsRequest()`   | Client-side correlation string      |
| `resultCode`    | `int32`  | if `IsResponse()`  | Result code (0 = success)           |
| `appCode`       | `int32`  | if `IsResponse()`  | Application-level error code        |
| `appString`     | `string` | if `appCode != 0`  | Error description string            |
| `appCodes`      | `array`  | if `appCode == 17` | Extended error code list            |

**Flag semantics:**

| Bit               | Meaning                                 |
| ----------------- | --------------------------------------- |
| `flags & 2 == 0`  | IsService - a request/response pair     |
| `flags & 1 != 0`  | IsResponse (only when IsService)        |
| `flags & 2 != 0`  | IsNotify - fire-and-forget notification |
| `flags & 16 != 0` | IsDiscardable                           |

## Annotated Example

The very first message the client sends after connecting is a **GetClientVersionInfo** request. It is 22 bytes total:

```
hex:  15 20 c2 5c 04 6d 0c 0c 18 41 6d 61 7a 69 6e 67 57 6f 72 6c 64 00
```

Breakdown:

| Bytes                                 | Value          | Meaning                            |
| ------------------------------------- | -------------- | ---------------------------------- |
| `15`                                  | 21             | Length prefix: payload is 21 bytes |
| `20 c2 5c 04 6d 0c 0c 18`             | bit stream     | Header + body                      |
| `41 6d 61 7a 69 6e 67 57 6f 72 6c 64` | `AmazingWorld` | String content (byte-aligned)      |
| `00`                                  | 0              | Null terminator                    |

Bit-by-bit payload decoding:

```
Bit  8:  [0]                -> message is not null
Bit  9:  [0]                -> header is not null
Bit 10:  [1][0]             -> flags compressed, 0 bytes -> 4-bit value
Bit 14:  [0000]             -> flags = 0
Bit 16:  [1][1][0]          -> svcClass compressed, 1 byte follows
Bit 19:  [00010010]         -> svcClass = 18 (USER_SERVER)
Bit 27:  [1][1][1][0]       -> msgType compressed, 2 bytes follow
Bit 31:  [0000001000110110] -> msgType = 566 (GET_CLIENT_VERSION_INFO)
Bit 47:  [1][0]             -> requestId compressed, 0 bytes -> 4-bit value
Bit 51:  [0001]             -> requestId = 1
Bit 53:  [1][0]             -> logCorrelator length compressed, 0 bytes -> 4-bit value
Bit 57:  [0000]             -> length = 0 (empty string)
Bit 59:  [0]                -> body is not null
Bit 60:  [1][1][0]          -> clientName length compressed, 1 byte follows
Bit 63:  [00001100]         -> length = 12 (bytes in "AmazingWorld")
Bit 71:  [0] align          -> skip to byte boundary
Bits 72–167:                -> "AmazingWorld" (12 bytes, 96 bits)
```

## Go Implementation

### Reading a message

```go
// 1. Read framing length from buffered reader
length, err := codec.ReadLength(stream)

// 2. Read 'length' bytes into data buffer
data := make([]byte, length)
io.ReadFull(stream, data)

// 3. Create a BitReader over the raw bytes
reader := bitprotocol.NewBitReader(data)

// 4. Parse the header
header, err := gsf.ReadHeader(reader)

// 5. Look up the handler
handler, ok := router.Lookup(header.SvcClass, header.MsgType)

// 6. Build request and response objects
conn := &gsf.Connection{remoteIP: remoteAddr}
req := gsf.NewRequest(ctx, header, reader, conn)
res := gsf.NewResponse(header, writer)

// 7. Call the handler
handler(res, req)
```

### Implementing a handler

```go
func GetClientVersionInfo(w gsf.ResponseWriter, r *gsf.Request) error {
    // Decode the request body
    req := &messages.GetClientVersionInfoRequest{}
    if err := r.Read(req); err != nil {
        return err
    }

    // Build and send the response
    res := &messages.GetClientVersionInfoResponse{
        ClientVersionInfo: "133852.true",
    }
    return w.Write(res)
}
```

### Defining a message type

Request and response structs implement `gsf.Deserializable` and `gsf.Serializable`:

```go
type GetClientVersionInfoRequest struct {
    ClientName string
}

func (req *GetClientVersionInfoRequest) Deserialize(reader gsf.ProtocolReader) {
    req.ClientName = reader.ReadString()
}

type GetClientVersionInfoResponse struct {
    ClientVersionInfo string
}

func (res *GetClientVersionInfoResponse) Serialize(writer gsf.ProtocolWriter) {
    writer.WriteString(res.ClientVersionInfo)
}
```

### Collection helpers

The `gsf` package provides generic helpers for encoding and decoding slices, maps, and nullable values:

```go
// Slices: int32 length prefix (−1 = nil), then elements
gsf.ReadSlice(reader, func() MyType { ... })
gsf.WriteSlice(writer, slice, func(v MyType) { ... })

// Maps: int32 length prefix (−1 = nil), then string key + value pairs
gsf.ReadMap(reader, func() MyType { ... })
gsf.WriteMap(writer, dict, func(v MyType) { ... })

// Nullable objects: bool null flag + optional value
gsf.ReadNullable(reader, func() MyType { ... })
gsf.WriteNullable(writer, value, func(v MyType) { ... })
```
