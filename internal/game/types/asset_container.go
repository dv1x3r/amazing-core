package types

import "github.com/dv1x3r/amazing-core/internal/game/gsf"

type AssetContainer struct {
	OID           OID
	AssetMap      map[string][]Asset
	AssetPackages []AssetPackage
}

func (ac *AssetContainer) Serialize(writer gsf.ProtocolWriter) {
	writer.WriteObject(&ac.OID)
	gsf.WriteMap(writer, ac.AssetMap, func(slice []Asset) {
		gsf.WriteSlice(writer, slice, func(value Asset) {
			writer.WriteObject(&value)
		})
	})
	gsf.WriteSlice(writer, ac.AssetPackages, func(value AssetPackage) {
		writer.WriteObject(&value)
	})
}

func (ac *AssetContainer) Deserialize(reader gsf.ProtocolReader) {
	reader.ReadObject(&ac.OID)
	ac.AssetMap = gsf.ReadMap(reader, func() []Asset {
		return gsf.ReadSlice(reader, func() Asset {
			var value Asset
			reader.ReadObject(&value)
			return value
		})
	})
	ac.AssetPackages = gsf.ReadSlice(reader, func() AssetPackage {
		var value AssetPackage
		reader.ReadObject(&value)
		return value
	})
}
