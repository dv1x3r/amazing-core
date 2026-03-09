# EnterBuilding

## Description

Sent during the `ZoneTransitionQueue`, after announcements are loaded.

Notifies the server that the player is entering a building.

During initial login with no active building context, this is a pass-through step. (?)

## Request

| Field         | Type       | Description                     |
| ------------- | ---------- | ------------------------------- |
| `LocID`       | `OID`      | OID of the zone (?)             |
| `BuildingID`  | `OID`      | OID of the building             |
| `Pos`         | `Position` | Player position                 |
| `Orientation` | `QTH`      | Player orientation (quaternion) |

## Response

| Field        | Type  | Description                 |
| ------------ | ----- | --------------------------- |
| `BuildingID` | `OID` | OID of the building entered |
