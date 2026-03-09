# GetRandomNames

## Description

Sent in two situations during the intro sequence:

- When the player requests a random Zing name;
- When building the family name picker;

Returns a list of randomly selected names of the requested type.

## Request

| Field          | Type     | Description                                                                                    |
| -------------- | -------- | ---------------------------------------------------------------------------------------------- |
| `Amount`       | `int32`  | Number of names to return                                                                      |
| `NamePartType` | `string` | `"second_name"` for Zing names, `"Family_1"`, `"Family_2"`, `"Family_3"` for family name parts |

## Response

| Field   | Type       | Description                     |
| ------- | ---------- | ------------------------------- |
| `Names` | `[]string` | List of randomly selected names |
