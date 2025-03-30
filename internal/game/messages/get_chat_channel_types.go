package messages

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types"
)

type GetChatChannelTypesRequest struct {
}

func (req *GetChatChannelTypesRequest) Deserialize(reader gsf.ProtocolReader) {
}

type GetChatChannelTypesResponse struct {
	ChatChannelTypes []types.ChatChannelType
}

func (res *GetChatChannelTypesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.ChatChannelTypes, func(value types.ChatChannelType) {
		writer.WriteObject(&value)
	})
}
