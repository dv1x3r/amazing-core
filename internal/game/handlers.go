package game

import (
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/game/services"
	"github.com/dv1x3r/amazing-core/internal/game/worldsync"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/messages"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/notify"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/appcode"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/chatchanneltypevalue"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/queststate"
)

type Handler struct {
	svc      *services.Service
	syncHub  *worldsync.Hub
	sessions *SessionManager
}

func NewHandler(svc *services.Service) *Handler {
	return &Handler{
		svc:      svc,
		syncHub:  worldsync.NewHub(),
		sessions: NewSessionManager(),
	}
}

// session returns the session of the request connection, errNoSession before LOGIN or RELOGIN.
func (h *Handler) session(r *gsf.Request) (*Session, error) {
	if s := h.sessions.Get(r.Conn()); s != nil {
		return s, nil
	}
	return nil, errNoSession
}

// platform returns the platform of the session, the requests before LOGIN get the default one.
func (h *Handler) platform(r *gsf.Request) gsf.Platform {
	if s := h.sessions.Get(r.Conn()); s != nil {
		return s.platform
	}
	return gsf.PlatformUnknown
}

const (
	SCENE_HOMELOT_SMALL  = "OTYxMTQ4NDU5NDE5MA"
	SCENE_HOMELOT_WINTER = "OTQ1MDc3NTY0MjEyNg"
	SCENE_EMPTY          = "OTYwNzA1NTQwOTE2Ng"
)

const (
	CHAT_PRIVATE_OID       = 1
	CHAT_PRIVATE_GROUP_OID = 2
	CHAT_LOCAL_OID         = 3
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
	siteFrame, err := h.svc.GetSiteFrame(r.Context(), h.platform(r), req.TypeValue)
	if err != nil {
		return err
	}
	res := &messages.GetSiteFrameResponse{}
	res.AssetDeliveryURL = h.svc.DeliveryURL()
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

	playerOID, err := h.svc.LoginOrCreatePlayer(r.Context(), req.LoginID, req.Password)
	if err != nil {
		return wrap.WithGSFError(err, loginAppCode(err))
	}
	platform := gsf.ParsePlatformFromMachineOS(req.ClientEnvInfo.MachineOS)
	session := h.sessions.Login(r.Conn(), playerOID, platform, req.LoginID)

	player, err := h.svc.GetPlayer(r.Context(), session.platform, playerOID)
	if err != nil {
		return wrap.WithGSFError(err, appcode.INTERNAL_ERROR)
	}

	playerInfo, err := h.svc.GetPlayerInfoTO(r.Context(), playerOID)
	if err != nil {
		return wrap.WithGSFError(err, appcode.INTERNAL_ERROR)
	}

	res := &messages.LoginResponse{}
	res.SessionOID = session.sessionOID
	res.ConversationID = session.conversationID
	res.AssetDeliveryURL = h.svc.DeliveryURL()
	res.Player = player
	res.PlayerInfoTO = playerInfo

	return w.Write(res)
}

func loginAppCode(err error) appcode.AppCode {
	switch {
	case errors.Is(err, services.ErrBlankCredentials), errors.Is(err, services.ErrInvalidCredentials):
		return appcode.AUTH
	case errors.Is(err, services.ErrCreatePlayerAccountFailed):
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
	h.sessions.Logout(r.Conn())
	res := &messages.LogoutResponse{}
	return w.Write(res)
}

// Relogin binds the new USER connection to the session after the client reconnects.
// The client sends the session OID and conversation ID from the LOGIN response,
// an unknown session gets the SESSION app code and GSFSession switches to SessionExpired.
func (h *Handler) Relogin(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.ReloginRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	if !h.sessions.Relogin(r.Conn(), req.SessionOID, req.ConversationID, req.LoginName) {
		return wrap.WithGSFError(errors.New("session expired"), appcode.SESSION)
	}
	res := &messages.ReloginResponse{}
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
	names, err := h.svc.GetNStringsByType(r.Context(), req.NamePartType, int(req.Amount))
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
	categories, err := h.svc.GetItemCategories(r.Context(), true)
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
	categories, err := h.svc.GetItemCategories(r.Context(), false)
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
	session, err := h.session(r)
	if err != nil {
		return err
	}
	items, err := h.svc.GetInventoryObjects(r.Context(), session.platform, session.playerOID)
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
	session, err := h.session(r)
	if err != nil {
		return err
	}
	avatars, err := h.svc.GetPlayerAvatars(r.Context(), session.platform, session.playerOID)
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
	avatar, err := h.svc.SetPlayerActiveAvatar(r.Context(), h.platform(r), req.PlayerAvatarOID)
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
	outfits, err := h.svc.GetOutfits(r.Context(), req.PlayerAvatarOID, req.PlayerOID)
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
	outfitOID, err := h.svc.CreateOutfit(r.Context(), req.PlayerAvatarOID, req.OutfitNo)
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
	err := h.svc.SetPlayerActiveOutfit(r.Context(), req.PlayerAvatarOID, req.PlayerAvatarOutfitOID, req.OutfitNo)
	if err != nil {
		return err
	}
	if session, err := h.session(r); err == nil {
		if err := h.syncHub.AppearanceChanged(session.playerOID); err != nil {
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
	items, err := h.svc.GetOutfitItems(r.Context(), h.platform(r), req.PlayerAvatarOutfitOID, req.PlayerOID)
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
	err := h.svc.AddOutfitItems(r.Context(), req.PlayerAvatarOutfitOID, req.InventoryOIDs, req.SlotOIDs)
	if err != nil {
		return err
	}
	if session, err := h.session(r); err == nil {
		if err := h.syncHub.AppearanceChanged(session.playerOID); err != nil {
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
	err := h.svc.RemoveOutfitItems(r.Context(), req.PlayerAvatarOutfitOID, req.InventoryOIDs)
	if err != nil {
		return err
	}
	if session, err := h.session(r); err == nil {
		if err := h.syncHub.AppearanceChanged(session.playerOID); err != nil {
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
	err := h.svc.ReplaceOutfitItems(r.Context(), req.OldInventoryOIDs, req.NewInventoryOIDs)
	if err != nil {
		return err
	}
	if session, err := h.session(r); err == nil {
		if err := h.syncHub.AppearanceChanged(session.playerOID); err != nil {
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
	zones, err := h.svc.GetZones(r.Context(), h.platform(r))
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

	if session, err := h.session(r); err == nil {
		if err := h.syncHub.SetLocation(session.playerOID, req.LocOID); err != nil {
			return err
		}
	}

	// The Home.PlayerMaze.HomeTheme.AssetMap["Scene_Unity3D"] asset drives scene loading via
	// LoadMazeCommand.cs -> LoadMainScene() -> AssetDownloadManager.cs -> LoadMainScene()

	if req.LocOID.Int64() == 292733975781503755 {
		// The zone container has the Main_Scene asset and the spawn points of the task items,
		// SpawnPoints.ParseZone reads them from ZoneManager.currentZone set by this response.
		zone, err := h.svc.GetZone(r.Context(), h.platform(r), req.LocOID)
		if err != nil {
			return err
		}
		res.ZoneInstance = &types.ZoneInstance{Zone: zone}

	} else {
		cdnid := SCENE_HOMELOT_SMALL
		if rand.Intn(2) == 1 {
			cdnid = SCENE_HOMELOT_WINTER
		}

		scene, err := h.svc.GetAssetByCDNID(r.Context(), cdnid)
		if err != nil {
			return err
		}

		// EmptyScene has the shared scene objects (AngleViewCamera, RadialMenuManager, MiniMap),
		// LoadMazeCommand loads it additively after the Main_Scene.
		emptyScene, err := h.svc.GetAssetByCDNID(r.Context(), SCENE_EMPTY)
		if err != nil {
			return err
		}

		container := types.AssetContainer{}
		container.AssetMap = types.AssetMap{}
		container.AssetMap["Scene_Unity3D"] = []types.Asset{scene, emptyScene}

		// ZoneManager takes the yard owner from PlayerHome.PlayerOID and HUDManager requests its details.
		session, err := h.session(r)
		if err != nil {
			return err
		}

		playerMaze := types.PlayerMaze{
			Name:       "coremaze",
			PlayerOID:  session.playerOID,
			MazePieces: []types.PlayerMazePiece{},
			HomeTheme:  container,
		}

		res.Home = &types.PlayerHome{
			PlayerMaze:  playerMaze,
			HomeTheme:   container,
			PlayerOID:   session.playerOID,
			PlayerMazes: []types.PlayerMaze{playerMaze},
		}
	}

	return w.Write(res)
}

// SyncLogin authenticates the player on the SYNC server.
// The SYNC server handles real-time positional and social updates.
// It also handles RELOGIN, the client sends it with the same payload after the connection is restored.
func (h *Handler) SyncLogin(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SyncLoginRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	if !h.sessions.AttachSync(r.Conn(), req.UOID) {
		return errNoSession
	}
	if err := h.syncHub.Join(r.Conn(), req.UOID); err != nil {
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
	return h.syncHub.Move(r.Conn(), req)
}

// SyncAction receives a local player action and relays it to other sync peers.
func (h *Handler) SyncAction(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &notify.Action{}
	if err := r.Read(req); err != nil {
		return err
	}
	throttled, err := h.syncHub.Action(r.Conn(), req)
	if throttled {
		r.MarkDropped()
	}
	return err
}

// HeartbeatNotify keeps the SYNC connection bound to the session of the player.
func (h *Handler) HeartbeatNotify(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &notify.Heartbeat{}
	if err := r.Read(req); err != nil {
		return err
	}
	h.sessions.AttachSync(r.Conn(), req.POID)
	return nil
}

// GetOtherPlayerDetails fetches the profile payload used by the client to spawn a remote player.
func (h *Handler) GetOtherPlayerDetails(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOtherPlayerDetailsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	details, err := h.svc.GetOtherPlayerDetails(r.Context(), h.platform(r), req.PlayerOID)
	if errors.Is(err, sql.ErrNoRows) {
		// PlayerInfoCache handles the error response (OtherPlayerDetailErrorHandler).
		return wrap.WithGSFError(err, appcode.NOT_FOUND)
	}
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
	res.ChatChannelTypes = []types.ChatChannelType{
		{OID: types.OIDFromInt64(CHAT_PRIVATE_OID), Value: chatchanneltypevalue.PRIVATE},
		{OID: types.OIDFromInt64(CHAT_PRIVATE_GROUP_OID), Value: chatchanneltypevalue.PRIVATE_GROUP},
		{OID: types.OIDFromInt64(CHAT_LOCAL_OID), Value: chatchanneltypevalue.LOCAL},
	}
	return w.Write(res)
}

// SendMessage accepts a chat message and broadcasts local chat to players in same location.
func (h *Handler) SendMessage(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SendMessageRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	message := req.Message
	if config.Get().Settings.BroChat {
		message = "bro"
	}

	res := &messages.SendMessageResponse{
		FilteredMessage: message,
		Categories:      []string{},
		OffenceGroup:    0,
		RecipientID:     req.RecipientID,
	}

	localChatTypeOID := types.OIDFromInt64(CHAT_LOCAL_OID)
	if req.Message != "" && req.TypeID == localChatTypeOID {
		session, err := h.session(r)
		if err != nil {
			return err
		}
		senderOID := session.playerOID
		if err := h.syncHub.Chat(senderOID, &notify.Chat{
			Msg:          res.FilteredMessage,
			TypeID:       localChatTypeOID,
			GroupID:      types.OID{},
			SenderID:     senderOID,
			Categories:   res.Categories,
			OffenceGroup: res.OffenceGroup,
		}); err != nil {
			return err
		}
	}

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
	npcs, err := h.svc.GetZoneNPCs(r.Context(), h.platform(r), req.ZoneOID)
	if err != nil {
		return err
	}
	res := &messages.GetPlayerNPCsResponse{}
	res.NPCs = npcs
	return w.Write(res)
}

// GetRuleCountList fetches rule counter values for the requested objects.
func (h *Handler) GetRuleCountList(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetRuleCountListRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetRuleCountListResponse{}
	res.ItemPairs = make([]types.OIDLongPair, 0, len(req.Items))
	for i, oid := range req.Items {
		pair := types.OIDLongPair{A: oid}
		// NPCsDownloader looks for the "SpawnPoint_<count>" location asset,
		// and falls back to "SpawnPoint_Default" or the first one.
		if i < len(req.Expressions) && req.Expressions[i] == "Counter.NPCSpawnPoint" {
			pair.B = gsf.NewNullable[int64](0)
		}
		// TODO: other expressions (Counter.ConstructionStage, wheel respawn timers) are not evaluated yet
		res.ItemPairs = append(res.ItemPairs, pair)
	}
	return w.Write(res)
}

// GetRuleInfo fetches the rules of an object. Task requests them for a quest before the quest dialog.
func (h *Handler) GetRuleInfo(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetRuleInfoRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetRuleInfoResponse{}
	// TODO: quest awards (AwardSets) for the quest dialog
	// Task.RuleInfoHandler continues only with the "Maps" group.
	res.Property.ChildrenGroup = map[string][]types.RuleProperty{"Maps": {}}
	res.ObjectOID = req.ObjectOID
	return w.Write(res)
}

// GetRuleInstanceInfo fetches the rules of an instance. Task requests them for a player quest on accept.
func (h *Handler) GetRuleInstanceInfo(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetRuleInstanceInfoRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetRuleInstanceInfoResponse{}
	res.Property.ChildrenGroup = map[string][]types.RuleProperty{"Maps": {}}
	res.InstanceOID = req.InstanceOID
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

// GetQuestFromParent fetches the quest offered by the parent (e.g. NPC).
func (h *Handler) GetQuestFromParent(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetQuestFromParentRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	isNPC, err := h.svc.IsNPC(r.Context(), req.ParentOID)
	if err != nil {
		return err
	}
	res := &messages.GetQuestFromParentResponse{}
	// TODO: MobquestController, MultiplayerQuestTrigger do not handle a null quest;
	if !isNPC {
		res.Quest = &types.Quest{}
		return w.Write(res)
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	// NPCController expects a null quest when the NPC offers nothing.
	res.Quest, err = h.svc.GetNPCQuest(r.Context(), session.platform, session.playerOID, req.ParentOID)
	if err != nil {
		return err
	}
	return w.Write(res)
}

// CreateQuest creates a player quest for the quest the player has taken.
func (h *Handler) CreateQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.CreateQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	playerQuestOID, err := h.svc.CreatePlayerQuest(r.Context(), session.playerOID, req.QuestOID)
	if err != nil {
		return err
	}
	res := &messages.CreateQuestResponse{}
	res.PlayerQuestOID = playerQuestOID
	// NPC quests have no children, TaskController reads the map only for the quest children.
	res.SubQuestMap = map[string]types.OID{}
	return w.Write(res)
}

// AcceptQuest accepts the player quest.
func (h *Handler) AcceptQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.AcceptQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	if err := h.svc.SetPlayerQuestState(r.Context(), session.playerOID, req.PlayerQuestOID, queststate.ACCEPTED); err != nil {
		return err
	}
	res := &messages.AcceptQuestResponse{}
	res.IsAccepted = true
	res.CanStart = true
	res.CanPlay = true
	res.JoinedPlayers = []types.OID{}
	return w.Write(res)
}

// StartQuest starts the player quest, the spawned items are sent back as is.
func (h *Handler) StartQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.StartQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	if err := h.svc.SetPlayerQuestState(r.Context(), session.playerOID, req.PlayerQuestOID, queststate.STARTED); err != nil {
		return err
	}
	res := &messages.StartQuestResponse{}
	res.IsStarted = true
	res.SpawnedItems = req.SpawnedItems
	res.ConfirmedPlayers = []types.OID{}
	return w.Write(res)
}

// AbandonQuest closes the player quest the player has quit or failed, the NPC offers it again.
func (h *Handler) AbandonQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.AbandonQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	state := queststate.CLOSED_ABANDONED
	if req.IsFailed {
		state = queststate.FAILED
	}
	if err := h.svc.SetPlayerQuestState(r.Context(), session.playerOID, req.PlayerQuestOID, state); err != nil {
		return err
	}
	res := &messages.AbandonQuestResponse{}
	res.IsAbandoned = true
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

// GetAssetsByOIDs fetches the asset containers by OIDs, e.g. the quest items.
func (h *Handler) GetAssetsByOIDs(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAssetsByOIDsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	assets, err := h.svc.GetAssetContainersByOIDs(r.Context(), h.platform(r), req.OIDs)
	if err != nil {
		return err
	}
	res := &messages.GetAssetsByOIDsResponse{}
	res.Assets = assets
	return w.Write(res)
}
