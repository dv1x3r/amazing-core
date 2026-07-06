package game

import (
	"errors"
	"math/rand"
	"time"

	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/game/worldsync"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/messages"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/notify"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/appcode"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/resultcode"
	"github.com/dv1x3r/amazing-core/internal/services"
	"github.com/dv1x3r/amazing-core/internal/services/auth"
)

type Handler struct {
	svc     services.Services
	syncHub *worldsync.Hub
}

func NewHandler(svc services.Services) *Handler {
	return &Handler{
		svc:     svc,
		syncHub: worldsync.NewHub(),
	}
}

const (
	SCENE_HOMELOT_SMALL  = "OTYxMTQ4NDU5NDE5MA"
	SCENE_HOMELOT_WINTER = "OTQ1MDc3NTY0MjEyNg"
	SCENE_SPRINGTIME     = "OTYwOTUyODk5OTk1MA"
)

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

	r.SetPlatform(gsf.ParsePlatformFromMachineOS(req.ClientEnvInfo.MachineOS))

	playerOID, err := h.svc.Auth.LoginOrCreatePlayer(r.Context(), req.LoginID, req.Password)
	if err != nil {
		return wrap.WithGSFError(err, int32(resultcode.ERR), int32(loginAppCode(err)))
	}
	r.SetPlayerOID(playerOID.Int64())

	player, err := h.svc.Player.GetGSFPlayer(r.Context(), r.Platform(), playerOID)
	if err != nil {
		return wrap.WithGSFError(err, int32(resultcode.ERR), int32(appcode.INTERNAL_ERROR))
	}

	res := &messages.LoginResponse{}
	res.AssetDeliveryURL = h.svc.Asset.DeliveryURL()
	res.Player = player

	return w.Write(res)
}

func loginAppCode(err error) appcode.AppCode {
	switch {
	case errors.Is(err, auth.ErrBlankCredentials), errors.Is(err, auth.ErrInvalidCredentials):
		return appcode.AUTH
	case errors.Is(err, auth.ErrCreatePlayerAccountFailed):
		return appcode.CREATE_PLAYER_ACCOUNT_FAILED
	default:
		return appcode.INTERNAL_ERROR
	}
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

// Heartbeat keeps the USER server session alive and refreshes stats/time payloads.
func (h *Handler) Heartbeat(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.HeartbeatRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.HeartbeatResponse{
		PlayerStats:       []types.PlayerStats{},
		CurrentServerTime: gsf.UnixTime{Time: time.Now().UTC()},
		QueueNotify:       []types.PlayerNotify{},
	}
	return w.Write(res)
}

// ── Registration ─────────────────────────────────────────────────────────────

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
	res.PlayerOID = types.OID{Class: 1, Type: 2, Server: 3, Number: 4}
	return w.Write(res)
}

// ── Inventory ────────────────────────────────────────────────────────────────

// GetPublicItemCategories sends public item categories to classify temporary player items.
// Requested during the new player registration.
func (h *Handler) GetPublicItemCategories(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPublicItemCategoriesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	categories, err := h.svc.Item.GetGSFItemCategories(r.Context(), true)
	if err != nil {
		return err
	}
	res := &messages.GetPublicItemCategoriesResponse{}
	res.ItemCategories = categories
	return w.Write(res)
}

// GetCMSItemCategories fetches the main item-category lookup into the InventoryManager.itemCategories.
// Systems rely on that to map category OIDs to item types such as Clothing, Decoration, Yard, MazePiece, and so on.
func (h *Handler) GetCMSItemCategories(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetCMSItemCategoriesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	categories, err := h.svc.Item.GetGSFItemCategories(r.Context(), false)
	if err != nil {
		return err
	}
	res := &messages.GetCMSItemCategoriesResponse{}
	res.ItemCategories = categories
	return w.Write(res)
}

// GetBuildObjects fetches player-owned build objects.
func (h *Handler) GetBuildObjects(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetBuildObjectsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetBuildObjectsResponse{}
	res.PlayerBuildObjects = []types.PlayerBuildObject{}
	return w.Write(res)
}

// GetInventoryObjects fetches player-owned objects for the inventory grid.
func (h *Handler) GetInventoryObjects(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetInventoryObjectsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	playerOID, ok := r.PlayerOID()
	if !ok {
		return errors.New("player oid is not set on session")
	}
	items, err := h.svc.Player.GetGSFInventoryObjects(r.Context(), r.Platform(), types.OIDFromInt64(playerOID))
	if err != nil {
		return err
	}
	res := &messages.GetInventoryObjectsResponse{}
	res.PlayerItems = items
	return w.Write(res)
}

// ── Avatars ──────────────────────────────────────────────────────────────────

// GetAvatars fetches the list of player avatars. The returned list is stored in AvatarManager.Instance.GSFPlayerAvatars.
func (h *Handler) GetAvatars(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAvatarsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	playerOID, ok := r.PlayerOID()
	if !ok {
		return errors.New("player oid is not set on session")
	}
	avatars, err := h.svc.Player.GetGSFPlayerAvatars(r.Context(), r.Platform(), types.OIDFromInt64(playerOID))
	if err != nil {
		return err
	}
	res := &messages.GetAvatarsResponse{}
	res.Avatars = avatars
	return w.Write(res)
}

// GetAvatarItems fetches item instances owned by the active player avatar.
// This handler was used before outfit presets were implemented, probably.
// It is retained for backwards compatibility and may be repurposed for active avatar effects.
// Clothing should be managed with GetOutfitItems instead.
func (h *Handler) GetAvatarItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAvatarItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetAvatarItemsResponse{}
	res.AvatarItems = []types.PlayerItem{}
	return w.Write(res)
}

// UpdatePlayerActiveAvatar handles the active-avatar change request and returns avatar data.
func (h *Handler) UpdatePlayerActiveAvatar(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.UpdatePlayerActiveAvatarRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	avatar, err := h.svc.Player.SetGSFPlayerActiveAvatar(r.Context(), r.Platform(), req.PlayerAvatarOID)
	if err != nil {
		return err
	}
	if err := h.syncHub.AppearanceChanged(avatar.PlayerOID); err != nil {
		return err
	}
	res := &messages.UpdatePlayerActiveAvatarResponse{}
	res.ActivePlayerAvatar = avatar
	return w.Write(res)
}

// ── Outfits ──────────────────────────────────────────────────────────────────

// GetOutfits fetches the saved outfits for a given player avatar.
// The results are stored as PresetOutfits on the AvatarAssets object.
func (h *Handler) GetOutfits(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOutfitsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	outfits, err := h.svc.Player.GetGSFOutfits(r.Context(), req.PlayerAvatarOID, req.PlayerOID)
	if err != nil {
		return err
	}
	res := &messages.GetOutfitsResponse{}
	res.PlayerAvatarOutfits = outfits
	return w.Write(res)
}

// AddOutfit creates a new player outfit instance.
func (h *Handler) AddOutfit(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.AddOutfitRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	outfitOID, err := h.svc.Player.CreateGSFOutfit(r.Context(), req.PlayerAvatarOID, req.OutfitNo)
	if err != nil {
		return err
	}
	res := &messages.AddOutfitResponse{}
	res.PlayerAvatarOutfitOID = outfitOID
	res.OutfitNo = req.OutfitNo
	return w.Write(res)
}

// SetCurrentOutfit handles the active-outfit change request and returns the update status.
func (h *Handler) SetCurrentOutfit(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SetCurrentOutfitRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	err := h.svc.Player.SetGSFPlayerActiveOutfit(r.Context(), req.PlayerAvatarOID, req.PlayerAvatarOutfitOID, req.OutfitNo)
	if err != nil {
		return err
	}
	if playerOID, ok := r.PlayerOID(); ok {
		if err := h.syncHub.AppearanceChanged(types.OIDFromInt64(playerOID)); err != nil {
			return err
		}
	}
	res := &messages.SetCurrentOutfitResponse{}
	res.IsUpdated = true
	return w.Write(res)
}

// GetOutfitItems fetches the item instances associated with a player avatar outfit.
func (h *Handler) GetOutfitItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOutfitItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	items, err := h.svc.Player.GetGSFOutfitItems(r.Context(), r.Platform(), req.PlayerAvatarOutfitOID, req.PlayerOID)
	if err != nil {
		return err
	}
	res := &messages.GetOutfitItemsResponse{}
	res.OutfitItems = items
	return w.Write(res)
}

// AddOutfitItems handles the item assign request and returns the update status.
func (h *Handler) AddOutfitItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.AddOutfitItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	// TODO: ownership check, acceptable slot sheck, outfit oid check
	err := h.svc.Player.AddGSFOutfitItems(r.Context(), req.PlayerAvatarOutfitOID, req.InventoryOIDs, req.SlotOIDs)
	if err != nil {
		return err
	}
	if playerOID, ok := r.PlayerOID(); ok {
		if err := h.syncHub.AppearanceChanged(types.OIDFromInt64(playerOID)); err != nil {
			return err
		}
	}
	res := &messages.AddOutfitItemsResponse{}
	res.IsUpdated = true
	return w.Write(res)
}

// RemoveOutfitItems handles the item removal request and returns the update status.
func (h *Handler) RemoveOutfitItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.RemoveOutfitItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	// TODO: ownership check, outfit oid check
	err := h.svc.Player.RemoveGSFOutfitItems(r.Context(), req.PlayerAvatarOutfitOID, req.InventoryOIDs)
	if err != nil {
		return err
	}
	if playerOID, ok := r.PlayerOID(); ok {
		if err := h.syncHub.AppearanceChanged(types.OIDFromInt64(playerOID)); err != nil {
			return err
		}
	}
	res := &messages.RemoveOutfitItemsResponse{}
	res.IsUpdated = true
	return w.Write(res)
}

// ReplaceOutfitItems handles the item swap request and returns the update status.
func (h *Handler) ReplaceOutfitItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.ReplaceOutfitItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	// TODO: ownership check, acceptable slot sheck, outfit oid check
	err := h.svc.Player.ReplaceGSFOutfitItems(r.Context(), req.OldInventoryOIDs, req.NewInventoryOIDs)
	if err != nil {
		return err
	}
	if playerOID, ok := r.PlayerOID(); ok {
		if err := h.syncHub.AppearanceChanged(types.OIDFromInt64(playerOID)); err != nil {
			return err
		}
	}
	res := &messages.ReplaceOutfitItemsResponse{}
	res.IsUpdated = true
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

// GetPublicItemsByOIDs handles requests for public item definitions by OIDs.
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
	zones, err := h.svc.Zone.GetGSFZones(r.Context(), r.Platform())
	if err != nil {
		return err
	}
	res := &messages.GetZonesResponse{}
	res.Zones = zones
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

	if playerOID, ok := r.PlayerOID(); ok {
		if err := h.syncHub.SetLocation(types.OIDFromInt64(playerOID), req.LocOID); err != nil {
			return err
		}
	}

	// The Home.PlayerMaze.HomeTheme.AssetMap["Scene_Unity3D"] asset drives scene loading via
	// LoadMazeCommand.cs -> LoadMainScene() -> AssetDownloadManager.cs -> LoadMainScene()

	if req.LocOID.Int64() == 292733975781503755 {
		scene, err := h.svc.Asset.GetGSFAssetByCDNID(r.Context(), SCENE_SPRINGTIME)
		if err != nil {
			return err
		}

		container := types.AssetContainer{}
		container.AssetMap = types.AssetMap{}
		container.AssetMap["Scene_Unity3D"] = []types.Asset{scene}

		zoneInstance := &types.ZoneInstance{}
		zoneInstance.Zone.OID = types.OIDFromInt64(292733975781503755)
		zoneInstance.Zone.AssetContainer = container

		res.ZoneInstance = zoneInstance

	} else {
		cdnid := SCENE_HOMELOT_SMALL
		if rand.Intn(2) == 1 {
			cdnid = SCENE_HOMELOT_WINTER
		}

		scene, err := h.svc.Asset.GetGSFAssetByCDNID(r.Context(), cdnid)
		if err != nil {
			return err
		}

		container := types.AssetContainer{}
		container.AssetMap = types.AssetMap{}
		container.AssetMap["Scene_Unity3D"] = []types.Asset{scene}

		playerMaze := types.PlayerMaze{
			Name:       "coremaze",
			MazePieces: []types.PlayerMazePiece{},
			HomeTheme:  container,
		}

		res.Home = &types.PlayerHome{
			PlayerMaze:  playerMaze,
			HomeTheme:   container,
			PlayerMazes: []types.PlayerMaze{playerMaze},
		}
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
	r.SetPlayerOID(req.UOID.Int64())
	if err := h.syncHub.Join(r.Session(), req.UOID); err != nil {
		return err
	}
	res := &messages.SyncLoginResponse{}
	return w.Write(res)
}

// MovePlayer receives local player movement and relays it to other sync peers.
func (h *Handler) MovePlayer(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &notify.Move{}
	if err := r.Read(req); err != nil {
		return err
	}
	return h.syncHub.Move(r.Session(), req)
}

// HeartbeatNotify keeps the SYNC server session associated with the active player.
func (h *Handler) HeartbeatNotify(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &notify.Heartbeat{}
	if err := r.Read(req); err != nil {
		return err
	}
	r.SetPlayerOID(req.POID.Int64())
	return nil
}

// GetOtherPlayerDetails fetches the profile payload used by the client to spawn a remote player.
func (h *Handler) GetOtherPlayerDetails(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOtherPlayerDetailsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	details, err := h.svc.Player.GetGSFOtherPlayerDetails(r.Context(), r.Platform(), req.PlayerOID)
	if err != nil {
		return err
	}
	res := &messages.GetOtherPlayerDetailsResponse{}
	res.OtherPlayerDetails = details
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

// GetStoreItems fetches the items placed in store.
func (h *Handler) GetStoreItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetStoreItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetStoreItemsResponse{}
	res.StoreItems = []types.StoreItem{}
	return w.Write(res)
}

// GetPlayerQuests fetches the quests associated with NPC.
func (h *Handler) GetPlayerQuests(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPlayerQuestsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPlayerQuestsResponse{}
	res.PlayerQuests = []types.PlayerQuest{}
	return w.Write(res)
}

// GetQuestFromParent fetches the quests associated with NPC.
func (h *Handler) GetQuestFromParent(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetQuestFromParentRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetQuestFromParentResponse{}
	return w.Write(res)
}

// CreateQuest requests a quest initialization.
func (h *Handler) CreateQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.CreateQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.CreateQuestResponse{}
	return w.Write(res)
}

// AcceptQuest requests a quest acceptance.
func (h *Handler) AcceptQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.AcceptQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.AcceptQuestResponse{}
	return w.Write(res)
}

// StartQuest requests a quest start.
func (h *Handler) StartQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.StartQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.StartQuestResponse{}
	return w.Write(res)
}

// CompleteQuest requests a quest start.
func (h *Handler) CompleteQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.CompleteQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.CompleteQuestResponse{}
	res.AwardSet = []types.QuestAwardElement{}
	return w.Write(res)
}

// GetPlayerQuests fetches the quests by OIDs.
func (h *Handler) GetPlayerQuestsByOIDs(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPlayerQuestsByOIDsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPlayerQuestsByOIDsResponse{}
	res.PlayerQuests = []types.PlayerQuest{}
	return w.Write(res)
}
