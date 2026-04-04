# AssetPackage

## Description

An [`AssetContainer`](./asset-container.md) with an additional tag used to conditionally load assets depending on which maze the player is entering.

## Fields

| Field           | Type                                   | Description                                                                        |
| --------------- | -------------------------------------- | ---------------------------------------------------------------------------------- |
| `OID`           | [`OID`](../types/oid.md)               | Inherited from [`AssetContainer`](./asset-container.md)                            |
| `AssetMap`      | [`map[string][]Asset`](./asset.md)     | Inherited from [`AssetContainer`](./asset-container.md)                            |
| `AssetPackages` | [`[]AssetPackage`](./asset-package.md) | Inherited from [`AssetContainer`](./asset-container.md)                            |
| `PTag`          | `string`                               | Tag matched against the current maze name to decide whether this package is loaded |
| `CreateDate`    | `time.Time`                            |                                                                                    |

## PTag matching

During `LoadMazeCommand:238`, the client iterates all packages from the home theme and applies this rule:

- `PTag == maze.name` -> package is loaded (maze-specific assets)
- `PTag == ""` (empty) -> package is always loaded (shared assets)
- `PTag` is anything else -> package is skipped

This allows the server to ship assets for multiple different mazes in a single response, while the client only downloads what is relevant to the current location. (?)
