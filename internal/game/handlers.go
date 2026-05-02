package game

import (
	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/messages"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/services"
)

type Handler struct {
	svc services.Set
}

func NewHandler(svc services.Set) *Handler {
	return &Handler{
		svc: svc,
	}
}

const playerID = 1
const maxOutfit = 1

var (
	dummyPlayerID      = types.OIDFromInt64(1)
	dummyAvatarID      = types.OIDFromInt64(1)
	dummyOutfitID      = types.OIDFromInt64(1)
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

// ── General ──────────────────────────────────────────────────────────────────

// GetClientVersionInfo validates the client name and version. Requested on game start.
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

// GetSiteFrame sends the main asset container with core assets. Requested on login.
func (h *Handler) GetSiteFrame(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetSiteFrameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	siteFrame, err := h.svc.SiteFrame.GetGSFSiteFrame(r.Context(), r.Platform(), req.TypeValue)
	if err != nil {
		return err
	}
	res := &messages.GetSiteFrameResponse{}
	res.AssetDeliveryURL = h.svc.Asset.DeliveryURL()
	res.SiteFrame = siteFrame
	return w.Write(res)
}

// Login handles client login, records client platform data, and returns
// the player's active avatar and asset delivery configuration.
func (h *Handler) Login(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.LoginRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	// TODO: this should be stored in the session
	// Due the disconnect we can lose this state
	// Could be restored via SessionOID during the Relogin
	r.SetPlatform(gsf.ParsePlatformFromMachineOS(req.ClientEnvInfo.MachineOS))

	player, err := h.svc.Player.GetGSFPlayer(r.Context(), r.Platform(), playerID)
	if err != nil {
		return err
	}

	res := &messages.LoginResponse{}
	res.AssetDeliveryURL = h.svc.Asset.DeliveryURL()
	res.MaxOutfit = maxOutfit
	res.Player = player

	return w.Write(res)
}

// Logout acknowledges the client logout request.
func (h *Handler) Logout(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.LogoutRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.LogoutResponse{}
	return w.Write(res)
}

// ── Registration ─────────────────────────────────────────────────────────────

// GetPublicItemCategories sends public item categories to classify temporary player items.
// Requested during the new player registration.
func (h *Handler) GetPublicItemCategories(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPublicItemCategoriesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPublicItemCategoriesResponse{}
	res.ItemCategories = dummyItemCategories()
	return w.Write(res)
}

// GetRandomNames sends random Zing names or family name parts.
// Requested during the new player registration or Zing rename.
func (h *Handler) GetRandomNames(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetRandomNamesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	names, err := h.svc.RandName.GetNStringsByType(r.Context(), req.NamePartType, int(req.Amount))
	if err != nil {
		return err
	}
	res := &messages.GetRandomNamesResponse{}
	res.Names = names
	return w.Write(res)
}

// ValidateName checks if Zing name is polite enough.
// Requested on Zing name submission.
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

// SelectPlayerName acknowledges the selected family name during registration.
func (h *Handler) SelectPlayerName(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SelectPlayerNameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.SelectPlayerNameResponse{}
	return w.Write(res)
}

// CheckUsername handles the registration username availability request.
func (h *Handler) CheckUsername(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.CheckUsernameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.CheckUsernameResponse{}
	return w.Write(res)
}

// RegisterPlayer handles the player account registration request.
func (h *Handler) RegisterPlayer(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.RegisterPlayerRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.RegisterPlayerResponse{}
	res.PlayerID = types.OID{Class: 1, Type: 2, Server: 3, Number: 4}
	return w.Write(res)
}

// ── Avatars & Outfits ────────────────────────────────────────────────────────

// GetAvatars fetches the list of player avatars. The returned list is stored in AvatarManager.Instance.GSFPlayerAvatars.
func (h *Handler) GetAvatars(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAvatarsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	avatars, err := h.svc.Player.GetGSFAvatars(r.Context(), r.Platform(), playerID)
	if err != nil {
		return err
	}
	res := &messages.GetAvatarsResponse{}
	res.Avatars = avatars
	return w.Write(res)
}

// GetOutfits fetches the saved outfits for a given player avatar.
// The results are stored as PresetOutfits on the AvatarAssets object.
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

// GetOutfitItems fetches the item instances associated with a player avatar outfit.
func (h *Handler) GetOutfitItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOutfitItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	hatAsset, err := h.svc.Asset.GetGSFAssetByCDNID(r.Context(), dummyHatAssetCDNID)
	if err != nil {
		return err
	}
	hatAsset = normalizePrefabAsset(hatAsset)
	hatIcon, err := h.svc.Asset.GetGSFAssetByCDNID(r.Context(), dummyHatIconCDNID)
	if err != nil {
		return err
	}
	hatIcon = normalizeImageAsset(hatIcon)
	res := &messages.GetOutfitItemsResponse{}
	res.OutfitItems = []types.PlayerItem{dummyHatPlayerItem(hatAsset, hatIcon)}
	return w.Write(res)
}

// UpdatePlayerActiveAvatar handles the active-avatar change request and returns avatar data.
func (h *Handler) UpdatePlayerActiveAvatar(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.UpdatePlayerActiveAvatarRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	avatar, err := h.svc.Player.SetGSFPlayerActiveAvatar(r.Context(), r.Platform(), req.PlayerAvatarID)
	if err != nil {
		return err
	}
	res := &messages.UpdatePlayerActiveAvatarResponse{}
	res.ActivePlayerAvatar = avatar
	return w.Write(res)
}

// GetAvatarItems fetches item instances owned by the active player avatar.
func (h *Handler) GetAvatarItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAvatarItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetAvatarItemsResponse{}
	res.AvatarItems = []types.PlayerItem{}
	return w.Write(res)
}

// ── Login ────────────────────────────────────────────────────────────────────

// GetTiers fetches all subscription tiers on initial login.
func (h *Handler) GetTiers(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetTiersRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetTiersResponse{}
	res.Tiers = []types.Tier{}
	return w.Write(res)
}

// GetCMSItemCategories fetches the main item-category lookup into the InventoryManager.itemCategories.
// Systems rely on that to map category OIDs to item types such as Clothing, Decoration, Yard, MazePiece, and so on.
func (h *Handler) GetCMSItemCategories(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetCMSItemCategoriesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetCMSItemCategoriesResponse{}
	res.ItemCategories = dummyItemCategories()
	return w.Write(res)
}

// GetPublicItemsByOIDs handles requests for public item definitions by object ID.
func (h *Handler) GetPublicItemsByOIDs(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPublicItemsByOIDsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPublicItemsByOIDsResponse{}
	res.Items = []types.Item{}
	return w.Write(res)
}

// GetZones fetches the list of available zones.
// The list is stored in ZoneManager.Instance.zones.
func (h *Handler) GetZones(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetZonesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetZonesResponse{}
	npcZone := types.Zone{}
	// LoadNPCsCommand calls SpawnPoints.Instance.ParseZone(NPCManager.HardCodedZoneId),
	// so the NPC zone must be present in this list.
	// LoadNPCsCommand.cs -> SpawnPoints.Instance.ParseZone(NPCManager.HardCodedZoneId)
	npcZone.OID = types.OID{Class: 4, Type: 16, Server: 0, Number: 2937912}
	npcZone.AssetMap = map[string][]types.Asset{}
	res.Zones = []types.Zone{npcZone}
	return w.Write(res)
}

// InitLocation sends the player home location data and the address of the SYNC server.
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

	dummyScene, err := h.svc.Dummy.GetValue(r.Context(), "map")
	if err != nil {
		return err
	}

	scene, err := h.svc.Asset.GetGSFAssetByCDNID(r.Context(), dummyScene)
	if err != nil {
		return err
	}

	// The Home.PlayerMaze.HomeTheme.AssetMap["Scene_Unity3D"] asset drives scene loading via
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

// SyncLogin authenticates the player on the SYNC server.
// The SYNC server handles real-time positional and social updates.
func (h *Handler) SyncLogin(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SyncLoginRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.SyncLoginResponse{}
	return w.Write(res)
}

// GetMazeItems fetches the items placed in the player maze.
func (h *Handler) GetMazeItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetMazeItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetMazeItemsResponse{}
	res.MazeItems = []types.PlayerItem{}
	return w.Write(res)
}

// GetChatChannelTypes fetches the available chat channel type definitions.
func (h *Handler) GetChatChannelTypes(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetChatChannelTypesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetChatChannelTypesResponse{}
	res.ChatChannelTypes = []types.ChatChannelType{}
	return w.Write(res)
}

// GetAnnouncements fetches login announcements to be displayed to the player on entry.
func (h *Handler) GetAnnouncements(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAnnouncementsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetAnnouncementsResponse{}
	res.Announcements = []types.Announcement{}
	return w.Write(res)
}

// EnterBuilding notifies the server that the player is entering a building.
func (h *Handler) EnterBuilding(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.EnterBuildingRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.EnterBuildingResponse{}
	return w.Write(res)
}

// GetOnlineStatuses fetches the online statuses of the player’s friends.
// The result is stored in FriendManager.Instance.statusList.
func (h *Handler) GetOnlineStatuses(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOnlineStatusesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetOnlineStatusesResponse{}
	res.OnlineStatuses = []types.OnlineStatus{}
	return w.Write(res)
}

// GetPlayerNPCs fetches the NPCs for the player zone.
// SpawnPoints.Instance.ParseZone(NPCManager.HardCodedZoneId) is called beforehand to set up spawn points.
// This requires the NPC zone to have been returned by GetZones.
func (h *Handler) GetPlayerNPCs(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPlayerNPCsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPlayerNPCsResponse{}
	res.NPCs = []types.NPC{}
	return w.Write(res)
}
