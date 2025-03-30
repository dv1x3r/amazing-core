package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type RuleContainer struct {
	AssetContainer
	RuleProperty    RuleProperty
	Locked          bool
	IsMultiplayer   bool
	IsPlayerHosted  bool
	IsPlayedOffline bool
}

func (rc *RuleContainer) Serialize(writer gsf.ProtocolWriter) {
	rc.AssetContainer.Serialize(writer)
	writer.WriteObject(&rc.RuleProperty)
	writer.WriteBool(rc.Locked)
	writer.WriteBool(rc.IsMultiplayer)
	writer.WriteBool(rc.IsPlayerHosted)
	writer.WriteBool(rc.IsPlayedOffline)
}

func (rc *RuleContainer) Deserialize(reader gsf.ProtocolReader) {
	rc.AssetContainer.Deserialize(reader)
	reader.ReadObject(&rc.RuleProperty)
	rc.Locked = reader.ReadBool()
	rc.IsMultiplayer = reader.ReadBool()
	rc.IsPlayerHosted = reader.ReadBool()
	rc.IsPlayedOffline = reader.ReadBool()
}
