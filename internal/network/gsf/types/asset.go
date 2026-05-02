package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// Asset is a single resource file, such as Unity asset bundle, treenode, json, image, or audio file.
type Asset struct {
	// OID identifier.
	OID OID

	// Asset classification (e.g. "Prefab_Unity3D").
	AssetTypeName string

	// CDN identifier used to construct the download URL: AssetDeliveryURL + CDNID.
	CDNID string

	// Resource name.
	ResName string

	// Group classification (e.g. "Main_Scene", "3D Components", "Locked").
	GroupName string

	// File size in bytes.
	FileSize int64
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
