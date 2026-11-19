package gsf

import "github.com/dv1x3r/amazing-core/internal/network/gsf/types/clientmessagetype"

type ResponseWriter interface {
	Header() *Header
	// Write sets the response body, the server sends it when the handler returns no error.
	// A handler that writes no body sends no response, e.g. for an inbound sync notify.
	// For a GSF error the server sends the error response instead.
	Write(body Serializable) error
	// Notify queues a notify to the client, the server sends it
	// after the response when the handler returns no error.
	Notify(msgType clientmessagetype.ClientMessageType, body Serializable)
}

type NotifyMessage struct {
	MsgType clientmessagetype.ClientMessageType
	Body    Serializable
}

type Response struct {
	header   *Header
	body     Serializable
	notifies []NotifyMessage
}

func NewResponse(header *Header) *Response {
	return &Response{
		header: &Header{
			Flags:         header.Flags | FlagResponse,
			SvcClass:      header.SvcClass,
			MsgType:       header.MsgType,
			RequestID:     header.RequestID,
			LogCorrelator: header.LogCorrelator,
		},
	}
}

func (res *Response) Header() *Header {
	return res.header
}

func (res *Response) Body() Serializable {
	return res.body
}

func (res *Response) Write(body Serializable) error {
	res.body = body
	return nil
}

func (res *Response) Notify(msgType clientmessagetype.ClientMessageType, body Serializable) {
	res.notifies = append(res.notifies, NotifyMessage{MsgType: msgType, Body: body})
}
