# GetZones

## Description

Sent during the `ZoneTransitionQueue`.

Fetches the list of available zones. The list is stored in `ZoneManager.Instance.zones`.

`LoadNPCsCommand` later calls `SpawnPoints.Instance.ParseZone(NPCManager.HardCodedZoneId)`, so the NPC zone must be present in this list.

The array must include the NPC zone with OID `{ Class: 4, Type: 16, Server: 0, Number: 2937912 }` (`NPCManager.HardCodedZoneId`). (?)

## Request

No fields.

## Response

| Field   | Type                         | Description   |
| ------- | ---------------------------- | ------------- |
| `Zones` | [`[]Zone`](../types/zone.md) | List of zones |
