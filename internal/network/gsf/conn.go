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

type Conn struct {
	remoteIP string

	stream   *bufio.ReadWriter
	codec    ProtocolCodec
	onNotify func(NotifyEvent)
	writeMu  sync.Mutex
}

func NewConn(remoteIP string, stream *bufio.ReadWriter, codec ProtocolCodec, onNotify func(NotifyEvent)) *Conn {
	return &Conn{
		remoteIP: remoteIP,
		stream:   stream,
		codec:    codec,
		onNotify: onNotify,
	}
}

func (c *Conn) RemoteIP() string {
	return c.remoteIP
}

func (c *Conn) Send(header *Header, body Serializable) error {
	writer := c.codec.NewWriter()
	if err := wrap.Panic(func() error {
		writer.WriteBool(false)
		writer.WriteObject(header)
		writer.WriteObject(body)
		return nil
	}); err != nil {
		return err
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return c.writeMessage(writer)
}

func (c *Conn) SendNotify(msgType clientmessagetype.ClientMessageType, body Serializable) error {
	header := &Header{
		Flags:    2,
		SvcClass: int32(serviceclass.CLIENT),
		MsgType:  int32(msgType),
	}

	startTime := time.Now()
	err := c.Send(header, body)

	if c.onNotify != nil {
		c.onNotify(NotifyEvent{
			Conn:    c,
			Header:  header,
			Body:    body,
			Err:     err,
			Latency: time.Since(startTime),
		})
	}

	return err
}

func (c *Conn) writeMessage(message ProtocolWriter) error {
	memoryStream := &bytes.Buffer{}
	message.CommitTo(memoryStream)
	if err := c.codec.WriteLength(c.stream, memoryStream.Len()+1); err != nil {
		return err
	}
	if _, err := c.stream.Write(memoryStream.Bytes()); err != nil {
		return err
	}
	if err := c.stream.WriteByte(0); err != nil {
		return err
	}
	return c.stream.Flush()
}
