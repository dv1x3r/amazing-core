package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type RaceMode struct {
	OID  OID
	Name string
}

func (rm *RaceMode) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&rm.OID)
	writer.WriteString(rm.Name)
}

func (rm *RaceMode) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&rm.OID)
	rm.Name = reader.ReadString()
}
