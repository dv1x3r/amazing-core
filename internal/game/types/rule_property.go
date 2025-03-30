package types

import (
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

type RuleProperty struct {
	ID               OID
	ParentID         OID
	Components       []string
	ParentComponents []string
	CreateTime       time.Time
	ModifiedTime     time.Time
	Properties       map[string]string
	ChildrenGroup    map[string][]RuleProperty
}

func (rp *RuleProperty) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&rp.ID)
	writer.WriteObject(&rp.ParentID)
	gsf.WriteSlice(writer, rp.Components, writer.WriteString)
	gsf.WriteSlice(writer, rp.ParentComponents, writer.WriteString)
	writer.WriteUtcDate(rp.CreateTime)
	writer.WriteUtcDate(rp.ModifiedTime)
	gsf.WriteMap(writer, rp.Properties, writer.WriteString)
	gsf.WriteMap(writer, rp.ChildrenGroup, func(slice []RuleProperty) {
		gsf.WriteSlice(writer, slice, func(value RuleProperty) {
			writer.WriteObject(&value)
		})
	})
}

func (rp *RuleProperty) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&rp.ID)
	reader.ReadObject(&rp.ParentID)
	rp.Components = gsf.ReadSlice(reader, reader.ReadString)
	rp.ParentComponents = gsf.ReadSlice(reader, reader.ReadString)
	rp.CreateTime = reader.ReadUtcDate()
	rp.ModifiedTime = reader.ReadUtcDate()
	rp.Properties = gsf.ReadMap(reader, reader.ReadString)
	rp.ChildrenGroup = gsf.ReadMap(reader, func() []RuleProperty {
		return gsf.ReadSlice(reader, func() RuleProperty {
			var value RuleProperty
			reader.ReadObject(&value)
			return value
		})
	})
}
