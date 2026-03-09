# Asset

## Description

Represents a single downloadable game file.

## Fields

| Field           | Type     | Description                                                                     |
| --------------- | -------- | ------------------------------------------------------------------------------- |
| `OID`           | `OID`    | `OID` identifier                                                                |
| `AssetTypeName` | `string` | Asset classification (e.g. `"Prefab_Unity3D"`)                                  |
| `CDNID`         | `string` | `CDN` identifier used to construct the download URL: `AssetDeliveryURL + CDNID` |
| `ResName`       | `string` | Resource name                                                                   |
| `GroupName`     | `string` | Group classification (e.g. `"Main_Scene"`, `"3D Components"`, `"Locked"`)       |
| `FileSize`      | `int64`  | File size in bytes                                                              |
