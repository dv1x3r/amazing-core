# GetPublicItemCategories

## Description

These categories are somehow related to starter clothing for player's avatar before registering.

I do not know how this affects the intro scene at the moment.

## Request

| Field              | Type        | Description                                           |
| ------------------ | ----------- | ----------------------------------------------------- |
| `LanglocalePairID` | `OID`       | "Class": 4, "Type": 19 "Server": 0, "Number": 9023265 |
| `TierID`           | `OID`       | "Class": 0, "Type": 0, "Server": 0, "Number": 0       |
| `BirthDate`        | `time.Time` | `#N/A`                                                |
| `RegistrationDate` | `time.Time` | `#N/A`                                                |
| `PreviewDate`      | `time.Time` | `#N/A`                                                |
| `IsPreviewEnabled` | `bool`      | `false`                                               |

## Response

| Field            | Type                                          | Description             |
| ---------------- | --------------------------------------------- | ----------------------- |
| `ItemCategories` | [`[]ItemCategory`](../types/item-category.md) | List of item categories |
