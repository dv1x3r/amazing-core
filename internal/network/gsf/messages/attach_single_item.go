package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// AttachSingleItemRequest puts an item on a mount of a parent item, a door or a window on a house.
type AttachSingleItemRequest struct {
	ParentOID types.OID
	ChildOID  types.OID
	// Ordinal is the USDecoration.ordinal of the mount.
	Ordinal int32
}

func (req *AttachSingleItemRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.ParentOID)
	reader.ReadObject(&req.ChildOID)
	req.Ordinal = reader.ReadInt32()
}

// AttachSingleItemResponse repeats the attached item, the client does not read it.
type AttachSingleItemResponse struct {
	ParentOID  types.OID
	ChildOID   types.OID
	Ordinal    int32
	OrdinalUpd int32
}

func (res *AttachSingleItemResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.ParentOID)
	writer.WriteObject(&res.ChildOID)
	writer.WriteInt32(res.Ordinal)
	writer.WriteInt32(res.OrdinalUpd)
}
