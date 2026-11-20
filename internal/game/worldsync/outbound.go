package worldsync

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/clientmessagetype"
)

// outboundNotify is a notify to send to a SYNC connection.
type outboundNotify struct {
	conn *gsf.Conn
	gsf.NotifyMessage
}

func newOutboundNotify(conn *gsf.Conn, msgType clientmessagetype.ClientMessageType, body gsf.Serializable) outboundNotify {
	return outboundNotify{
		conn: conn,
		NotifyMessage: gsf.NotifyMessage{
			MsgType: msgType,
			Body:    body,
		},
	}
}

// sendAll sends the notifies in order. A failed send is not the sender's error,
// the OnNotify hook logs it and the dead connection is closed by its own read loop.
func sendAll(notifies []outboundNotify) {
	for _, n := range notifies {
		_ = n.conn.SendNotify(n.MsgType, n.Body)
	}
}
