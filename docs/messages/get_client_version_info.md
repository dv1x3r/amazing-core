# Get Client Version Info

## Request

| Property    | Type   | Description                     |
| ----------- | ------ | ------------------------------- |
| client_name | string | Should always be "AmazingWorld" |

### Request Bit Stream

1. client_name

## Response

| Property            | Type   | Description                        |
| ------------------- | ------ | ---------------------------------- |
| client_version_info | string | Client version supported by server |
|                     |        | Latest version is "133852.true"    |
|                     |        | "true" means Force Update          |

### Response Bit Stream

1. client_version_info
