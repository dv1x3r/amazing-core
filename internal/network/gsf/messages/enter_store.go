package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// EnterStoreRequest is sent when the store window opens.
type EnterStoreRequest struct {
	StoreOID types.OID
}

func (req *EnterStoreRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.StoreOID)
}

// EnterStoreResponse acknowledges the store enter.
type EnterStoreResponse struct{}

func (res *EnterStoreResponse) Serialize(writer gsf.ProtocolWriter) {}
