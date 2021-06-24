# Logging in

\* Not all login requests are presented. I mentioned only ones, which are required to enter the game.

| Action                     | Where                                              | What happens                                                                         |
| -------------------------- | -------------------------------------------------- | ------------------------------------------------------------------------------------ |
| 1. Start login process     | ClientManager.OnCreateSession()                    | Sending GSFLoginSvc                                                                  |
|                            | ClientManager.LoginResponseHandler()               |                                                                                      |
|                            | ClientManager.StartLogin()                         |                                                                                      |
| 2. Get Tiers               | ClientManager.GetTiers()                           | Sending GSFGetTiersSvc                                                               |
|                            | ClientManager.GetTiersResponseHandler()            |                                                                                      |
| 3. Get SiteFrame           | ClientManager.GetSiteFrame()                       | Sending GSFGetSiteFrameSvc                                                           |
|                            | ClientManager.GetSiteFrameResponseHandler()        | Filling SiteContentFactory.Instance lists from assetMap                              |
|                            | ClientManager.AnimationLoaded()                    | Filling ClientManager lists from assetMap                                            |
|                            | ClientManager.LoadPreloadedAssets()                | Download assets included into Preload_PrefabUnity3D assetMap and ShadersList.unity3d |
|                            | ClientManager.PreloadComplete()                    |                                                                                      |
| 4. Load Avatar             | LoadLoginScene.LoadLoginAvatar()                   | Pick activePlayerAvatar asset                                                        |
|                            | LoadLoginScene.LoadAvatar()                        | Create *new AvatarAssembler*                                                         |
|                            | AvatarAssembler.LoadAvatar()                       | Download Player_Base.unity3d                                                         |
|                            | AvatarAssembler.HandleBaseLoaded()                 | Download asset defined in the PlayerAvatar.Avatar.assetMap['Prefab_Unity3D'\]        |
|                            | LoadLoginScene.UpdateAvatar()                      |                                                                                      |
| 5. Initialize transition   | ClientManager.LoadZone()                           |                                                                                      |
|                            | ZoneManager.TravelToHomeYard()                     | Initialize tasks queue *new OutdoorMazeLoader*                                       |
| 6. Start transition tasks  | AWLoadingScreen.Begin()                            | Start load queue script                                                              |
| 7. Get Outfit Items        | LoadCurrentOutfitCommand.Begin()                   | Sending GSFGetOutfitItemsSvc                                                         |
| 8. Get Avatars             | LoadAvatarsCommand.Begin()                         | Sending GSFGetAvatarsSvc                                                             |
| 9. Get Outfits             | LoadAvatarsCommand.LoadedAvatarsHandler()          | Sending GSFGetOutfitsSvc                                                             |
| 10. Get Zones              | LoadZonesCommand.Begin()                           | Sending GSFGetZonesSvc                                                               |
| 11. Init Location          | LoadHomeZoneCommand.Begin()                        | Sending GSFInitLocationSvc                                                           |
| 12. Get Maze Items         | LoadMazeItemsCommand.Begin()                       | Sending GSFGetMazeItemsSvc                                                           |
| 13. Get Chat Channel Types | LoadChatChannelTypesCommand.Begin()                | Sending GSFGetChatChannelTypesSvc                                                    |
| 14. Get Announcements      | LoadGlobalAnnouncementsCommand.Begin()             | Sending GSFGetAnnouncementsSvc                                                       |
| 15. Sync Login             | SyncManager.SyncLogin()                            | Sending GSFSyncLoginSvc                                                              |
| 16. Enter Building         | LoadEnterBuildingMazeCommand.Begin()               | Sending GSFEnterBuildingSvc                                                          |
| 17. Get Online Statuses    | InitFriendMangerNotificationManagerCommand.Begin() | Sending GSFGetOnlineStatusesSvc                                                      |
| 18. Get Player NPCs        | LoadNPCsCommand.Begin()                            | Sending GSFGetPlayerNpcsSvc                                                          |

## Dummy parameters

### Login

```
- asset_delivery_url = 'http://localhost:8080/' # base url for downloadable assets
- player.active_player_avatar.avatar.assetmap = {'Prefab_Unity3D': [Asset(Player_Avatar.unity3d)\]} # to pass AvatarAssembler.HandleBaseLoaded
```

### Avatars

```
- avatars = [
    # activePlayerAvatar should be in AvatarManager.Instance.GSFPlayerAvatars: LoadAvatarsCommand.cs -> Step2()
    PlayerAvatar.Avatar.assetmap = {'Prefab_Unity3D': [Asset(Player_Avatar.unity3d)\]}
]
```

### SiteFrame

```
- asset_delivery_url = 'http://localhost:8080/' # base url for downloadable assets (site_frame)
- asset_map = {
    'Config_Text': [],              # DressAvatarManager.cs -> LoadSlotIds -> ClientManager.Instance.configList
    'Preload_PrefabUnity3D': [],    # OutdoorMazeLoader.cs -> LoadPreloadAssetsCommand() -> preloadList
    'Audio': [],                    # OutdoorMazeLoader.cs -> LoadPreloadAssetsCommand() -> audioClipList
    'Amazing_Core': [
        Asset('Player_Base.unity3d'),   # LoadLoginScene.cs -> LoadAvatar -> DownloadManager.LoadAsset("Player_Base.unity3d")
        Asset('PlayerCamera.unity3d'),  # OutdoorMazeLoader.cs -> LoadSharedPrefabsCommand -> DownloadManager.LoadAsset("PlayerCamera.unity3d")
    ]
}
```

### Zones

```
- zones = [
    # LoadNPCsCommand.cs -> SpawnPoints.Instance.ParseZone(NPCManager.HardCodedZoneId)
    Zone(ObjectID(4, 16, 0, 2937912)) 
]
```

### InitLocation

```
- sync_server_ip = 'localhost'
- sync_server_port = 8182

# LoadMazeCommand.cs -> LoadMainScene() -> AssetDownloadManager.cs -> LoadMainScene()
- home = PlayerHome(PlayerMaze(HomeTheme({'Scene_Unity3D': 'HomeLotSmall.unity3d'})))   # home lot
- home = PlayerHome(PlayerMaze(HomeTheme({'Scene_Unity3D': 'Springtime003.unity3d'})))  # main map
```
