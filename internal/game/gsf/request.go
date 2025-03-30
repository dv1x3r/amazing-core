package gsf

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
)

type Request struct {
	header *Header
	body   Deserializable
	reader ProtocolReader
	ctx    context.Context

	RemoteAddr string
}

func NewRequest(ctx context.Context, header *Header, reader ProtocolReader) *Request {
	return &Request{
		header: header,
		reader: reader,
		ctx:    ctx,
	}
}

func (req *Request) Header() *Header {
	return req.header
}

func (req *Request) Body() Deserializable {
	return req.body
}

func (req *Request) Context() context.Context {
	return req.ctx
}

func (req *Request) Read(body Deserializable) error {
	req.body = body
	return wrap.Panic(func() error {
		req.reader.ReadObject(req.body)
		return nil
	})
}
