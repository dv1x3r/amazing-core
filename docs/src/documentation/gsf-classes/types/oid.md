# OID

## Description

A unique object GSFOID identifier usually stored as a big `int64` number, that can be sliced into specific groups of bits to get the Class, Type, Server and Object number.

```
[Class][Type][Server][Number...]
  8b     8b     8b       40b
```

## Fields

| Field    | Type    | Description                                             |
| -------- | ------- | ------------------------------------------------------- |
| `Class`  | `byte`  | For cache files is always 0                             |
| `Type`   | `byte`  | For cache files is always 0                             |
| `Server` | `byte`  | Server version? Goes from 1 to 8 (higher = newer asset) |
| `Number` | `int64` | Object number                                           |
