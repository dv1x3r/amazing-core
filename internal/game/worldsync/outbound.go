package worldsync

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/clientmessagetype"
)

// OutboundNotify is one prepared notify send.
type OutboundNotify struct {
	Session *gsf.Session
	MsgType clientmessagetype.ClientMessageType
	Body    gsf.Serializable
}

// NewOutboundNotify prepares one notify send after hub state is unlocked.
func NewOutboundNotify(session *gsf.Session, msgType clientmessagetype.ClientMessageType, body gsf.Serializable) OutboundNotify {
	return OutboundNotify{
		Session: session,
		MsgType: msgType,
		Body:    body,
	}
}

// SendAll sends prepared notifies in order.
func SendAll(notifies []OutboundNotify) error {
	for _, notify := range notifies {
		if err := notify.Session.SendNotify(notify.MsgType, notify.Body); err != nil {
			return err
		}
	}
	return nil
}
