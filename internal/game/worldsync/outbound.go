package worldsync

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/clientmessagetype"
)

// OutboundNotify is one prepared notify send.
type OutboundNotify struct {
	Conn    *gsf.Conn
	MsgType clientmessagetype.ClientMessageType
	Body    gsf.Serializable
}

// NewOutboundNotify prepares one notify send after hub state is unlocked.
func NewOutboundNotify(conn *gsf.Conn, msgType clientmessagetype.ClientMessageType, body gsf.Serializable) OutboundNotify {
	return OutboundNotify{
		Conn:    conn,
		MsgType: msgType,
		Body:    body,
	}
}

// SendAll sends prepared notifies in order.
func SendAll(notifies []OutboundNotify) error {
	for _, notify := range notifies {
		if err := notify.Conn.SendNotify(notify.MsgType, notify.Body); err != nil {
			return err
		}
	}
	return nil
}
