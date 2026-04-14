# GetPublicItemCategories

## Description

This is the public item-category lookup used by the registration flow.

Its purpose is to give the intro scene enough category data to classify and render preview items that are shown before the player has fully registered.

The returned categories are added into `InventoryManager` and then used when the public items are converted into temporary `GSFPlayerItem`s.

This message is for public preview content, while [`GetCMSItemCategories`](../login/get-cms-item-categories.md) is for the actual logged-in game state.

## Request

| Field              | Type                     | Description                                           |
| ------------------ | ------------------------ | ----------------------------------------------------- |
| `LanglocalePairID` | [`OID`](../types/oid.md) | "Class": 4, "Type": 19 "Server": 0, "Number": 9023265 |
| `TierID`           | [`OID`](../types/oid.md) | "Class": 0, "Type": 0, "Server": 0, "Number": 0       |
| `BirthDate`        | `time.Time`              | `#N/A`                                                |
| `RegistrationDate` | `time.Time`              | `#N/A`                                                |
| `PreviewDate`      | `time.Time`              | `#N/A`                                                |
| `IsPreviewEnabled` | `bool`                   | `false`                                               |

## Response

| Field            | Type                                          | Description             |
| ---------------- | --------------------------------------------- | ----------------------- |
| `ItemCategories` | [`[]ItemCategory`](../types/item-category.md) | List of item categories |
