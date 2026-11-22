package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// MazePiece is a maze piece with its assets, build time and category.
type MazePiece struct {
	AssetContainer
	BuildTime        int32
	IsSingleUse      bool
	PieceCategoryOID OID
	Ordinal          int32
}

func (mp *MazePiece) Serialize(writer gsf.ProtocolWriter) {
	mp.AssetContainer.Serialize(writer)
	writer.WriteInt32(mp.BuildTime)
	writer.WriteBool(mp.IsSingleUse)
	writer.WriteObject(&mp.PieceCategoryOID)
	writer.WriteInt32(mp.Ordinal)
}

func (mp *MazePiece) Deserialize(reader gsf.ProtocolReader) {
	mp.AssetContainer.Deserialize(reader)
	mp.BuildTime = reader.ReadInt32()
	mp.IsSingleUse = reader.ReadBool()
	reader.ReadObject(&mp.PieceCategoryOID)
	mp.Ordinal = reader.ReadInt32()
}
