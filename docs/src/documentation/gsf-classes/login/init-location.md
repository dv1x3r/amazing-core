# InitLocation

## Description

Sent during the `ZoneTransitionQueue` queue, after zones are loaded.

Returns the player's home location data and the address of the SYNC server.

The `Home.PlayerMaze.HomeTheme.AssetMap["Scene_Unity3D"]` asset drives scene loading via `LoadMazeCommand` -> `AssetDownloadManager.LoadMainScene()`.

`Home.PlayerMaze.HomeTheme.AssetMap` must contain a `"Scene_Unity3D"` entry used by `LoadMazeCommand` to load the main scene asset,
e.g. `Springtime003.unity3d` (main map) or `HomeLotSmall.unity3d` (home lot). (?)

`ZoneManager.InitLocationHomeResponseHandler()` stores `syncServerIP`, `syncServerToken`, and `syncServerPort` in `GameSettings`.
These are later used by `SyncManager` to open its session and send `SyncLogin`.

## Request

| Field   | Type                     | Description |
| ------- | ------------------------ | ----------- |
| `LocID` | [`OID`](../types/oid.md) |             |

## Response

| Field             | Type                                        | Description                                     |
| ----------------- | ------------------------------------------- | ----------------------------------------------- |
| `ZoneInstance`    | [`ZoneInstance`](../types/zone-instance.md) |                                                 |
| `Village`         | [`Village`](../types/village.md)            |                                                 |
| `Home`            | [`PlayerHome`](../types/player-home.md)     |                                                 |
| `SyncServerToken` | `string`                                    | Token used to authenticate with the SYNC server |
| `SyncServerIP`    | `string`                                    | IP address of the SYNC server                   |
| `SyncServerPort`  | `int32`                                     | Port of the SYNC server                         |
