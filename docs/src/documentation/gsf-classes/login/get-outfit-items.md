# GetOutfitItems

## Description

Sent as the first step of the `ZoneTransitionQueue`.

Fetches the items currently equipped(?) by the player's active avatar.

## Request

| Field                  | Type                     | Description                 |
| ---------------------- | ------------------------ | --------------------------- |
| `PlayerAvatarOutfitID` | [`OID`](../types/oid.md) | Requested avatar outfit OID |
| `PlayerID`             | [`OID`](../types/oid.md) | Requested player OID        |

## Response

| Field         | Type                                      | Description          |
| ------------- | ----------------------------------------- | -------------------- |
| `OutfitItems` | [`[]PlayerItem`](../types/player-item.md) | List of player items |
