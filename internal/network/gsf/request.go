package gsf

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
)

type Request struct {
	header  *Header
	body    Deserializable
	reader  ProtocolReader
	ctx     context.Context
	session *Session
}

func NewRequest(ctx context.Context, header *Header, reader ProtocolReader, session *Session) *Request {
	return &Request{
		header:  header,
		reader:  reader,
		ctx:     ctx,
		session: session,
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

func (req *Request) Session() *Session {
	return req.session
}

func (req *Request) RemoteIP() string {
	return req.session.RemoteIP()
}

func (req *Request) Platform() Platform {
	return req.session.Platform()
}

func (req *Request) SetPlatform(platform Platform) {
	req.session.SetPlatform(platform)
}

func (req *Request) PlayerOID() (int64, bool) {
	return req.session.PlayerOID()
}

func (req *Request) SetPlayerOID(playerOID int64) {
	req.session.SetPlayerOID(playerOID)
}

func (req *Request) Read(body Deserializable) error {
	req.body = body
	return wrap.Panic(func() error {
		req.reader.ReadObject(req.body)
		return nil
	})
}
