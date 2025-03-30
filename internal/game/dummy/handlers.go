package dummy

import (
	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/messages"
	"github.com/dv1x3r/amazing-core/internal/game/types"
)

/*
GetClientVersionInfo returns the server version for a specific client.

Scenarios:
  - Client startup.

Returns:
  - A string in the format "version.force_update" where:
  - version is the current game version.
  - force_update is a boolean represented as "true" or "false", indicating whether a forced update is required.
*/
func GetClientVersionInfo(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetClientVersionInfoRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	res := &messages.GetClientVersionInfoResponse{}
	if req.ClientName == "AmazingWorld" {
		res.ClientVersionInfo = "133852.true"
	}

	return w.Write(res)
}

/*
GetPublicItemCategories returns item categories array.

Scenarios:
  - New player registration.

Returns:
  - ItemCategories array.
*/
func GetPublicItemCategories(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPublicItemCategoriesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPublicItemCategoriesResponse{}
	return w.Write(res)
}

/*
GetPublicItemsByOIDs returns items array.

Scenarios:
  - New player registration.

Returns:
  - Items array.
*/
func GetPublicItemsByOIDs(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPublicItemsByOIDsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPublicItemsByOIDsResponse{}
	return w.Write(res)
}

/*
ValidateName returns the result of checking the username for bad words.

Scenarios:
  - New player registration.

Returns:
  - An empty string when validation passed.
  - Non empty string when validation failed.
*/
func ValidateName(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.ValidateNameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	res := &messages.ValidateNameResponse{}
	if req.Name == "fuck" {
		res.FilterName = "-"
	}

	return w.Write(res)
}

/*
SelectPlayerName processes the family name selection.

Scenarios:
  - New player registration.

Returns:
  - An empty message.
*/
func SelectPlayerName(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SelectPlayerNameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.SelectPlayerNameResponse{}
	return w.Write(res)
}

/*
CheckUsername checks username and password in the registration form.

Scenarios:
  - New player registration.

Returns:
  - An empty message.
*/
func CheckUsername(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.CheckUsernameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.CheckUsernameResponse{}
	return w.Write(res)
}

/*
RegisterPlayer processes the registration form submission.

Scenarios:
  - New player registration.

Returns:
  - New PlayerID.
*/
func RegisterPlayer(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.RegisterPlayerRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.RegisterPlayerResponse{}
	res.PlayerID = types.OID{Class: 1, Type: 2, Server: 3, Number: 4}
	return w.Write(res)
}

/*
Login processes the login form submission.

Scenarios:
  - Initial login.

Returns:
  - LoginResponse object.
*/
func Login(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.LoginRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	res := &messages.LoginResponse{}

	// base url for downloadable assets
	res.AssetDeliveryURL = config.Get().Settings.AssetDeliveryURL

	// to pass AvatarAssembler.HandleBaseLoaded
	// !(item.resName == "PF__Avatar.unity3d")
	res.Player.ActivePlayerAvatar.Avatar.AssetMap = map[string][]types.Asset{
		"Prefab_Unity3D": {
			{
				OID:           types.OID{},
				AssetTypeName: "asset_type",
				CDNID:         "Player_Avatar.unity3d",
				ResName:       "Player_Avatar.unity3d",
				GroupName:     "asset_group",
				FileSize:      59109,
			},
		},
	}

	return w.Write(res)
}

/*
GetTiers returns tiers array.

Scenarios:
  - User login: ClientManager.GetTiersResponseHandler()

Returns:
  - Tiers array.
*/
func GetTiers(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetTiersRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetTiersResponse{}
	res.Tiers = []types.Tier{}
	return w.Write(res)
}

/*
GetSiteFrame returns the SiteFrame.

Scenarios:

- User login:
  - ClientManager.GetSiteFrame
  - Filling SiteContentFactory.Instance lists from assetMap
  - Filling ClientManager lists from assetMap
  - Download assets included into Preload_PrefabUnity3D assetMap and ShadersList.unity3d
  - Preloaded objects are downloaded to the Cached folder ClientManager.cs -> LoadPreloadedAssets()
  - Download Player_Base.unity3d AvatarAssembler.LoadAvatar()
  - Download cdn defined in the PlayerAvatar.Avatar.assetMap['Prefab_Unity3D'] AvatarAssembler.HandleBaseLoaded()

Returns:
  - SiteFrame object.
  - AssetDeliveryURL string.
*/
func GetSiteFrame(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetSiteFrameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	res := &messages.GetSiteFrameResponse{}

	// base url for downloadable assets (site_frame)
	res.AssetDeliveryURL = config.Get().Settings.AssetDeliveryURL // + cdn.cdn_id

	res.SiteFrame.AssetMap = map[string][]types.Asset{
		"Config_Text":           {}, // DressAvatarManager.cs -> LoadSlotIds -> ClientManager.Instance.configList
		"Preload_PrefabUnity3D": {}, // OutdoorMazeLoader.cs -> LoadPreloadAssetsCommand() -> preloadList
		"Audio":                 {}, // OutdoorMazeLoader.cs -> LoadPreloadAssetsCommand() -> audioClipList
		// this is used to load hardcoded assets (instead of using Resources.Load())
		"Amazing_Core": {
			// LoadLoginScene.cs -> LoadAvatar -> DownloadManager.LoadAsset("Player_Base.unity3d")
			{
				OID:           types.OID{},
				AssetTypeName: "asset_type",
				CDNID:         "OTU2NTAzNTgyMzExOA",
				ResName:       "Player_Base.unity3d",
				GroupName:     "asset_group",
				FileSize:      565066,
			},
			// OutdoorMazeLoader.cs -> LoadSharedPrefabsCommand -> DownloadManager.LoadAsset("PlayerCamera.unity3d")
			{
				OID:           types.OID{},
				AssetTypeName: "asset_type",
				CDNID:         "OTQyNDc5ODIyMDMwMg",
				ResName:       "PlayerCamera.unity3d",
				GroupName:     "asset_group",
				FileSize:      1878,
			},
			// ClientManager.cs -> LoadPreloadAssetsCommand -> DownloadManager.LoadAsset("ShadersList.unity3d")
			{
				OID:           types.OID{},
				AssetTypeName: "asset_type",
				CDNID:         "OTYyNDQwNDA5OTA4Ng",
				ResName:       "ShadersList.unity3d",
				GroupName:     "asset_group",
				FileSize:      91142,
			},
		},
	}

	return w.Write(res)
}

/*
GetOutfitItems returns outfit items array.

Scenarios:
  - User login: LoadCurrentOutfitCommand.Begin()

Returns:
  - PlayerItems array.
*/
func GetOutfitItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOutfitItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetOutfitItemsResponse{}
	res.OutfitItems = []types.PlayerItem{}
	return w.Write(res)
}

/*
GetAvatars returns player avatars array.

Scenarios:
  - User login: LoadAvatarsCommand.Begin()

Returns:
  - PlayerAvatars array.
*/
func GetAvatars(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAvatarsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	res := &messages.GetAvatarsResponse{}

	// activePlayerAvatar should be in AvatarManager.Instance.GSFPlayerAvatars: LoadAvatarsCommand.cs -> Step2()
	// PlayerAvatar.Avatar.assetmap = {'Prefab_Unity3D': [Asset(Player_Avatar.unity3d)\]}
	// !(item.resName == "PF__Avatar.unity3d"

	avatar := types.Avatar{}
	avatar.AssetMap = map[string][]types.Asset{
		"Prefab_Unity3D": {
			{
				OID:           types.OID{},
				AssetTypeName: "asset_type",
				CDNID:         "Player_Avatar.unity3d",
				ResName:       "Player_Avatar.unity3d",
				GroupName:     "asset_group",
				FileSize:      59109,
			},
		},
	}

	res.Avatars = []types.PlayerAvatar{{Avatar: avatar}}

	return w.Write(res)
}

/*
GetOutfits returns player avatar outfits array.

Scenarios:
  - User login: LoadAvatarsCommand.LoadedAvatarsHandler()

Returns:
  - PlayerAvatarOutfits array.
*/
func GetOutfits(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOutfitsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetOutfitsResponse{}
	res.PlayerAvatarOutfits = []types.PlayerAvatarOutfit{}
	return w.Write(res)
}

/*
GetZones returns zones array.

Scenarios:
  - User login: LoadZonesCommand.Begin()

Returns:
  - Zones array.
*/
func GetZones(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetZonesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	res := &messages.GetZonesResponse{}

	npcZone := types.Zone{}

	// LoadNPCsCommand.cs -> SpawnPoints.Instance.ParseZone(NPCManager.HardCodedZoneId)
	npcZone.OID = types.OID{Class: 4, Type: 16, Server: 0, Number: 2937912}
	npcZone.AssetMap = map[string][]types.Asset{}

	res.Zones = []types.Zone{npcZone}

	return w.Write(res)
}

/*
InitLocation returns data required for initial location load, and address for SYNC server.

Scenarios:
  - User login: LoadHomeZoneCommand.Begin()

Returns:
  - SyncServer address.
  - PlayerHome object.
  - ZoneInstance object.
  - Village object.
*/
func InitLocation(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.InitLocationRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	res := &messages.InitLocationResponse{}

	res.SyncServerIP = config.Get().Settings.SyncServerIP
	res.SyncServerPort = int32(config.Get().Settings.SyncServerPort)

	// LoadMazeCommand.cs -> LoadMainScene() -> AssetDownloadManager.cs -> LoadMainScene()
	homeTheme := types.AssetContainer{
		AssetMap: map[string][]types.Asset{
			"Scene_Unity3D": {
				{
					OID:           types.OID{},
					AssetTypeName: "asset_type",
					CDNID:         "non_existing_cdn_id",
					ResName:       "Springtime003.unity3d",
					// ResName:       "HomeLotSmall.unity3d",
					GroupName: "Main_Scene",
					FileSize:  0,
				},
			},
		},
	}

	playerMaze := types.PlayerMaze{
		Name:       "coremaze",
		MazePieces: []types.PlayerMazePiece{},
		HomeTheme:  homeTheme,
	}

	res.Home = types.PlayerHome{
		PlayerMaze:  playerMaze,
		HomeTheme:   homeTheme,
		PlayerMazes: []types.PlayerMaze{playerMaze},
	}

	return w.Write(res)
}

/*
GetMazeItems returns maze items array.

Scenarios:
  - User login: LoadMazeItemsCommand.Begin()

Returns:
  - PlayerItems array.
*/
func GetMazeItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetMazeItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetMazeItemsResponse{}
	res.MazeItems = []types.PlayerItem{}
	return w.Write(res)
}

/*
GetChatChannelTypes returns channel types array.

Scenarios:
  - User login: LoadChatChannelTypesCommand.Begin()

Returns:
  - ChatChannelTypes array.
*/
func GetChatChannelTypes(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetChatChannelTypesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetChatChannelTypesResponse{}
	res.ChatChannelTypes = []types.ChatChannelType{}
	return w.Write(res)
}

/*
GetAnnouncements returns announcements array.

Scenarios:
  - User login: LoadGlobalAnnouncementsCommand.Begin()

Returns:
  - Announcements array.
*/
func GetAnnouncements(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAnnouncementsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetAnnouncementsResponse{}
	res.Announcements = []types.Announcement{}
	return w.Write(res)
}

/*
SyncLogin processes the sync server connection initialization.

Scenarios:
  - User login: SyncManager.SyncLogin()

Returns:
  - An empty message.
*/
func SyncLogin(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SyncLoginRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.SyncLoginResponse{}
	return w.Write(res)
}

/*
EnterBuilding returns building id.

Scenarios:
  - User login: LoadEnterBuildingMazeCommand.Begin()

Returns:
  - BuildingID object.
*/
func EnterBuilding(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.EnterBuildingRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.EnterBuildingResponse{}
	return w.Write(res)
}

/*
GetOnlineStatuses returns online status array.

Scenarios:
  - User login: InitFriendMangerNotificationManagerCommand.Begin()

Returns:
  - OnlineStatuses array.
*/
func GetOnlineStatuses(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOnlineStatusesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetOnlineStatusesResponse{}
	res.OnlineStatuses = []types.OnlineStatus{}
	return w.Write(res)
}

/*
GetPlayerNPCs returns player NPCs array.

Scenarios:
  - User login: LoadNPCsCommand.Begin()

Returns:
  - NPC array.
*/
func GetPlayerNPCs(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPlayerNPCsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPlayerNPCsResponse{}
	res.NPCs = []types.NPC{}
	return w.Write(res)
}
