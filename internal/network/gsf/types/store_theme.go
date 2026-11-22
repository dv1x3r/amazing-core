package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// StoreTheme is a World Shop category, or a subcategory when it has a parent.
type StoreTheme struct {
	RuleContainer
	// ParentStoreThemeOID is nil for a category, WorldShopSVCS checks parentStoreThemeId == null.
	ParentStoreThemeOID *OID
	Ordinal             int32
}

func (st *StoreTheme) Serialize(writer gsf.ProtocolWriter) {
	st.RuleContainer.Serialize(writer)
	writer.WriteObject(st.ParentStoreThemeOID)
	writer.WriteInt32(st.Ordinal)
}

func (st *StoreTheme) Deserialize(reader gsf.ProtocolReader) {
	st.RuleContainer.Deserialize(reader)
	var parent OID
	if reader.ReadObject(&parent) {
		st.ParentStoreThemeOID = &parent
	}
	st.Ordinal = reader.ReadInt32()
}
