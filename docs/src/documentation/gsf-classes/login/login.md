# Login

## Description

Sent when the player submits the login form.

## Request

| Field                  | Type                                                           | Description                                            |
| ---------------------- | -------------------------------------------------------------- | ------------------------------------------------------ |
| `LoginID`              | `string`                                                       | Player username                                        |
| `Password`             | `string`                                                       | Player password                                        |
| `SitePIN`              | `int32`                                                        | Hardcoded to `1234`                                    |
| `LanguageLocalePairID` | `OID`                                                          | "Class": 4, "Type": 19, "Server": 0, "Number": 9023265 |
| `UserQueueingToken`    | `string`                                                       | Hardcoded to `"Token"`                                 |
| `ClientEnvInfo`        | [`ClientEnvironmentData`](../types/client-environment-data.md) | OS, resolution, Unity version, etc.                    |
| `Token`                | `string`                                                       |                                                        |
| `LoginType`            | `int32`                                                        |                                                        |
| `CNL`                  | `string`                                                       | Channel: `steam`                                       |

## Response

| Field                     | Type                                          | Description                       |
| ------------------------- | --------------------------------------------- | --------------------------------- |
| `SiteInfo`                | [`SiteInfo`](../types/site-info.md)           |                                   |
| `Status`                  | [`SessionStatus`](../types/session-status.md) |                                   |
| `SessionID`               | `OID`                                         |                                   |
| `ConversationID`          | `int64`                                       |                                   |
| `AssetDeliveryURL`        | `string`                                      | Base URL used for asset downloads |
| `Player`                  | [`Player`](../types/player.md)                |                                   |
| `MaxOutfit`               | `int16`                                       |                                   |
| `PlayerStats`             | [`[]PlayerStats`](../types/player-stats.md)   |                                   |
| `PlayerInfoTO`            | [`PlayerInfoTO`](../types/player-info-to.md)  |                                   |
| `CurrentServerTime`       | `time.Time`                                   |                                   |
| `SystemLockoutTime`       | `time.Time`                                   |                                   |
| `SystemShutdownTime`      | `time.Time`                                   |                                   |
| `ClientInactivityTimeout` | `int32`                                       |                                   |
| `CNL`                     | `string`                                      |                                   |
