package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetChatChannelTypesRequest requests available chat channel definitions.
type GetChatChannelTypesRequest struct {
}

func (req *GetChatChannelTypesRequest) Deserialize(reader gsf.ProtocolReader) {
}

// GetChatChannelTypesResponse contains available chat channel definitions.
type GetChatChannelTypesResponse struct {
	ChatChannelTypes []types.ChatChannelType
}

func (res *GetChatChannelTypesResponse) Serialize(writer gsf.ProtocolWriter) {
	gsf.WriteSlice(writer, res.ChatChannelTypes, func(value types.ChatChannelType) {
		writer.WriteObject(&value)
	})
}
