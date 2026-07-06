package gsf

type ResponseWriter interface {
	Header() *Header
	Body() Serializable
	Write(Serializable) error
}

type Response struct {
	header *Header
	body   Serializable
}

func NewResponse(header *Header) *Response {
	return &Response{
		header: &Header{
			Flags:         header.Flags | 1,
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
