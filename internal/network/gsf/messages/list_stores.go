package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// ListStoresRequest requests the stores of a zone or of an NPC.
type ListStoresRequest struct {
	ZoneOID types.OID
	NPCOID  types.OID
}

func (req *ListStoresRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.ZoneOID)
	reader.ReadObject(&req.NPCOID)
}

// ListStoresResponse contains the stores, the client takes the first one.
type ListStoresResponse struct {
	Stores []types.Store
}

func (res *ListStoresResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Stores, func(value types.Store) {
		writer.WriteObject(&value)
	})
}
