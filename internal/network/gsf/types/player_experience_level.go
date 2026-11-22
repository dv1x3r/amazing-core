package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// PlayerExperienceLevel is a level with its experience and play time range.
type PlayerExperienceLevel struct {
	AssetContainer
	Level              int16
	MinExperienceValue int32
	MaxExperienceValue int32
	MinPlayTimeMinutes int32
	MaxPlayTimeMinutes int32
}

func (l *PlayerExperienceLevel) Serialize(writer gsf.ProtocolWriter) {
	l.AssetContainer.Serialize(writer)
	writer.WriteInt16(l.Level)
	writer.WriteInt32(l.MinExperienceValue)
	writer.WriteInt32(l.MaxExperienceValue)
	writer.WriteInt32(l.MinPlayTimeMinutes)
	writer.WriteInt32(l.MaxPlayTimeMinutes)
}

func (l *PlayerExperienceLevel) Deserialize(reader gsf.ProtocolReader) {
	l.AssetContainer.Deserialize(reader)
	l.Level = reader.ReadInt16()
	l.MinExperienceValue = reader.ReadInt32()
	l.MaxExperienceValue = reader.ReadInt32()
	l.MinPlayTimeMinutes = reader.ReadInt32()
	l.MaxPlayTimeMinutes = reader.ReadInt32()
}
