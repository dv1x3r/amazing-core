# Login Messages

When the player submits their credentials on the login screen, the client opens a session and begins the login sequence.

After the initial login response the client runs two load queues:

- **ClientServices**: starts in `ClientManager.StartLogin()`, initialises background services (including connection to the SYNC server).

- **ZoneTransitionQueue**: starts in `ClientManager.LoadZone()` after the `SiteFrame` assets are preloaded, runs the full zone-entry sequence.

After the `SiteFrame` load, `AnimationLoaded()` is then called inline, which starts preloading assets.

| #  | Message                                              | Invoked at                                        |
| -- | ---------------------------------------------------- | ------------------------------------------------- |
| 1  | [`Login`](./login.md)                                | `ClientManager.cs:567`                            |
| 2  | [`GetTiers`](./get-tiers.md)                         | `ClientManager.cs:652`                            |
| 3  | [`GetSiteFrame`](../system/get-site-frame.md)        | `ClientManager.cs:647`                            |
| 4  | [`GetOutfitItems`](./get-outfit-items.md)            | `LoadCurrentOutfitCommand.cs:5`                   |
| 5  | [`GetAvatars`](./get-avatars.md)                     | `LoadAvatarsCommand.cs:13`                        |
| 6  | [`GetOutfits`](./get-outfits.md)                     | `AvatarAssetsLoader.cs:138`                       |
| 7  | [`GetZones`](./get-zones.md)                         | `LoadZonesCommand.cs:18`                          |
| 8  | [`InitLocation`](./init-location.md)                 | `LoadHomeZoneCommand.cs:23`                       |
| 9  | [`SyncLogin`](./sync-login.md)                       | `SyncManager.cs:331`                              |
| 10 | [`GetMazeItems`](./get-maze-items.md)                | `LoadMazeItemsCommand.cs:28`                      |
| 11 | [`GetChatChannelTypes`](./get-chat-channel-types.md) | `LoadChatChannelTypesCommand.cs:7`                |
| 12 | [`GetAnnouncements`](./get-announcements.md)         | `LoadGlobalAnnouncementsCommand.cs:11`            |
| 13 | [`EnterBuilding`](./enter-building.md)               | `LoadEnterBuildingMazeCommand.cs:49`              |
| 14 | [`GetOnlineStatuses`](./get-online-statuses.md)      | `InitFriendMangerNotificationManagerCommand.cs:9` |
| 15 | [`GetPlayerNPCs`](./get-player-npcs.md)              | `LoadNPCsCommand.cs:29`                           |
