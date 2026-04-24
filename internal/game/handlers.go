package game

import (
	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/messages"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/services/asset"
	"github.com/dv1x3r/amazing-core/internal/services/dummy"
	"github.com/dv1x3r/amazing-core/internal/services/randname"
	"github.com/dv1x3r/amazing-core/internal/services/siteframe"
)

type Handler struct {
	assetService     *asset.Service
	dummyService     *dummy.Service
	randnameService  *randname.Service
	siteFrameService *siteframe.Service
}

func NewHandler(
	assetService *asset.Service,
	dummyService *dummy.Service,
	randnameService *randname.Service,
	siteFrameService *siteframe.Service,
) *Handler {
	return &Handler{
		assetService:     assetService,
		dummyService:     dummyService,
		randnameService:  randnameService,
		siteFrameService: siteFrameService,
	}
}

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

func (h *Handler) GetClientVersionInfo(w gsf.ResponseWriter, r *gsf.Request) error {
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

func (h *Handler) GetPublicItemCategories(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPublicItemCategoriesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPublicItemCategoriesResponse{}
	res.ItemCategories = dummyItemCategories()
	return w.Write(res)
}

func (h *Handler) GetCMSItemCategories(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetCMSItemCategoriesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetCMSItemCategoriesResponse{}
	res.ItemCategories = dummyItemCategories()
	return w.Write(res)
}

func (h *Handler) GetPublicItemsByOIDs(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPublicItemsByOIDsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPublicItemsByOIDsResponse{}
	res.Items = []types.Item{}
	return w.Write(res)
}

func (h *Handler) GetRandomNames(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetRandomNamesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	names, err := h.randnameService.GetNStringsByType(r.Context(), req.NamePartType, int(req.Amount))
	if err != nil {
		return err
	}
	res := &messages.GetRandomNamesResponse{}
	res.Names = names
	return w.Write(res)
}

func (h *Handler) ValidateName(w gsf.ResponseWriter, r *gsf.Request) error {
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

func (h *Handler) SelectPlayerName(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SelectPlayerNameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.SelectPlayerNameResponse{}
	return w.Write(res)
}

func (h *Handler) CheckUsername(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.CheckUsernameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.CheckUsernameResponse{}
	return w.Write(res)
}

func (h *Handler) RegisterPlayer(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.RegisterPlayerRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.RegisterPlayerResponse{}
	res.PlayerID = types.OID{Class: 1, Type: 2, Server: 3, Number: 4}
	return w.Write(res)
}

func (h *Handler) Login(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.LoginRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	// TODO: this should be stored in the session
	// Due the disconnect we can lose this state
	// Could be restored via SessionOID during the Relogin
	r.SetPlatform(gsf.ParsePlatformFromMachineOS(req.ClientEnvInfo.MachineOS))

	res := &messages.LoginResponse{}
	res.AssetDeliveryURL = h.assetService.DeliveryURL()

	avatarAsset, err := h.assetService.GetGSFAssetByCDNID(r.Context(), dummyAvatarCDNID)
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

func (h *Handler) GetTiers(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetTiersRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetTiersResponse{}
	res.Tiers = []types.Tier{}
	return w.Write(res)
}

func (h *Handler) GetSiteFrame(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetSiteFrameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	siteFrame, err := h.siteFrameService.GetGSFSiteFrame(r.Context(), r.Platform(), req.TypeValue)
	if err != nil {
		return err
	}
	res := &messages.GetSiteFrameResponse{}
	res.AssetDeliveryURL = h.assetService.DeliveryURL()
	res.SiteFrame = siteFrame
	return w.Write(res)
}

func (h *Handler) GetOutfitItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOutfitItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	hatAsset, err := h.assetService.GetGSFAssetByCDNID(r.Context(), dummyHatAssetCDNID)
	if err != nil {
		return err
	}
	hatAsset = normalizePrefabAsset(hatAsset)
	hatIcon, err := h.assetService.GetGSFAssetByCDNID(r.Context(), dummyHatIconCDNID)
	if err != nil {
		return err
	}
	hatIcon = normalizeImageAsset(hatIcon)
	res := &messages.GetOutfitItemsResponse{}
	res.OutfitItems = []types.PlayerItem{dummyHatPlayerItem(hatAsset, hatIcon)}
	return w.Write(res)
}

func (h *Handler) GetAvatars(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAvatarsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	avatarAsset, err := h.assetService.GetGSFAssetByCDNID(r.Context(), dummyAvatarCDNID)
	if err != nil {
		return err
	}
	avatarAsset = normalizePrefabAsset(avatarAsset)
	res := &messages.GetAvatarsResponse{}
	res.Avatars = []types.PlayerAvatar{dummyAvatarWithAsset(avatarAsset)}
	return w.Write(res)
}

func (h *Handler) GetOutfits(w gsf.ResponseWriter, r *gsf.Request) error {
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

func (h *Handler) GetZones(w gsf.ResponseWriter, r *gsf.Request) error {
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

func (h *Handler) InitLocation(w gsf.ResponseWriter, r *gsf.Request) error {
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

	dummyScene, err := h.dummyService.GetValue(r.Context(), "map")
	if err != nil {
		return err
	}

	scene, err := h.assetService.GetGSFAssetByCDNID(r.Context(), dummyScene)
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

func (h *Handler) GetMazeItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetMazeItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetMazeItemsResponse{}
	res.MazeItems = []types.PlayerItem{}
	return w.Write(res)
}

func (h *Handler) GetChatChannelTypes(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetChatChannelTypesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetChatChannelTypesResponse{}
	res.ChatChannelTypes = []types.ChatChannelType{}
	return w.Write(res)
}

func (h *Handler) GetAnnouncements(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAnnouncementsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetAnnouncementsResponse{}
	res.Announcements = []types.Announcement{}
	return w.Write(res)
}

func (h *Handler) SyncLogin(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SyncLoginRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.SyncLoginResponse{}
	return w.Write(res)
}

func (h *Handler) EnterBuilding(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.EnterBuildingRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.EnterBuildingResponse{}
	return w.Write(res)
}

func (h *Handler) GetOnlineStatuses(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOnlineStatusesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetOnlineStatusesResponse{}
	res.OnlineStatuses = []types.OnlineStatus{}
	return w.Write(res)
}

func (h *Handler) GetPlayerNPCs(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPlayerNPCsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPlayerNPCsResponse{}
	res.NPCs = []types.NPC{}
	return w.Write(res)
}

func (h *Handler) Logout(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.LogoutRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.LogoutResponse{}
	return w.Write(res)
}
