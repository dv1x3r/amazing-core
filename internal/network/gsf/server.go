package gsf

import (
	"bufio"
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/resultcode"
)

type Server struct {
	Addr   string
	Router *Router
	Codec  ProtocolCodec
	Hooks  ServerHooks

	listener net.Listener
	cancel   context.CancelFunc
	conns    sync.Map
	wg       sync.WaitGroup
}

type ServerHooks struct {
	OnConnect    func(conn *Conn)
	OnDisconnect func(conn *Conn, reason string)
	OnUnhandled  func(conn *Conn, header *Header, data []byte)
	OnRequest    func(event RequestEvent)
	OnNotify     func(event NotifyEvent)
}

type RequestEvent struct {
	Conn           *Conn
	RequestHeader  *Header
	ResponseHeader *Header
	RequestBody    Deserializable
	ResponseBody   Serializable
	Err            error
	Latency        time.Duration
	Dropped        bool
}

type NotifyEvent struct {
	Conn    *Conn
	Header  *Header
	Body    Serializable
	Err     error
	Latency time.Duration
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.listener = listener

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return err
			}
			time.Sleep(1 * time.Second)
			continue
		}

		s.conns.Store(conn, struct{}{})
		s.wg.Go(func() {
			defer s.conns.Delete(conn)
			defer conn.Close()
			s.handleConn(ctx, conn)
		})
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.listener == nil {
		return nil
	}

	s.listener.Close()
	s.cancel()

	// unblock connections that are waiting for read input
	s.conns.Range(func(k, _ any) bool {
		conn := k.(net.Conn)
		conn.SetReadDeadline(time.Now())
		return true
	})

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// all connections finished gracefully
	case <-ctx.Done():
		// force close all connections
		s.conns.Range(func(k, _ any) bool {
			conn := k.(net.Conn)
			conn.Close()
			return true
		})
	}

	return nil
}

func (s *Server) handleConn(ctx context.Context, netConn net.Conn) {
	stream := bufio.NewReadWriter(bufio.NewReader(netConn), bufio.NewWriter(netConn))
	remoteAddr, _, _ := net.SplitHostPort(netConn.RemoteAddr().String())
	conn := NewConn(remoteAddr, stream, s.Codec, s.Hooks.OnNotify)
	if s.Hooks.OnConnect != nil {
		s.Hooks.OnConnect(conn)
	}

	for {
		if err := s.processRequest(ctx, conn); err != nil {
			if s.Hooks.OnDisconnect != nil {
				s.Hooks.OnDisconnect(conn, disconnectReason(err))
			}
			return
		}
	}
}

func disconnectReason(err error) string {
	if errors.Is(err, io.EOF) {
		return "eof"
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return "timeout"
	}
	return err.Error()
}

func (s *Server) processRequest(ctx context.Context, conn *Conn) error {
	data, err := s.readMessage(conn.stream)
	if err != nil {
		return err
	}

	reader := s.Codec.NewReader(data)

	header := &Header{}
	err = wrap.Panic(func() error {
		if reader.ReadBool() {
			return fmt.Errorf("null request")
		}
		reader.ReadObject(header)
		return nil
	})
	if err != nil {
		return err
	}

	handler, ok := s.Router.Lookup(header.SvcClass, header.MsgType)
	if !ok {
		if s.Hooks.OnUnhandled != nil {
			s.Hooks.OnUnhandled(conn, header, data)
		}
		return nil
	}

	req := NewRequest(ctx, header, reader, conn)
	res := NewResponse(header)

	startTime := time.Now()
	handlerErr := handler(res, req)
	sendErr := respond(conn, res, handlerErr)

	if s.Hooks.OnRequest != nil {
		s.Hooks.OnRequest(RequestEvent{
			Conn:           conn,
			RequestHeader:  req.Header(),
			ResponseHeader: res.Header(),
			RequestBody:    req.Body(),
			ResponseBody:   res.Body(),
			Err:            cmp.Or(sendErr, handlerErr),
			Latency:        time.Since(startTime),
			Dropped:        req.Dropped(),
		})
	}

	// A GSF error is sent to the client as the error response
	// and keeps the connection, any other handler error closes it.
	if sendErr != nil || errors.As(handlerErr, new(wrap.GSFError)) {
		return sendErr
	}
	return handlerErr
}

// respond sends the error response for a GSF error, or the response body if the handler wrote one.
func respond(conn *Conn, res *Response, handlerErr error) error {
	var gsfErr wrap.GSFError
	if errors.As(handlerErr, &gsfErr) {
		res.Header().ResultCode = int32(resultcode.ERR)
		res.Header().AppCode = int32(gsfErr.AppCode())
		res.Header().AppString = gsfErr.Error()
		return conn.Send(res.Header(), res.Body())
	}
	if handlerErr != nil || res.Body() == nil {
		return nil
	}
	return conn.Send(res.Header(), res.Body())
}

func (s *Server) readMessage(stream *bufio.ReadWriter) ([]byte, error) {
	length, err := s.Codec.ReadLength(stream)
	if err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, fmt.Errorf("empty message")
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(stream, data); err != nil {
		return nil, err
	}

	if data[len(data)-1] != 0 {
		return nil, fmt.Errorf("invalid message")
	}

	return data, nil
}
