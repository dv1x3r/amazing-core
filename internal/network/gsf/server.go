package gsf

import (
	"bufio"
	"bytes"
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
	OnConnect    func(remoteIP string)
	OnDisconnect func(remoteIP string, reason string)
	OnUnhandled  func(remoteIP string, header *Header, data []byte)
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
	if s.Hooks.OnConnect != nil {
		s.Hooks.OnConnect(remoteAddr)
	}

	for {
		err := s.processRequest(ctx, stream, remoteAddr)
		if err == nil {
			continue
		}

		if errors.Is(err, io.EOF) {
			if s.Hooks.OnDisconnect != nil {
				s.Hooks.OnDisconnect(remoteAddr, "eof")
			}
			break
		}

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			if s.Hooks.OnDisconnect != nil {
				s.Hooks.OnDisconnect(remoteAddr, "timeout")
			}
			break
		}

		if s.Hooks.OnDisconnect != nil {
			s.Hooks.OnDisconnect(remoteAddr, err.Error())
		}
		break
	}
}

func (s *Server) processRequest(ctx context.Context, stream *bufio.ReadWriter, remoteAddr string) error {
	data, err := s.readMessage(stream)
	if err != nil {
		return err
	}

	reader := s.Codec.NewReader(data)
	writer := s.Codec.NewWriter()

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
			s.Hooks.OnUnhandled(remoteAddr, header, data)
		}
		return nil
	}

	req := NewRequest(ctx, header, reader)
	res := NewResponse(header, writer)
	req.RemoteIP = remoteAddr

	if err = handler(res, req); err != nil {
		return err
	}

	if res.Body() != nil {
		return s.writeMessage(stream, writer)
	}

	return nil
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

func (s *Server) writeMessage(stream *bufio.ReadWriter, message ProtocolWriter) error {
	memoryStream := &bytes.Buffer{}
	message.CommitTo(memoryStream)

	if err := s.Codec.WriteLength(stream, memoryStream.Len()+1); err != nil {
		return err
	}

	stream.Write(memoryStream.Bytes())
	stream.WriteByte(0)
	return stream.Flush()
}
