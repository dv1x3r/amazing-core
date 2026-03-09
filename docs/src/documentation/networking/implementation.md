# Go Implementation

## Reading a message

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
req := gsf.NewRequest(ctx, header, reader)
res := gsf.NewResponse(header, writer)

// 7. Call the handler
handler(res, req)
```

## Implementing a handler

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

## Defining a message type

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

## Collection helpers

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
