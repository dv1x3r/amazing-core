package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetOtherPlayerDetailsRequest contains the other player.
type GetOtherPlayerDetailsRequest struct {
	PlayerOID types.OID
}

func (req *GetOtherPlayerDetailsRequest) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&req.PlayerOID)
}

// GetOtherPlayerDetailsResponse contains the Zing, the clothes and the profile of the other player.
type GetOtherPlayerDetailsResponse struct {
	OtherPlayerDetails types.OtherPlayerDetails
}

func (res *GetOtherPlayerDetailsResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&res.OtherPlayerDetails)
}
