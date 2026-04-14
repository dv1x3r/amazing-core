package dummy

import (
	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/messages"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/services/asset"
)

var (
	DummyService *Service
	AssetService *asset.Service
)

var (
	dummyPlayerID      = types.OIDFromInt64(72057594037927937)
	dummyAvatarID      = types.OIDFromInt64(72057594037927938)
	dummyOutfitID      = types.OIDFromInt64(72057594037927939)
	dummyHatItemID     = types.OIDFromInt64(72057594037927940)
	dummyHatTemplateID = types.OIDFromInt64(72057594037927941)
	dummyClothingCatID = types.OIDFromInt64(72057594037927942)
	dummyHatSlotID     = types.OIDFromInt64(289356276061314068)
	dummyAvatarCDNID   = "OTQ3ODg2NDg5NjAxNA"
	dummyHatAssetCDNID = "OTYyOTU0MTA3MjkxMA"
	dummyHatIconCDNID  = "OTYyOTM4NjkzMjIzOA"
	avatarSlotIDsCDNID = "OTU3MDUxOTg3NTU5OA"
	playerBaseCDNID    = "OTQ4NzQwNDQ5ODk1OA"
	playerCameraCDNID  = "OTQ4NzQyMTI3NjE3NA"
	shadersListCDNID   = "OTYyMzU1MDU1ODIyMg"
)

func dummyItemCategories() []types.ItemCategory {
	return []types.ItemCategory{
		{
			RuleContainer: types.RuleContainer{
				AssetContainer: types.AssetContainer{
					OID:      dummyClothingCatID,
					AssetMap: map[string][]types.Asset{},
				},
			},
			Name: "Clothing",
		},
	}
}

func normalizePrefabAsset(asset types.Asset) types.Asset {
	asset.AssetTypeName = "Prefab_Unity3D"
	asset.GroupName = ""
	return asset
}

func normalizeImageAsset(asset types.Asset) types.Asset {
	asset.AssetTypeName = "Images"
	if asset.GroupName == "" {
		asset.GroupName = "Inventory Icon"
	}
	return asset
}

func dummyAvatarWithAsset(avatarAsset types.Asset) types.PlayerAvatar {
	return types.PlayerAvatar{
		OID:                  dummyAvatarID,
		PlayerID:             dummyPlayerID,
		Name:                 "dummy-zing",
		PlayerAvatarOutfitID: dummyOutfitID,
		OutfitNo:             1,
		Avatar: types.Avatar{
			AssetContainer: types.AssetContainer{
				AssetMap: map[string][]types.Asset{
					"Prefab_Unity3D": {avatarAsset},
				},
			},
			MaxOutfits: 1,
			Name:       "dummy-zing",
		},
	}
}

func dummyHatPlayerItem(hatAsset, hatIcon types.Asset) types.PlayerItem {
	return types.PlayerItem{
		OID:                  dummyHatItemID,
		SlotID:               dummyHatSlotID,
		PlayerAvatarOutfitID: dummyOutfitID,
		PlayerAvatarID:       dummyAvatarID,
		PlayerID:             dummyPlayerID,
		Item: types.Item{
			AssetContainer: types.AssetContainer{
				OID: dummyHatTemplateID,
				AssetMap: map[string][]types.Asset{
					"Prefab_Unity3D": {hatAsset},
					"Images":         {hatIcon},
				},
			},
			Name:              "dummy-hat",
			ItemCategories:    dummyItemCategories(),
			AcceptableSlotIds: []types.OID{dummyHatSlotID},
		},
		Quantity: 1,
	}
}

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
	res.ItemCategories = dummyItemCategories()
	return w.Write(res)
}

func GetCMSItemCategories(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetCMSItemCategoriesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetCMSItemCategoriesResponse{}
	res.ItemCategories = dummyItemCategories()
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
	res.Items = []types.Item{}
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
	ctx := r.Context()
	req := &messages.LoginRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	res := &messages.LoginResponse{}

	// base url for downloadable assets
	res.AssetDeliveryURL = config.Get().Settings.AssetDeliveryURL

	avatarAsset, err := AssetService.GetGSFAssetByCDNID(ctx, dummyAvatarCDNID)
	if err != nil {
		return err
	}
	avatarAsset = normalizePrefabAsset(avatarAsset)

	res.Player.OID = dummyPlayerID
	res.Player.ActivePlayerAvatar = dummyAvatarWithAsset(avatarAsset)
	res.Player.IsQA = true
	res.MaxOutfit = 1

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
	ctx := r.Context()
	req := &messages.GetSiteFrameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	res := &messages.GetSiteFrameResponse{}

	// base url for downloadable assets (site_frame) + cdn.cdn_id
	res.AssetDeliveryURL = config.Get().Settings.AssetDeliveryURL

	// Player_Base.unity3d
	playerBase, err := AssetService.GetGSFAssetByCDNID(ctx, playerBaseCDNID)
	if err != nil {
		return err
	}

	// PlayerCamera.unity3d
	playerCamera, err := AssetService.GetGSFAssetByCDNID(ctx, playerCameraCDNID)
	if err != nil {
		return err
	}

	// ShadersList.unity3d
	shadersList, err := AssetService.GetGSFAssetByCDNID(ctx, shadersListCDNID)
	if err != nil {
		return err
	}

	// Avatar_SlotIds.txt
	// slot registry for wearable/equippable avatar parts
	// the client looks for an asset whose resName is Avatar_SlotIds.txt (DressAvatarManager.cs:917)
	avatarSlotIDs, err := AssetService.GetGSFAssetByCDNID(ctx, avatarSlotIDsCDNID)
	if err != nil {
		return err
	}

	res.SiteFrame.AssetMap = map[string][]types.Asset{}

	// DressAvatarManager.cs -> LoadSlotIds -> ClientManager.Instance.configList
	res.SiteFrame.AssetMap["Config_Text"] = []types.Asset{
		avatarSlotIDs,
	}

	// OutdoorMazeLoader.cs -> LoadPreloadAssetsCommand() -> preloadList
	res.SiteFrame.AssetMap["Preload_PrefabUnity3D"] = []types.Asset{}

	// OutdoorMazeLoader.cs -> LoadPreloadAssetsCommand() -> audioClipList
	res.SiteFrame.AssetMap["Audio"] = []types.Asset{}

	// this is used to load hardcoded assets (instead of using Resources.Load()) "Amazing_Core": {
	res.SiteFrame.AssetMap["Prefab_Unity3D"] = []types.Asset{
		// LoadLoginScene.cs -> LoadAvatar -> DownloadManager.LoadAsset("Player_Base.unity3d")
		playerBase,
		// OutdoorMazeLoader.cs -> LoadSharedPrefabsCommand -> DownloadManager.LoadAsset("PlayerCamera.unity3d")
		playerCamera,
		// ClientManager.cs -> LoadPreloadAssetsCommand -> DownloadManager.LoadAsset("ShadersList.unity3d")
		shadersList,
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
	ctx := r.Context()
	req := &messages.GetOutfitItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	hatAsset, err := AssetService.GetGSFAssetByCDNID(ctx, dummyHatAssetCDNID)
	if err != nil {
		return err
	}
	hatAsset = normalizePrefabAsset(hatAsset)
	hatIcon, err := AssetService.GetGSFAssetByCDNID(ctx, dummyHatIconCDNID)
	if err != nil {
		return err
	}
	hatIcon = normalizeImageAsset(hatIcon)

	res := &messages.GetOutfitItemsResponse{}
	res.OutfitItems = []types.PlayerItem{dummyHatPlayerItem(hatAsset, hatIcon)}
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
	ctx := r.Context()
	req := &messages.GetAvatarsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	avatarAsset, err := AssetService.GetGSFAssetByCDNID(ctx, dummyAvatarCDNID)
	if err != nil {
		return err
	}
	avatarAsset = normalizePrefabAsset(avatarAsset)

	res := &messages.GetAvatarsResponse{}
	res.Avatars = []types.PlayerAvatar{dummyAvatarWithAsset(avatarAsset)}
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
	res.PlayerAvatarOutfits = []types.PlayerAvatarOutfit{
		{
			OID:            dummyOutfitID,
			PlayerID:       dummyPlayerID,
			PlayerAvatarID: dummyAvatarID,
			OutfitNo:       1,
		},
	}
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
	ctx := r.Context()
	req := &messages.InitLocationRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	res := &messages.InitLocationResponse{}

	res.SyncServerIP = config.Get().Settings.SyncServerIP
	res.SyncServerPort = int32(config.Get().Settings.SyncServerPort)

	// dummyScene := "OTYwOTUyODk5OTk1MA" // Springtime003.unity3d
	// dummyScene := "OTYxMTQ4NDU5NDE5MA" // HomeLotSmall.unity3d
	// dummyScene := "OTQ1MDc3NTY0MjEyNg" // HomeLot_Winter.unity3d

	dummyScene, err := DummyService.GetValue(ctx, "map")
	if err != nil {
		return err
	}

	scene, err := AssetService.GetGSFAssetByCDNID(ctx, dummyScene)
	if err != nil {
		return err
	}

	// LoadMazeCommand.cs -> LoadMainScene() -> AssetDownloadManager.cs -> LoadMainScene()
	homeTheme := types.AssetContainer{}
	homeTheme.AssetMap = map[string][]types.Asset{}
	homeTheme.AssetMap["Scene_Unity3D"] = []types.Asset{scene}

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

func Logout(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.LogoutRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.LogoutResponse{}
	return w.Write(res)
}
