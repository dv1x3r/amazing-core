package messages

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// GetRequiredExperienceRequest contains the level.
type GetRequiredExperienceRequest struct {
	Level int16
}

func (req *GetRequiredExperienceRequest) Deserialize(reader gsf.ProtocolReader) {
	req.Level = reader.ReadInt16()
}

// GetRequiredExperienceResponse contains the experience and play time range of the level.
type GetRequiredExperienceResponse struct {
	MinExperience         int32
	MaxExperience         int32
	MinTimePlayed         int32
	MaxTimePlayed         int32
	PlayerExperienceLevel types.PlayerExperienceLevel
}

func (res *GetRequiredExperienceResponse) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteInt32(res.MinExperience)
	writer.WriteInt32(res.MaxExperience)
	writer.WriteInt32(res.MinTimePlayed)
	writer.WriteInt32(res.MaxTimePlayed)
	writer.WriteObject(&res.PlayerExperienceLevel)
}
