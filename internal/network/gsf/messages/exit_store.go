package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// ExitStoreRequest is sent when the store window closes.
type ExitStoreRequest struct {
	StoreOID types.OID
}

func (req *ExitStoreRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.StoreOID)
}

// ExitStoreResponse acknowledges the store exit.
type ExitStoreResponse struct{}

func (res *ExitStoreResponse) Serialize(writer gsf.ProtocolWriter) {}
