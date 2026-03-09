# ValidateName

## Description

Sent when the player confirms their Zing name during the sign up.

Checks whether the chosen player name contains any filtered or prohibited words.

## Request

| Field  | Type     | Description                 |
| ------ | -------- | --------------------------- |
| `Name` | `string` | The player name to validate |

## Response

| Field        | Type     | Description                                                                                     |
| ------------ | -------- | ----------------------------------------------------------------------------------------------- |
| `FilterName` | `string` | The client does not really care about the value. A non-empty value means the name was rejected. |
