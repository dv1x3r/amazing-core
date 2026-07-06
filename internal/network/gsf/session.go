package gsf

import (
	"bufio"
	"bytes"
	"sync"
	"time"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/clientmessagetype"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/serviceclass"
)

type Session struct {
	remoteIP     string
	platform     Platform
	playerOID    int64
	hasPlayerOID bool

	stream   *bufio.ReadWriter
	codec    ProtocolCodec
	onNotify func(NotifyEvent)
	stateMu  sync.RWMutex
	writeMu  sync.Mutex
}

func NewSession(remoteIP string, stream *bufio.ReadWriter, codec ProtocolCodec, onNotify func(NotifyEvent)) *Session {
	return &Session{
		remoteIP: remoteIP,
		stream:   stream,
		codec:    codec,
		onNotify: onNotify,
	}
}

func (s *Session) RemoteIP() string {
	return s.remoteIP
}

func (s *Session) Platform() Platform {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.platform
}

func (s *Session) SetPlatform(platform Platform) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	s.platform = platform
}

func (s *Session) PlayerOID() (int64, bool) {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.playerOID, s.hasPlayerOID
}

func (s *Session) SetPlayerOID(playerOID int64) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	s.playerOID = playerOID
	s.hasPlayerOID = true
}

func (s *Session) Send(header *Header, body Serializable) error {
	writer := s.codec.NewWriter()
	if err := wrap.Panic(func() error {
		writer.WriteBool(false)
		writer.WriteObject(header)
		writer.WriteObject(body)
		return nil
	}); err != nil {
		return err
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	return s.writeMessage(writer)
}

func (s *Session) SendNotify(msgType clientmessagetype.ClientMessageType, body Serializable) error {
	header := &Header{
		Flags:    2,
		SvcClass: int32(serviceclass.CLIENT),
		MsgType:  int32(msgType),
	}

	startTime := time.Now()
	err := s.Send(header, body)

	if s.onNotify != nil {
		s.onNotify(NotifyEvent{
			Session: s,
			Header:  header,
			Body:    body,
			Err:     err,
			Latency: time.Since(startTime),
		})
	}

	return err
}

func (s *Session) writeMessage(message ProtocolWriter) error {
	memoryStream := &bytes.Buffer{}
	message.CommitTo(memoryStream)
	if err := s.codec.WriteLength(s.stream, memoryStream.Len()+1); err != nil {
		return err
	}
	if _, err := s.stream.Write(memoryStream.Bytes()); err != nil {
		return err
	}
	if err := s.stream.WriteByte(0); err != nil {
		return err
	}
	return s.stream.Flush()
}
