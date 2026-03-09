# GetAvatars

## Description

Sent during the `ZoneTransitionQueue`, after outfit items are loaded.

Fetches the list of player avatars. The returned list is stored in `AvatarManager.Instance.GSFPlayerAvatars`.

`LoadAvatarsCommand.Step2()` calls `AvatarManager.LoadAvatarAssets()` for the active player avatar, which in turn calls `GetOutfits`.

The active player avatar must be present in the returned array. (?)

Each `PlayerAvatar.Avatar.AssetMap` must contain a `"Prefab_Unity3D"` entry with at least one asset whose `ResName` is not `"PF__Avatar.unity3d"`. (?)

## Request

| Field       | Type    | Description               |
| ----------- | ------- | ------------------------- |
| `Start`     | `int32` | Start index, `0`          |
| `Max`       | `int32` | Max results, `-1` for all |
| `FilterIDs` | `[]OID` | Empty filter list         |

## Response

| Field     | Type                                          | Description            |
| --------- | --------------------------------------------- | ---------------------- |
| `Avatars` | [`[]PlayerAvatar`](../types/player-avatar.md) | List of player avatars |
