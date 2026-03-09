# CheckUsername

## Description

Sent when the player submits the registration form.

Verifies that the chosen username is not already taken.

## Request

| Field      | Type     | Description           |
| ---------- | -------- | --------------------- |
| `Username` | `string` | The username to check |
| `Password` | `string` | The account password  |

## Response

No fields.

Success is implied by a non-error response.

On `AppCode 301` the server signals the username is already in use.
