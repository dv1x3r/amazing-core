# GetAnnouncements

## Description

Sent during the `ZoneTransitionQueue`.

Fetches login announcements to be displayed to the player on entry.

## Request

| Field      | Type   | Description |
| ---------- | ------ | ----------- |
| `UnMarked` | `bool` |             |

## Response

| Field           | Type                                         | Description           |
| --------------- | -------------------------------------------- | --------------------- |
| `Announcements` | [`[]Announcement`](../types/announcement.md) | List of announcements |
