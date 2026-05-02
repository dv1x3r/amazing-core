package types

import "github.com/dv1x3r/amazing-core/internal/network/gsf"

// AssetMap maps asset type names to the assets in that group.
type AssetMap map[string][]Asset

// AssetContainer is a group of assets. It can contain multiple assets and/or asset packages.
type AssetContainer struct {
	// OID identifier.
	OID OID

	// Dictionary keyed by AssetTypeName containing a list of assets.
	AssetMap AssetMap

	// List of asset packages for conditional/maze-specific asset groups.
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
