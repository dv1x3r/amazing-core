package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

type StoreTheme struct {
	RuleContainer
	ParentStoreThemeOID OID
	Ordinal             int32
}

func (st *StoreTheme) Serialize(writer gsf.ProtocolWriter) {
	st.RuleContainer.Serialize(writer)
	writer.WriteObject(&st.ParentStoreThemeOID)
	writer.WriteInt32(st.Ordinal)
}

func (st *StoreTheme) Deserialize(reader gsf.ProtocolReader) {
	st.RuleContainer.Deserialize(reader)
	reader.ReadObject(&st.ParentStoreThemeOID)
	st.Ordinal = reader.ReadInt32()
}
