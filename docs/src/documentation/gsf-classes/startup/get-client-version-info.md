# GetClientVersionInfo

## Description

Sent on game start, before the main menu is interactive.
A temporary session is created for this check.
The client needs to know whether it is running an outdated version of the game.
On receiving the response, it compares the server version number against `GameSettings.clientLocalVersion`.

## Request

| Field        | Type     | Description                   |
| ------------ | -------- | ----------------------------- |
| `ClientName` | `string` | Hardcoded to `"AmazingWorld"` |

## Response

| Field               | Type     | Description                                                                                                                                           |
| ------------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| `ClientVersionInfo` | `string` | Version string in `"<version>.<forceUpdate>"` format, e.g. `"133852.true"`. The `forceUpdate` part is either `true` (blocking) or `false` (optional). |
