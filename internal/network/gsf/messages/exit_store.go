package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// ExitStoreRequest contains the closed store.
type ExitStoreRequest struct {
	StoreOID types.OID
}

func (req *ExitStoreRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.StoreOID)
}

// ExitStoreResponse is empty.
type ExitStoreResponse struct{}

func (res *ExitStoreResponse) Serialize(writer gsf.ProtocolWriter) {}
