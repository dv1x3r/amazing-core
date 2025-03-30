package gsf

import (
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
)

type Header struct {
	Flags         int32
	SvcClass      int32
	MsgType       int32
	RequestID     int32
	LogCorrelator string
	ResultCode    int32
	AppCode       int32
	AppString     string
	AppCodeArray  []struct {
		Code int32
		Str  string
	}
}

func ReadHeader(reader ProtocolReader) (*Header, error) {
	header := &Header{}
	return header, wrap.Panic(func() error {
		if reader.ReadBool() {
			return fmt.Errorf("null request")
		}
		reader.ReadObject(header)
		return nil
	})
}

func (h *Header) IsService() bool {
	return h.Flags&2 == 0
}

func (h *Header) IsResponse() bool {
	return h.IsService() && ((h.Flags & 1) != 0)
}

func (h *Header) SetResponse(value bool) {
	if value {
		h.Flags |= 1
	} else {
		h.Flags &= -1
	}
}

func (h *Header) IsRequest() bool {
	return h.IsService() && !h.IsResponse()
}

func (h *Header) IsNotify() bool {
	return h.Flags&2 != 0
}

func (h *Header) IsDiscardable() bool {
	return h.Flags&0x10 != 0
}

func (h *Header) SetDiscardable(value bool) {
	if value {
		h.Flags |= 16
	} else {
		h.Flags &= -17
	}
}

func (h *Header) Serialize(writer ProtocolWriter) {
	writer.WriteInt32(h.Flags)
	writer.WriteInt32(h.SvcClass)
	writer.WriteInt32(h.MsgType)
	if h.IsService() {
		writer.WriteInt32(h.RequestID)
	}
	if h.IsRequest() {
		writer.WriteString(h.LogCorrelator)
	}
	if h.IsResponse() {
		writer.WriteInt32(h.ResultCode)
		writer.WriteInt32(h.AppCode)
		if h.AppCode != 0 {
			writer.WriteString(h.AppString)
		}
		if h.AppCode == 17 {
			writer.WriteInt32(int32(len(h.AppCodeArray)))
			for i := range h.AppCodeArray {
				writer.WriteInt32(h.AppCodeArray[i].Code)
				writer.WriteString(h.AppCodeArray[i].Str)
			}
		}
	}
}

func (h *Header) Deserialize(reader ProtocolReader) {
	h.Flags = reader.ReadInt32()
	h.SvcClass = reader.ReadInt32()
	h.MsgType = reader.ReadInt32()
	if h.IsService() {
		h.RequestID = reader.ReadInt32()
	}
	if h.IsRequest() {
		h.LogCorrelator = reader.ReadString()
	}
}
