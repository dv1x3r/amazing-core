# GetCMSItemCategories

## Description

This is the main item-category lookup for the logged-in game.

The response is cached into `InventoryManager.itemCategories`, and many systems rely on that cache to map category OIDs to semantic item types such as `Clothing`, `Decoration`, `Yard`, `MazePiece`, and so on.

## Request

No fields.

## Response

| Field            | Type                                          | Description             |
| ---------------- | --------------------------------------------- | ----------------------- |
| `ItemCategories` | [`[]ItemCategory`](../types/item-category.md) | List of item categories |
