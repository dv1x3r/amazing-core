package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type ItemCategory struct {
	RuleContainer
	CreateDate time.Time
	IsOutdoor  bool
	IsWalkover bool
	ParentID   OID
	Name       string
	ShowInDock bool
}

func (ic *ItemCategory) Serialize(writer gsf.ProtocolWriter) {
	ic.RuleContainer.Serialize(writer)
	writer.WriteUtcDate(ic.CreateDate)
	writer.WriteBool(ic.IsOutdoor)
	writer.WriteBool(ic.IsWalkover)
	writer.WriteObject(&ic.ParentID)
	writer.WriteString(ic.Name)
	writer.WriteBool(ic.ShowInDock)
}

func (ic *ItemCategory) Deserialize(reader gsf.ProtocolReader) {
	ic.RuleContainer.Deserialize(reader)
	ic.CreateDate = reader.ReadUtcDate()
	ic.IsOutdoor = reader.ReadBool()
	ic.IsWalkover = reader.ReadBool()
	reader.ReadObject(&ic.ParentID)
	ic.Name = reader.ReadString()
	ic.ShowInDock = reader.ReadBool()
}
