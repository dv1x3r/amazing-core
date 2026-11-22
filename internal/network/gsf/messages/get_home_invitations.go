package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetHomeInvitationsRequest contains the player and the filter of the home invitations.
type GetHomeInvitationsRequest struct {
	PlayerOID       types.OID
	IsHost          bool
	IsStatusPending bool
}

func (req *GetHomeInvitationsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerOID)
	req.IsHost = reader.ReadBool()
	req.IsStatusPending = reader.ReadBool()
}

// GetHomeInvitationsResponse contains the home invitations, always an empty list.
type GetHomeInvitationsResponse struct{}

func (res *GetHomeInvitationsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(0)
}
