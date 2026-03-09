# GetOnlineStatuses

## Description

Sent during the `ZoneTransitionQueue`.

Fetches the online statuses of the player's friends. The result is stored in `FriendManager.Instance.statusList`.

Only sent once per session (`InitFriendMangerNotificationManagerCommand.hasDone` guards against repetition).

## Request

No fields.

## Response

| Field            | Type                                          | Description             |
| ---------------- | --------------------------------------------- | ----------------------- |
| `OnlineStatuses` | [`[]OnlineStatus`](../types/online-status.md) | List of online statuses |
