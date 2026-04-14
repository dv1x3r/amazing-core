# Login Messages

When the player submits their credentials on the login screen, the client opens a session and begins the login sequence.

After the initial login response the client runs two load queues:

- **ClientServices**: starts in `ClientManager.StartLogin()`, initialises background services (including connection to the SYNC server).

- **ZoneTransitionQueue**: starts in `ClientManager.LoadZone()` after the `SiteFrame` assets are preloaded, runs the full zone-entry sequence.

After the `SiteFrame` load, `AnimationLoaded()` is then called inline, which starts preloading assets.

| #  | Message                                                | Invoked at                                        |
| -- | ------------------------------------------------------ | ------------------------------------------------- |
| 1  | [`Login`](./login.md)                                  | `ClientManager.cs:567`                            |
| 2  | [`GetTiers`](./get-tiers.md)                           | `ClientManager.cs:652`                            |
| 3  | [`GetSiteFrame`](../system/get-site-frame.md)          | `ClientManager.cs:647`                            |
| 4  | [`GetCMSItemCategories`](./get-cms-item-categories.md) | `InventoryManager.cs:846`                         |
| 5  | [`GetOutfitItems`](./get-outfit-items.md)              | `LoadCurrentOutfitCommand.cs:5`                   |
| 6  | [`GetAvatars`](./get-avatars.md)                       | `LoadAvatarsCommand.cs:13`                        |
| 7  | [`GetOutfits`](./get-outfits.md)                       | `AvatarAssetsLoader.cs:138`                       |
| 8  | [`GetZones`](./get-zones.md)                           | `LoadZonesCommand.cs:18`                          |
| 9  | [`InitLocation`](./init-location.md)                   | `LoadHomeZoneCommand.cs:23`                       |
| 10 | [`SyncLogin`](./sync-login.md)                         | `SyncManager.cs:331`                              |
| 11 | [`GetMazeItems`](./get-maze-items.md)                  | `LoadMazeItemsCommand.cs:28`                      |
| 12 | [`GetChatChannelTypes`](./get-chat-channel-types.md)   | `LoadChatChannelTypesCommand.cs:7`                |
| 13 | [`GetAnnouncements`](./get-announcements.md)           | `LoadGlobalAnnouncementsCommand.cs:11`            |
| 14 | [`EnterBuilding`](./enter-building.md)                 | `LoadEnterBuildingMazeCommand.cs:49`              |
| 15 | [`GetOnlineStatuses`](./get-online-statuses.md)        | `InitFriendMangerNotificationManagerCommand.cs:9` |
| 16 | [`GetPlayerNPCs`](./get-player-npcs.md)                | `LoadNPCsCommand.cs:29`                           |
