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
	conn   *Connection
}

func NewRequest(ctx context.Context, header *Header, reader ProtocolReader, conn *Connection) *Request {
	return &Request{
		header: header,
		reader: reader,
		ctx:    ctx,
		conn:   conn,
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

func (req *Request) RemoteIP() string {
	return req.conn.RemoteIP()
}

func (req *Request) Platform() Platform {
	return req.conn.Platform()
}

func (req *Request) SetPlatform(platform Platform) {
	req.conn.SetPlatform(platform)
}

func (req *Request) Read(body Deserializable) error {
	req.body = body
	return wrap.Panic(func() error {
		req.reader.ReadObject(req.body)
		return nil
	})
}
