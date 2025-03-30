package gsf

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/dv1x3r/amazing-core/internal/lib/logger"
)

type Server struct {
	Addr   string
	Router Router
	Codec  ProtocolCodec

	listener net.Listener
	cancel   context.CancelFunc
	conns    sync.Map
	wg       sync.WaitGroup
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
		if errors.Is(err, net.ErrClosed) {
			return nil
		}

		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		s.conns.Store(conn, struct{}{})
		s.wg.Add(1)

		go func() {
			defer s.wg.Done()
			defer s.conns.Delete(conn)
			defer conn.Close()
			s.handleConn(ctx, conn)
		}()
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
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
	logger.Get().Info(fmt.Sprintf("[gsf] %s connected", remoteAddr))

	for {
		err := s.processRequest(ctx, stream, remoteAddr)
		if err == nil {
			continue
		}

		if errors.Is(err, io.EOF) {
			logger.Get().Info(fmt.Sprintf("[gsf] %s disconnected", remoteAddr))
			break
		}

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			logger.Get().Warn(fmt.Sprintf("[gsf] %s disconnected: timeout", remoteAddr))
			break
		}

		logger.Get().Error(fmt.Sprintf("[gsf] %s disconnected: "+err.Error(), remoteAddr))
		break
	}
}

func (s *Server) processRequest(ctx context.Context, stream *bufio.ReadWriter, remoteAddr string) error {
	data, err := s.readMessage(stream)
	if err != nil {
		return err
	}

	codecReader := s.Codec.NewReader(data)
	codecWriter := s.Codec.NewWriter()

	requestHeader, err := ReadHeader(codecReader)
	if err != nil {
		return err
	}

	handler, ok := s.Router.Lookup(requestHeader.SvcClass, requestHeader.MsgType)
	if !ok {
		logger.Get().Warn(fmt.Sprintf("[gsf] %s Unhandled: %+v", remoteAddr, requestHeader),
			slog.Any("hex", fmt.Sprintf("%x", data)),
		)
		return nil
	}

	req := NewRequest(ctx, requestHeader, codecReader)
	res := NewResponse(requestHeader, codecWriter)
	req.RemoteAddr = remoteAddr

	if err = handler(res, req); err != nil {
		return err
	}

	if res.Body() != nil {
		return s.writeMessage(stream, codecWriter)
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
