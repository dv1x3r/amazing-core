package messages

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types"
)

type SyncLoginRequest struct {
	UID        types.OID
	Token      string
	MaxVisSize int32
}

func (req *SyncLoginRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.UID)
	req.Token = reader.ReadString()
	req.MaxVisSize = reader.ReadInt32()
}

type SyncLoginResponse struct {
}

func (res *SyncLoginResponse) Serialize(writer gsf.ProtocolWriter) {
}
