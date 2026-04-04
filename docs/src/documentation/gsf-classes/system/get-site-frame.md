# GetSiteFrame

## Description

- Sent when the intro scene begins loading (on registration).
- Sent after tiers are loaded, if `SessionID` is valid (on login).

Fetches the site frame configuration for the current session.

The response also carries `AssetDeliveryURL`, which the client stores in `GameSettings.assetDeliveryURL` and uses for asset downloads.

## Request

| Field              | Type                     | Description                                           |
| ------------------ | ------------------------ | ----------------------------------------------------- |
| `TypeValue`        | `int32`                  | Hardcoded to `1`                                      |
| `LanglocalePairID` | [`OID`](../types/oid.md) | "Class": 4, "Type": 19 "Server": 0, "Number": 9023265 |
| `TierID`           | [`OID`](../types/oid.md) | "Class": 0, "Type": 0, "Server": 0, "Number": 0       |
| `BirthDate`        | `time.Time`              | `#N/A`                                                |
| `RegistrationDate` | `time.Time`              | `#N/A`                                                |
| `PreviewDate`      | `time.Time`              | `#N/A`                                                |
| `IsPreviewEnabled` | `bool`                   | `false`                                               |

## Response

| Field              | Type                                  | Description                       |
| ------------------ | ------------------------------------- | --------------------------------- |
| `SiteFrame`        | [`SiteFrame`](../types/site-frame.md) | `SiteFrame` object                |
| `AssetDeliveryURL` | `string`                              | Base URL used for asset downloads |
