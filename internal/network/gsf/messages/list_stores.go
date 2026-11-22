package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// ListStoresRequest contains the zone or the NPC of the stores.
type ListStoresRequest struct {
	ZoneOID types.OID
	NPCOID  types.OID
}

func (req *ListStoresRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.ZoneOID)
	reader.ReadObject(&req.NPCOID)
}

// ListStoresResponse contains the stores.
type ListStoresResponse struct {
	Stores []types.Store
}

func (res *ListStoresResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Stores, func(value types.Store) {
		writer.WriteObject(&value)
	})
}
