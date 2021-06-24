# Message Header

## Properties

| Property       | Type   | Description                                                   |
| -------------- | ------ | ------------------------------------------------------------- |
| flags          | int    |                                                               |
| service_class  | int    | ServiceClass Enum                                             |
| message_type   | int    | UserMessageTypes / SyncMessageTypes / ClientMessageTypes Enum |
| request_id     | int    |                                                               |
| log_correlator | string |                                                               |
| result_code    | int    | ResultCode Enum                                               |
| app_code       | int    | AppCode Enum                                                  |
| is_service     | bool   | flags & 2 == 0                                                |
| is_response    | bool   | is_service & (flags & 1) != 0                                 |
| is_request     | bool   | is_service & (not is_response)                                |
| is_notify      | bool   | flags & 2 != 0                                                |
| is_discardable | bool   | flags & 16 != 0                                               |

## Request Bit Stream

1. flags
2. service_class
3. message_type
4. if is_service:
   1. request_id
5. if is_request:
   1. log_correlator

## Response Bit Stream

1. flags
2. service_class
3. message_type
4. if is_service:
   1. request_id
5. if is_request:
   1. log_correlator
6. if is_response:
   1. result_code
   2. app_code
   3. if app_code != 0:
      1. app_code_name
      2. if app_code == 17:
         1. app_codes
