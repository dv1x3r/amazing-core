package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// Building is a rule-backed building object within a zone.
type Building struct {
	RuleContainer
	Position   Position
	Dimensions Dimensions
	SpawnPoint string
	ZoneOID    OID
}

func (b *Building) Serialize(writer gsf.ProtocolWriter) {
	b.RuleContainer.Serialize(writer)
	writer.WriteObject(&b.Position)
	writer.WriteObject(&b.Dimensions)
	writer.WriteString(b.SpawnPoint)
	writer.WriteObject(&b.ZoneOID)
}

func (b *Building) Deserialize(reader gsf.ProtocolReader) {
	b.RuleContainer.Deserialize(reader)
	reader.ReadObject(&b.Position)
	reader.ReadObject(&b.Dimensions)
	b.SpawnPoint = reader.ReadString()
	reader.ReadObject(&b.ZoneOID)
}
