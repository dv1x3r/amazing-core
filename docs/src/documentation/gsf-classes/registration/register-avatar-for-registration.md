# RegisterAvatarForRegistration

## Description

Sent after a successful [`RegisterPlayer`](./register-player.md) response.

Associates the chosen avatar and starter inventory items with the newly created player account.

## Request

| Field               | Type     | Description                                       |
| ------------------- | -------- | ------------------------------------------------- |
| `PlayerID`          | `OID`    | The OID returned by `RegisterPlayer`              |
| `SecretCode`        | `string` |                                                   |
| `Name`              | `string` | Avatar name                                       |
| `Bio`               | `string` |                                                   |
| `AvatarID`          | `OID`    | OID of the selected avatar                        |
| `GivenInventoryIDs` | `[]OID`  | List of starter item OIDs selected during intro   |
| `GivenItemSlotIDs`  | `[]OID`  | List of inventory slot OIDs for the starter items |

## Response

| Field                  | Type                                        | Description                            |
| ---------------------- | ------------------------------------------- | -------------------------------------- |
| `PlayerAvatar`         | [`PlayerAvatar`](../types/player-avatar.md) | The created avatar object              |
| `InvalidCodeCount`     | `int32`                                     | Number of invalid secret code attempts |
| `InvalidCodeThreshold` | `int32`                                     | Maximum allowed invalid code attempts  |
