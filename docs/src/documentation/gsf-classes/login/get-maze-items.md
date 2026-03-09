# GetMazeItems

## Description

Sent during the `ZoneTransitionQueue`, after the home location is initialised.

Fetches the items placed in the player's current maze.

## Request

| Field          | Type  | Description          |
| -------------- | ----- | -------------------- |
| `PlayerMazeID` | `OID` | Requested maze OID   |
| `PlayerID`     | `OID` | Requested player OID |

## Response

| Field       | Type                                      | Description          |
| ----------- | ----------------------------------------- | -------------------- |
| `MazeItems` | [`[]PlayerItem`](../types/player-item.md) | List of player items |
