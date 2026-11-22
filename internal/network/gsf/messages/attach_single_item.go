package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// AttachSingleItemRequest contains the item to put on a mount of the parent item.
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

// AttachSingleItemResponse contains the attached item and its mount.
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
