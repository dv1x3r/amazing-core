package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetAvatarsRequest requests the player's avatar list.
type GetAvatarsRequest struct {
	// Start index, 0.
	Start int32

	// Max results, -1 for all.
	Max int32

	// Filter list.
	FilterIDs []types.OID
}

func (req *GetAvatarsRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Start = reader.ReadInt32()
	req.Max = reader.ReadInt32()
	req.FilterIDs = gsf.ReadSlice(reader, func() types.OID {
		var value types.OID
		reader.ReadObject(&value)
		return value
	})
}

// GetAvatarsResponse contains player avatar records.
type GetAvatarsResponse struct {
	Avatars []types.PlayerAvatar
}

func (res *GetAvatarsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Avatars, func(value types.PlayerAvatar) {
		writer.WriteObject(&value)
	})
}
