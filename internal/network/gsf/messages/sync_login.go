package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// SyncLoginRequest authenticates the player with the SYNC server.
type SyncLoginRequest struct {
	UID types.OID

	// Sync server token from the InitLocation response.
	Token string

	MaxVisSize int32
}

func (req *SyncLoginRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.UID)
	req.Token = reader.ReadString()
	req.MaxVisSize = reader.ReadInt32()
}

// SyncLoginResponse acknowledges SYNC server authentication.
type SyncLoginResponse struct {
}

func (res *SyncLoginResponse) Serialize(writer gsf.ProtocolWriter) {
}
