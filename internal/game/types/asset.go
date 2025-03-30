package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type Asset struct {
	OID           OID
	AssetTypeName string
	CDNID         string
	ResName       string
	GroupName     string
	FileSize      int64
}

func (a *Asset) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&a.OID)
	writer.WriteString(a.AssetTypeName)
	writer.WriteString(a.CDNID)
	writer.WriteString(a.ResName)
	writer.WriteString(a.GroupName)
	writer.WriteInt64(a.FileSize)
}

func (a *Asset) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&a.OID)
	a.AssetTypeName = reader.ReadString()
	a.CDNID = reader.ReadString()
	a.ResName = reader.ReadString()
	a.GroupName = reader.ReadString()
	a.FileSize = reader.ReadInt64()
}
