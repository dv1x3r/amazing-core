# SelectPlayerName

## Description

Sent when the player confirms their **family name** by pressing the Next button.

Reserves the chosen family name on the server so no other player can claim it during the same session(?).

## Request

| Field  | Type     | Description                |
| ------ | -------- | -------------------------- |
| `Name` | `string` | The family name to reserve |

## Response

No fields.

Success is implied by a non-error response.

On `AppCode 71` (duplicate) the client shows an error dialog.
