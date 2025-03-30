package messages

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types"
)

type GetAvatarsRequest struct {
	Start     int32
	Max       int32
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

type GetAvatarsResponse struct {
	Avatars []types.PlayerAvatar
}

func (res *GetAvatarsResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.Avatars, func(value types.PlayerAvatar) {
		writer.WriteObject(&value)
	})
}
