package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// DetachItemsRequest contains the parent item and the items to take off its mounts.
type DetachItemsRequest struct {
	ParentOID  types.OID
	ChildItems []types.OID
}

func (req *DetachItemsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.ParentOID)
	req.ChildItems = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
}

// DetachItemsResponse is empty.
type DetachItemsResponse struct{}

func (res *DetachItemsResponse) Serialize(writer gsf.ProtocolWriter) {}
