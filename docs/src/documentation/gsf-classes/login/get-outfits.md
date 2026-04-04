# GetOutfits

## Description

Sent during the `ZoneTransitionQueue` after avatars are loaded, once per avatar whose assets are being loaded.

Fetches the saved outfits for a given player avatar. The results are stored as `PresetOutfits` on the `AvatarAssets` object.

## Request

| Field            | Type                     | Description          |
| ---------------- | ------------------------ | -------------------- |
| `PlayerAvatarID` | [`OID`](../types/oid.md) | Requested avatar OID |
| `PlayerID`       | [`OID`](../types/oid.md) | Requested player OID |

## Response

| Field                 | Type                                                       | Description            |
| --------------------- | ---------------------------------------------------------- | ---------------------- |
| `PlayerAvatarOutfits` | [`[]PlayerAvatarOutfit`](../types/player-avatar-outfit.md) | List of avatar outfits |
