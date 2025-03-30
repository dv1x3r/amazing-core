package gsf

import "github.com/dv1x3r/amazing-core/internal/lib/wrap"

type ResponseWriter interface {
	Header() *Header
	Body() Serializable
	Write(Serializable) error
}

type Response struct {
	header *Header
	body   Serializable
	writer ProtocolWriter
}

func NewResponse(header *Header, writer ProtocolWriter) *Response {
	return &Response{
		header: &Header{
			Flags:         header.Flags | 1,
			SvcClass:      header.SvcClass,
			MsgType:       header.MsgType,
			RequestID:     header.RequestID,
			LogCorrelator: header.LogCorrelator,
		},
		writer: writer,
	}
}

func (res *Response) Header() *Header {
	return res.header
}

func (res *Response) Body() Serializable {
	return res.body
}

func (res *Response) Write(body Serializable) error {
	res.body = body
	return wrap.Panic(func() error {
		res.writer.WriteBool(false)
		res.writer.WriteObject(res.header)
		res.writer.WriteObject(res.body)
		return nil
	})
}
