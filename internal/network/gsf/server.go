package gsf

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
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
	OnConnect    func(session *Session)
	OnDisconnect func(session *Session, reason string)
	OnUnhandled  func(session *Session, header *Header, data []byte)
	OnRequest    func(event RequestEvent)
	OnNotify     func(event NotifyEvent)
}

type RequestEvent struct {
	Session        *Session
	RequestHeader  *Header
	ResponseHeader *Header
	RequestBody    Deserializable
	ResponseBody   Serializable
	Err            error
	Latency        time.Duration
}

type NotifyEvent struct {
	Session *Session
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
			} else {
				time.Sleep(1 * time.Second)
				continue
			}
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

func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	stream := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	remoteAddr, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
	session := NewSession(remoteAddr, stream, s.Codec, s.Hooks.OnNotify)
	if s.Hooks.OnConnect != nil {
		s.Hooks.OnConnect(session)
	}

	for {
		err := s.processRequest(ctx, session)
		if err == nil {
			continue
		}

		if errors.Is(err, io.EOF) {
			if s.Hooks.OnDisconnect != nil {
				s.Hooks.OnDisconnect(session, "eof")
			}
			break
		}

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			if s.Hooks.OnDisconnect != nil {
				s.Hooks.OnDisconnect(session, "timeout")
			}
			break
		}

		if s.Hooks.OnDisconnect != nil {
			s.Hooks.OnDisconnect(session, err.Error())
		}
		break
	}
}

func (s *Server) processRequest(ctx context.Context, session *Session) error {
	data, err := s.readMessage(session.stream)
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
			s.Hooks.OnUnhandled(session, header, data)
		}
		return nil
	}

	req := NewRequest(ctx, header, reader, session)
	res := NewResponse(header)

	startTime := time.Now()
	handlerErr := handler(res, req)
	eventErr := handlerErr
	returnErr := handlerErr

	if handlerErr == nil {
		if res.Body() != nil {
			eventErr = session.Send(res.Header(), res.Body())
			returnErr = eventErr
		}
	} else {
		var gsfErr wrap.GSFError
		if errors.As(handlerErr, &gsfErr) {
			res.Header().ResultCode = gsfErr.ResultCode()
			res.Header().AppCode = gsfErr.AppCode()
			res.Header().AppString = gsfErr.Error()
			if err := session.Send(res.Header(), res.Body()); err != nil {
				eventErr = err
				returnErr = err
			} else {
				returnErr = nil
			}
		}
	}

	if s.Hooks.OnRequest != nil {
		s.Hooks.OnRequest(RequestEvent{
			Session:        session,
			RequestHeader:  req.Header(),
			ResponseHeader: res.Header(),
			RequestBody:    req.Body(),
			ResponseBody:   res.Body(),
			Err:            eventErr,
			Latency:        time.Since(startTime),
		})
	}

	return returnErr
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
