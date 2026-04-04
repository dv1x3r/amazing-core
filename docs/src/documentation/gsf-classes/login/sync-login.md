# SyncLogin

## Description

Sent when the SYNC server opens TCP session.

The SYNC session is started by `SyncManager.Instance.Start()`, which is called from `LoadServicesCommand.cs:21`.

The sync server address and token come from the `InitLocation` response stored in `GameSettings`.

Authenticates the player on the SYNC server. The SYNC server handles real-time positional and social updates.

## Request

| Field        | Type                     | Description                                        |
| ------------ | ------------------------ | -------------------------------------------------- |
| `UID`        | [`OID`](../types/oid.md) |                                                    |
| `Token`      | `string`                 | Sync server token from the `InitLocation` response |
| `MaxVisSize` | `int32`                  |                                                    |

## Response

No fields.
