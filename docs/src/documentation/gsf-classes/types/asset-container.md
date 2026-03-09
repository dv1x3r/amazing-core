# AssetContainer

## Description

A container that groups assets by type and bundles optional asset packages.

## Fields

| Field           | Type                                   | Description                                                       |
| --------------- | -------------------------------------- | ----------------------------------------------------------------- |
| `OID`           | `OID`                                  | Container identifier                                              |
| `AssetMap`      | [`map[string][]Asset`](./asset.md)     | Dictionary keyed by `AssetTypeName` containing a list of assets   |
| `AssetPackages` | [`[]AssetPackage`](./asset-package.md) | List of asset packages for conditional/maze-specific asset groups |
