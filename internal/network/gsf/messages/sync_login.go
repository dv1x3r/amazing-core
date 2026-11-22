package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// SyncLoginRequest contains the player and the token for the SYNC server.
type SyncLoginRequest struct {
	UOID types.OID

	// Sync server token from the InitLocation response.
	Token string

	MaxVisSize int32
}

func (req *SyncLoginRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.UOID)
	req.Token = reader.ReadString()
	req.MaxVisSize = reader.ReadInt32()
}

// SyncLoginResponse is empty.
type SyncLoginResponse struct {
}

func (res *SyncLoginResponse) Serialize(writer gsf.ProtocolWriter) {
}
