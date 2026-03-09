# GetTiers

## Description

Sent immediately after a successful login response.

Fetches all subscription tiers. The client matches the player's `tierID` from `PlayerInfoTO` against the returned list to determine whether the account is paid or free-play.

On success, `GetSiteFrame()` is called if `SessionID` is valid.

## Request

No fields.

## Response

| Field   | Type                         | Description   |
| ------- | ---------------------------- | ------------- |
| `Tiers` | [`[]Tier`](../types/tier.md) | List of tiers |
