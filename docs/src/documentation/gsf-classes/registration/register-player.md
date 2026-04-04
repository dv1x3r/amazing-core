# RegisterPlayer

## Description

Sent when the player submits the registration form.

Creates the player account on the server.

On success, the returned `PlayerID` is immediately used by the client to call [`RegisterAvatarForRegistration`](./register-avatar-for-registration.md),
which associates the chosen avatar and starter items with the new account.

## Request

| Field                 | Type                     | Description                                     |
| --------------------- | ------------------------ | ----------------------------------------------- |
| `Token`               | `string`                 |                                                 |
| `Password`            | `string`                 | Chosen password                                 |
| `ParentEmailAddress`  | `string`                 | Chosen email                                    |
| `BirthDate`           | `time.Time`              |                                                 |
| `Gender`              | `string`                 | Player gender (always `"U"` for unspecified)    |
| `LocationID`          | [`OID`](../types/oid.md) | "Class": 0, "Type": 0, "Server": 0, "Number": 0 |
| `Username`            | `string`                 | Chosen username                                 |
| `Worldname`           | `string`                 | Chosen family name                              |
| `ChatAllowed`         | `bool`                   | Whether the player has enabled chat             |
| `CNL`                 | `string`                 | Channel: `steam`                                |
| `ReferredByWorldname` | `string`                 |                                                 |
| `LoginType`           | `int32`                  |                                                 |

## Response

| Field      | Type                     | Description                    |
| ---------- | ------------------------ | ------------------------------ |
| `PlayerID` | [`OID`](../types/oid.md) | The newly created player's OID |
