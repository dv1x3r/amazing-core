# GetPlayerNPCs

## Description

Sent during the `ZoneTransitionQueue`, as the last major data fetch.

Fetches the NPCs for the player's home zone.

`SpawnPoints.Instance.ParseZone(NPCManager.HardCodedZoneId)` is called beforehand to set up spawn points - this requires the NPC zone to have been returned by `GetZones`.

## Request

| r Field    | Type                     | Description          |
| ---------- | ------------------------ | -------------------- |
| `PlayerID` | [`OID`](../types/oid.md) | Requested player OID |
| `ZoneID`   | [`OID`](../types/oid.md) | Requested zone OID   |

## Response

| Field  | Type                       | Description  |
| ------ | -------------------------- | ------------ |
| `NPCs` | [`[]NPC`](../types/npc.md) | List of NPCs |
