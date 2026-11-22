package game

import (
	"database/sql"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/game/services"
	"github.com/dv1x3r/amazing-core/internal/game/worldsync"
	"github.com/dv1x3r/amazing-core/internal/lib/ratelimit"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/messages"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/notify"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/appcode"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/chatchanneltypevalue"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/clientmessagetype"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/notificationtype"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/queststate"

	"golang.org/x/time/rate"
)

type Handler struct {
	svc      *services.Service
	world    *worldsync.World
	sessions *SessionManager

	// loginLimiter limits LOGIN by IP, every LOGIN hashes the password and can create a player.
	loginLimiter *ratelimit.Limiter
}

func NewHandler(svc *services.Service, sessions *SessionManager, world *worldsync.World) *Handler {
	return &Handler{
		svc:          svc,
		world:        world,
		sessions:     sessions,
		loginLimiter: ratelimit.New(rate.Every(time.Second), 5, 10*time.Minute),
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

// The chat channel type OIDs of GET_CHAT_CHANNEL_TYPES.
const (
	chatPrivateOID      = 1
	chatPrivateGroupOID = 2
	chatLocalOID        = 3
)

// maxChatMessageLength is in characters, the real limit of the client chat field is unknown.
const maxChatMessageLength = 256

// ── General ──────────────────────────────────────────────────────────────────

// GetClientVersionInfo lets the game start, the client checks its version against the reply on launch.
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

// GetSiteFrame gives the client the core assets and the CDN address it loads every other asset from.
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
	res.AssetDeliveryURL = config.Get().Settings.AssetDeliveryURL
	res.SiteFrame = siteFrame
	return w.Write(res)
}

// Login opens the player session and gives the client its player, active avatar and asset address.
// It also stores the client platform, the asset containers are picked by it.
func (h *Handler) Login(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.LoginRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	if !h.loginLimiter.Allow(r.Conn().RemoteIP()) {
		return wrap.WithGSFError(errors.New("too many login attempts"), appcode.AUTH)
	}

	playerOID, err := h.svc.LoginOrCreatePlayer(r.Context(), req.LoginID, req.Password)
	if err != nil {
		return wrap.WithGSFError(err, loginAppCode(err))
	}
	platform := gsf.ParsePlatformFromMachineOS(req.ClientEnvInfo.MachineOS)

	player, err := h.svc.GetPlayer(r.Context(), platform, playerOID)
	if err != nil {
		return wrap.WithGSFError(err, appcode.INTERNAL_ERROR)
	}

	playerInfo, err := h.svc.GetPlayerInfoTO(r.Context(), playerOID)
	if err != nil {
		return wrap.WithGSFError(err, appcode.INTERNAL_ERROR)
	}

	// RELOGIN sends back the player name of the response as the login name.
	session := h.sessions.Login(r.Conn(), h.world.NewPlayer(playerOID), platform, playerInfo.PlayerName)

	res := &messages.LoginResponse{}
	res.SessionOID = session.sessionOID
	res.ConversationID = session.conversationID
	res.AssetDeliveryURL = config.Get().Settings.AssetDeliveryURL
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

// Logout closes the session of the connection when the player leaves the game.
func (h *Handler) Logout(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.LogoutRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	h.sessions.Logout(r.Conn())
	res := &messages.LogoutResponse{}
	return w.Write(res)
}

// Relogin keeps the session after a reconnect by binding the new USER connection to it.
// An unknown session gets the SESSION app code, and the client treats the session as expired.
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

// Heartbeat keeps the USER connection alive and gives the client the server time.
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

// GetTiers answers the login with no subscription tiers, tiers are not implemented.
func (h *Handler) GetTiers(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetTiersRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetTiersResponse{}
	res.Tiers = []types.Tier{}
	return w.Write(res)
}

// GetRequiredExperience lets the level-up reward play, RewardSequenceScriptLevelUp reads the asset map of the level.
// Experience is not implemented, so the level is an empty container.
func (h *Handler) GetRequiredExperience(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetRequiredExperienceRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetRequiredExperienceResponse{}
	res.PlayerExperienceLevel.Level = req.Level
	res.PlayerExperienceLevel.AssetMap = types.AssetMap{}
	res.PlayerExperienceLevel.AssetPackages = []types.AssetPackage{}
	return w.Write(res)
}

// GetStatsType answers PlayerStatsUpdater with no statistic types, statistics are not implemented.
func (h *Handler) GetStatsType(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetStatsTypeRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.GetStatsTypeResponse{})
}

// ── Registration ─────────────────────────────────────────────────────────────

// GetRandomNames offers random Zing names and family name parts at registration and on Zing rename.
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

// ValidateName lets the client accept a new Zing name, the name filter is a stub.
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

// SelectPlayerName lets the registration continue, every family name is accepted.
func (h *Handler) SelectPlayerName(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SelectPlayerNameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.SelectPlayerNameResponse{}
	return w.Write(res)
}

// CheckUsername lets the registration continue, every username is free for now.
func (h *Handler) CheckUsername(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.CheckUsernameRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.CheckUsernameResponse{}
	return w.Write(res)
}

// RegisterPlayer lets the registration finish, the player is created at the first login.
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

// GetPublicItemCategories lets the registration classify the temporary items of the new player.
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

// GetCMSItemCategories lets the client map item categories to item types like Clothing, Decoration or Yard,
// the inventory and the item placement depend on it.
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

// GetBuildObjects answers the builder with no maze pieces, maze pieces are not implemented.
func (h *Handler) GetBuildObjects(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetBuildObjectsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetBuildObjectsResponse{}
	res.PlayerBuildObjects = []types.PlayerBuildObject{}
	return w.Write(res)
}

// GetInventoryObjects fills the inventory with the items the player has not placed, attached or worn.
func (h *Handler) GetInventoryObjects(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetInventoryObjectsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	items, err := h.svc.GetInventoryObjects(r.Context(), session.platform, session.player.OID)
	if err != nil {
		return err
	}
	res := &messages.GetInventoryObjectsResponse{}
	res.PlayerItems = items
	return w.Write(res)
}

// GetPublicItemsByOIDs answers with no item definitions, the request is not implemented.
func (h *Handler) GetPublicItemsByOIDs(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPublicItemsByOIDsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPublicItemsByOIDsResponse{}
	res.Items = []types.Item{}
	return w.Write(res)
}

// ── Avatars ──────────────────────────────────────────────────────────────────

// GetAvatars gives the client the Zings of the player for the avatar switcher.
func (h *Handler) GetAvatars(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAvatarsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	avatars, err := h.svc.GetPlayerAvatars(r.Context(), session.platform, session.player.OID)
	if err != nil {
		return err
	}
	res := &messages.GetAvatarsResponse{}
	res.Avatars = avatars
	return w.Write(res)
}

// GetAvatarItems answers an old clothing request with no items, the clothes come from GetOutfitItems.
func (h *Handler) GetAvatarItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAvatarItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetAvatarItemsResponse{}
	res.AvatarItems = []types.PlayerItem{}
	return w.Write(res)
}

// UpdatePlayerActiveAvatar saves the Zing the player switched to and shows it to the players nearby.
func (h *Handler) UpdatePlayerActiveAvatar(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.UpdatePlayerActiveAvatarRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	avatar, err := h.svc.SetPlayerActiveAvatar(r.Context(), session.platform, session.player.OID, req.PlayerAvatarOID)
	if err != nil {
		return err
	}
	h.world.AppearanceChanged(session.player)
	res := &messages.UpdatePlayerActiveAvatarResponse{}
	res.ActivePlayerAvatar = avatar
	return w.Write(res)
}

// ── Outfits ──────────────────────────────────────────────────────────────────

// GetOutfits gives the dressing room the saved outfits of a Zing.
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

// AddOutfit saves a new outfit of a Zing in the dressing room.
func (h *Handler) AddOutfit(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.AddOutfitRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	outfitOID, err := h.svc.CreateOutfit(r.Context(), session.player.OID, req.PlayerAvatarOID, req.OutfitNo)
	if err != nil {
		return err
	}
	res := &messages.AddOutfitResponse{}
	res.PlayerAvatarOutfitOID = outfitOID
	res.OutfitNo = req.OutfitNo
	return w.Write(res)
}

// SetCurrentOutfit saves which outfit the Zing wears.
func (h *Handler) SetCurrentOutfit(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SetCurrentOutfitRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	if err := h.svc.SetPlayerActiveOutfit(r.Context(), session.player.OID, req.PlayerAvatarOID, req.OutfitNo); err != nil {
		return err
	}
	h.world.AppearanceChanged(session.player)
	res := &messages.SetCurrentOutfitResponse{}
	res.IsUpdated = true
	return w.Write(res)
}

// GetOutfitItems tells the client what a Zing wears in an outfit.
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

// AddOutfitItems saves the clothes the player puts on.
func (h *Handler) AddOutfitItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.AddOutfitItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	// TODO: acceptable slot check
	if err := h.svc.AddOutfitItems(r.Context(), session.player.OID, req.PlayerAvatarOutfitOID, req.InventoryOIDs, req.SlotOIDs); err != nil {
		return err
	}
	h.world.AppearanceChanged(session.player)
	res := &messages.AddOutfitItemsResponse{}
	res.IsUpdated = true
	return w.Write(res)
}

// RemoveOutfitItems saves the clothes the player takes off.
func (h *Handler) RemoveOutfitItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.RemoveOutfitItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	// TODO: outfit oid check
	if err := h.svc.RemoveOutfitItems(r.Context(), session.player.OID, req.InventoryOIDs); err != nil {
		return err
	}
	h.world.AppearanceChanged(session.player)
	res := &messages.RemoveOutfitItemsResponse{}
	res.IsUpdated = true
	return w.Write(res)
}

// ReplaceOutfitItems saves the clothes the player swaps in the same slots.
func (h *Handler) ReplaceOutfitItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.ReplaceOutfitItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	// TODO: acceptable slot check, outfit oid check
	if err := h.svc.ReplaceOutfitItems(r.Context(), session.player.OID, req.OldInventoryOIDs, req.NewInventoryOIDs); err != nil {
		return err
	}
	h.world.AppearanceChanged(session.player)
	res := &messages.ReplaceOutfitItemsResponse{}
	res.IsUpdated = true
	return w.Write(res)
}

// ── Location ─────────────────────────────────────────────────────────────────

// GetZones gives the client the zones it loads the world and the NPC spawn points from.
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

// InitLocation loads the location of the player and tells the client where the SYNC server is.
func (h *Handler) InitLocation(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.InitLocationRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	session, err := h.session(r)
	if err != nil {
		return err
	}
	h.world.Leave(session.player)

	res := &messages.InitLocationResponse{}
	res.SyncServerToken = session.syncToken
	res.SyncServerIP = config.Get().Settings.SyncServerIP
	res.SyncServerPort = int32(config.Get().Settings.SyncServerPort)

	// The zone container has the Main_Scene asset and the spawn points of the task items,
	// SpawnPoints.ParseZone reads them from ZoneManager.currentZone set by this response.
	zone, err := h.svc.GetZone(r.Context(), session.platform, req.LocOID)
	if err == nil {
		// One instance per zone, its OID is the zone OID, like ZoneManager.defaultCurrentHomeZone.
		res.ZoneInstance = &types.ZoneInstance{OID: req.LocOID, Zone: zone}
		return w.Write(res)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// Not a zone: the location of a home is the OID of its owner, the visitors send it too.
	res.Home, err = h.svc.GetPlayerHome(r.Context(), session.platform, req.LocOID, config.Get().Settings.Winter)
	if errors.Is(err, sql.ErrNoRows) {
		return wrap.WithGSFError(err, appcode.NOT_FOUND)
	}
	if err != nil {
		return err
	}
	return w.Write(res)
}

// GetMazeItems furnishes the yard or a house interior, for the owner and for the visitors.
func (h *Handler) GetMazeItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetMazeItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	items, err := h.svc.GetMazeItems(r.Context(), h.platform(r), req.PlayerMazeOID, req.PlayerOID)
	if err != nil {
		return err
	}
	res := &messages.GetMazeItemsResponse{}
	res.MazeItems = items
	return w.Write(res)
}

// PlaceMazeItem saves where the player put an item in the yard or a house.
// IsUpdated is false for an item or a maze of another player.
func (h *Handler) PlaceMazeItem(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.PlaceMazeItemRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	err = h.svc.PlaceMazeItem(r.Context(), session.player.OID, req.InventoryOID, req.PlayerMazeOID, req.ObjectPosition, req.Rotation)
	if err != nil && !errors.Is(err, services.ErrMazeItemNotFound) {
		return err
	}
	return w.Write(&messages.PlaceMazeItemResponse{IsUpdated: err == nil})
}

// RemoveMazeItem saves that the player took an item back to the inventory.
func (h *Handler) RemoveMazeItem(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.RemoveMazeItemRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	err = h.svc.RemoveMazeItem(r.Context(), session.player.OID, req.InventoryOID)
	if err != nil && !errors.Is(err, services.ErrMazeItemNotFound) {
		return err
	}
	return w.Write(&messages.RemoveMazeItemResponse{IsUpdated: err == nil})
}

// RemoveAndDetachMazeItem saves that the player took an item off its parent, e.g. when a house is emptied.
func (h *Handler) RemoveAndDetachMazeItem(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.RemoveAndDetachMazeItemRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	err = h.svc.RemoveMazeItem(r.Context(), session.player.OID, req.InventoryOID)
	if err != nil && !errors.Is(err, services.ErrMazeItemNotFound) {
		return err
	}
	return w.Write(&messages.RemoveAndDetachMazeItemResponse{IsUpdated: err == nil})
}

// DetachItems takes the upgrades off a house before its sale, the client sells them first.
func (h *Handler) DetachItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.DetachItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	if err := h.svc.DetachMazeItems(r.Context(), session.player.OID, req.ParentOID, req.ChildItems); err != nil {
		return err
	}
	return w.Write(&messages.DetachItemsResponse{})
}

// AttachSingleItem saves a door or a window the player put on a house.
func (h *Handler) AttachSingleItem(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.AttachSingleItemRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	err = h.svc.AttachMazeItem(r.Context(), session.player.OID, req.ParentOID, req.ChildOID, req.Ordinal)
	if errors.Is(err, services.ErrMazeItemNotFound) {
		return wrap.WithGSFError(err, appcode.NOT_FOUND)
	}
	if err != nil {
		return err
	}
	return w.Write(&messages.AttachSingleItemResponse{
		ParentOID:  req.ParentOID,
		ChildOID:   req.ChildOID,
		Ordinal:    req.Ordinal,
		OrdinalUpd: req.Ordinal,
	})
}

// EnterBuilding moves the player into the loaded location, so the players nearby see it at the spawn point.
func (h *Handler) EnterBuilding(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.EnterBuildingRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	if session, err := h.session(r); err == nil {
		h.world.Enter(session.player, req.LocOID, req.Pos, req.Orientation)
	}
	res := &messages.EnterBuildingResponse{}
	return w.Write(res)
}

// GetAllBuildings answers the zone load with no buildings, buildings are not implemented.
func (h *Handler) GetAllBuildings(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetAllBuildingsRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.GetAllBuildingsResponse{Buildings: []types.Building{}})
}

// ── Social ───────────────────────────────────────────────────────────────────

// GetOtherPlayerDetails lets the client draw another player: its Zing, clothes and name.
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

// GetChatChannelTypes enables the chat, the client maps the channel OIDs to its channels.
func (h *Handler) GetChatChannelTypes(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetChatChannelTypesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetChatChannelTypesResponse{}
	res.ChatChannelTypes = []types.ChatChannelType{
		{OID: types.OIDFromInt64(chatPrivateOID), Value: chatchanneltypevalue.PRIVATE},
		{OID: types.OIDFromInt64(chatPrivateGroupOID), Value: chatchanneltypevalue.PRIVATE_GROUP},
		{OID: types.OIDFromInt64(chatLocalOID), Value: chatchanneltypevalue.LOCAL},
	}
	return w.Write(res)
}

// SendMessage delivers a local chat message to the players in the same location.
func (h *Handler) SendMessage(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.SendMessageRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	// Local chat goes to every player nearby, a modified client could send megabytes.
	if utf8.RuneCountInString(req.Message) > maxChatMessageLength {
		return wrap.WithGSFError(errors.New("chat message is too long"), appcode.INPUT)
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

	localChatTypeOID := types.OIDFromInt64(chatLocalOID)
	if req.Message != "" && req.TypeID == localChatTypeOID {
		session, err := h.session(r)
		if err != nil {
			return err
		}
		h.world.Chat(session.player, &notify.Chat{
			Msg:          res.FilteredMessage,
			TypeID:       localChatTypeOID,
			GroupID:      types.OID{},
			SenderID:     session.player.OID,
			Categories:   res.Categories,
			OffenceGroup: res.OffenceGroup,
		})
	}

	return w.Write(res)
}

// GetAnnouncements answers the login with no announcements, they are not implemented.
func (h *Handler) GetAnnouncements(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAnnouncementsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetAnnouncementsResponse{}
	res.Announcements = []types.Announcement{}
	return w.Write(res)
}

// GetOnlineStatuses answers the friend list with no online statuses, friends are not implemented.
func (h *Handler) GetOnlineStatuses(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetOnlineStatusesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetOnlineStatusesResponse{}
	res.OnlineStatuses = []types.OnlineStatus{}
	return w.Write(res)
}

// GetFriendshipRequestCounts answers the friend list with no requests, friends are not implemented.
func (h *Handler) GetFriendshipRequestCounts(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetFriendshipRequestCountsRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.GetFriendshipRequestCountsResponse{})
}

// ListLimits answers the friend list with no limits, friends are not implemented.
func (h *Handler) ListLimits(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.ListLimitsRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.ListLimitsResponse{})
}

// GetHomeInvitations answers the inbox with no home invitations, they are not implemented.
func (h *Handler) GetHomeInvitations(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetHomeInvitationsRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.GetHomeInvitationsResponse{})
}

// GetNotificationByPlayerID answers the inbox with no stored notifications, they are not implemented.
func (h *Handler) GetNotificationByPlayerID(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetNotificationByPlayerIDRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.GetNotificationByPlayerIDResponse{})
}

// GetCrispActions answers with no Crisp actions, they are not implemented.
func (h *Handler) GetCrispActions(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetCrispActionsRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.GetCrispActionsResponse{})
}

// ── NPC and Quests ───────────────────────────────────────────────────────────

// GetPlayerNPCs spawns the NPCs of the zone.
// The client parses the spawn points of the zone from GetZones before it.
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

// GetRuleCountList gives the rule counters, NPCsDownloader picks the NPC spawn point by one.
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

// GetRuleInfo lets the quest dialog open, the quest rules are not sent yet.
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

// GetRuleInstanceInfo lets the player accept a quest, the quest rules are not sent yet.
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

// GetNPCRelationships fills the relationship meters of the NPCs, on the zone load and after a reward.
func (h *Handler) GetNPCRelationships(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetNPCRelationshipsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	relationships, err := h.svc.GetNPCRelationships(r.Context(), session.player.OID, req.ZoneOID)
	if err != nil {
		return err
	}
	return w.Write(&messages.GetNPCRelationshipsResponse{Relationships: relationships})
}

// GetNPCRelationshipLevels answers with no relationship levels, their points are unknown.
func (h *Handler) GetNPCRelationshipLevels(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetNPCRelationshipLevelsRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.GetNPCRelationshipLevelsResponse{Levels: []types.NPCRelationshipLevel{}})
}

// GetNPCsWithQuestOffer lets the task alert show the portrait of an NPC with a quest for the player.
func (h *Handler) GetNPCsWithQuestOffer(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetNPCsWithQuestOfferRequest{}); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	npcs, err := h.svc.GetNPCsWithQuestOffer(r.Context(), session.platform, session.player.OID)
	if err != nil {
		return err
	}
	return w.Write(&messages.GetNPCsWithQuestOfferResponse{NPCs: npcs})
}

// GetPlayerQuests answers the NPC with no player quests, the quest history is not implemented.
func (h *Handler) GetPlayerQuests(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPlayerQuestsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPlayerQuestsResponse{}
	res.PlayerQuests = []types.PlayerQuest{}
	return w.Write(res)
}

// GetPlayerQuestsByOIDs answers with no player quests, the quest history is not implemented.
func (h *Handler) GetPlayerQuestsByOIDs(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetPlayerQuestsByOIDsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetPlayerQuestsByOIDsResponse{}
	res.PlayerQuests = []types.PlayerQuest{}
	return w.Write(res)
}

// GetQuestFromParent tells the client which quest an NPC offers now.
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
	res.Quest, err = h.svc.GetNPCQuest(r.Context(), session.platform, session.player.OID, req.ParentOID)
	if err != nil {
		return err
	}
	return w.Write(res)
}

// CreateQuest starts the quest the player took: the player quest and the player quests of its steps.
func (h *Handler) CreateQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.CreateQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	playerQuestOID, subQuests, err := h.svc.CreatePlayerQuest(r.Context(), session.player.OID, req.QuestOID)
	if err != nil {
		if errors.Is(err, services.ErrQuestNotFound) {
			return wrap.WithGSFError(err, appcode.NOT_FOUND)
		}
		return err
	}
	res := &messages.CreateQuestResponse{}
	res.PlayerQuestOID = playerQuestOID
	// TaskController takes the steps of a parent quest that are in the map.
	res.SubQuestMap = subQuests
	return w.Write(res)
}

// AcceptQuest saves that the player accepted the quest.
func (h *Handler) AcceptQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.AcceptQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	if err := h.svc.SetPlayerQuestState(r.Context(), session.player.OID, req.PlayerQuestOID, queststate.ACCEPTED); err != nil {
		if errors.Is(err, services.ErrQuestNotFound) {
			return wrap.WithGSFError(err, appcode.NOT_FOUND)
		}
		return err
	}
	res := &messages.AcceptQuestResponse{}
	res.IsAccepted = true
	res.CanStart = true
	res.CanPlay = true
	res.JoinedPlayers = []types.OID{}
	return w.Write(res)
}

// StartQuest saves the start of the quest, the client spawns the quest items from the reply.
func (h *Handler) StartQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.StartQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	if err := h.svc.SetPlayerQuestState(r.Context(), session.player.OID, req.PlayerQuestOID, queststate.STARTED); err != nil {
		if errors.Is(err, services.ErrQuestNotFound) {
			return wrap.WithGSFError(err, appcode.NOT_FOUND)
		}
		return err
	}
	res := &messages.StartQuestResponse{}
	res.IsStarted = true
	res.SpawnedItems = req.SpawnedItems
	res.ConfirmedPlayers = []types.OID{}
	return w.Write(res)
}

// AddQuestItem counts a collected quest item for the task counter.
// When every amount is collected, the QuestCompleted notification follows the response.
func (h *Handler) AddQuestItem(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.AddQuestItemRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	completed, err := h.svc.AddPlayerQuestItem(r.Context(), session.player.OID, req.PlayerQuestOID, req.ItemOID, req.Count)
	if err != nil {
		return err
	}
	if completed {
		w.Notify(clientmessagetype.NOTIFICATION, &notify.Notification{
			OID:       types.OID{Number: int64(notificationtype.QUEST_COMPLETED)},
			POID:      session.player.OID,
			ObjectOID: req.PlayerQuestOID,
		})
	}
	res := &messages.AddQuestItemResponse{}
	res.QuestInventoryOIDs = []types.OID{}
	return w.Write(res)
}

// CompleteQuest closes the player quest when the player takes the reward.
// After the last step, the QuestCompleted notification of its quest follows the response.
func (h *Handler) CompleteQuest(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.CompleteQuestRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	closed, parentOID, err := h.svc.ClosePlayerQuest(r.Context(), session.player.OID, req.PlayerQuestOID)
	if err != nil {
		return err
	}
	if parentOID != nil {
		w.Notify(clientmessagetype.NOTIFICATION, &notify.Notification{
			OID:       types.OID{Number: int64(notificationtype.QUEST_COMPLETED)},
			POID:      session.player.OID,
			ObjectOID: *parentOID,
		})
	}
	res := &messages.CompleteQuestResponse{}
	res.IsCompleted = closed
	res.AwardSet = []types.QuestAwardElement{}
	return w.Write(res)
}

// AbandonQuest drops or fails the player quest and its open steps, so the NPC offers the quest again.
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
	if err := h.svc.AbandonPlayerQuest(r.Context(), session.player.OID, req.PlayerQuestOID, state); err != nil {
		if errors.Is(err, services.ErrQuestNotFound) {
			return wrap.WithGSFError(err, appcode.NOT_FOUND)
		}
		return err
	}
	res := &messages.AbandonQuestResponse{}
	res.IsAbandoned = true
	return w.Write(res)
}

// GetPlayerMissions answers with no missions, achievements and collections are not implemented.
func (h *Handler) GetPlayerMissions(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetPlayerMissionsRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.GetPlayerMissionsResponse{PlayerMissions: []types.PlayerQuest{}})
}

// GetCMSMissions answers with no mission definitions, missions are not implemented.
func (h *Handler) GetCMSMissions(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetCMSMissionsRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.GetCMSMissionsResponse{Missions: []types.Quest{}})
}

// QuestEventInProgress answers with no mob quest events, mob quests are not implemented.
func (h *Handler) QuestEventInProgress(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.QuestEventInProgressRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.QuestEventInProgressResponse{})
}

// AddEvent accepts a rule engine event of the player, e.g. an item use, events give no rewards yet.
func (h *Handler) AddEvent(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.AddEventRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.AddEventResponse{QuestAwardElements: []types.QuestAwardElement{}})
}

// GetCMSNotifications gives the notification definitions the client shows the notifications with.
func (h *Handler) GetCMSNotifications(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetCMSNotificationsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetCMSNotificationsResponse{}
	// The original OIDs are unknown, the OID of a definition is its type.
	first, last := notificationtype.FRIEND_REQUEST_PENDING, notificationtype.QUEST_MORE_PLAYERS_NEEDED
	res.Notifications = make([]types.Notification, 0, last-first+1)
	for t := first; t <= last; t++ {
		res.Notifications = append(res.Notifications, types.Notification{
			OID:              types.OID{Number: int64(t)},
			Category:         types.NotificationCategory{OID: types.OID{Number: 1}},
			NotificationType: t,
		})
	}
	return w.Write(res)
}

// GetAssetsByOIDs lets the client load the asset containers it knows only by OID, e.g. the quest items.
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

// ── Store ────────────────────────────────────────────────────────────────────

// GetStoreItems fills a store window with its products.
func (h *Handler) GetStoreItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetStoreItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	// The only store of this message is a house of the player selling its upgrades, other stores get none.
	items, err := h.svc.GetHouseUpgrades(r.Context(), session.platform, session.player.OID, req.StoreOID)
	if err != nil {
		return err
	}
	res := &messages.GetStoreItemsResponse{}
	res.StoreItems = items
	return w.Write(res)
}

// ListStores gives an NPC its store, NPCController shows the NPC Shop button only when there is one.
// Every NPC has the same store with every item.
func (h *Handler) ListStores(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.ListStoresRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.ListStoresResponse{}
	res.Stores = []types.Store{}
	isNPC, err := h.svc.IsNPC(r.Context(), req.NPCOID)
	if err != nil {
		return err
	}
	if isNPC {
		res.Stores = append(res.Stores, types.Store{
			AssetContainer: types.AssetContainer{OID: services.StoreOID, AssetMap: types.AssetMap{}, AssetPackages: []types.AssetPackage{}},
			ZoneOID:        req.ZoneOID,
			NPCOID:         req.NPCOID,
		})
	}
	return w.Write(res)
}

// GetAllStoreItems fills a store with free products: every item in an NPC store,
// the items of the opened subcategory in the World Shop.
func (h *Handler) GetAllStoreItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetAllStoreItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetAllStoreItemsResponse{}
	res.StoreItems = []types.StoreItem{}
	var theme *types.OID
	switch req.StoreOID {
	case services.StoreOID:
	case services.WorldShopOID:
		// The World Shop asks for the items of one subcategory.
		theme = &req.CategoryOID
	default:
		return w.Write(res)
	}
	items, err := h.svc.GetStoreItems(r.Context(), h.platform(r), theme)
	if err != nil {
		return err
	}
	res.StoreItems = items
	return w.Write(res)
}

// GetStoreThemes builds the World Shop categories and subcategories, the other stores have none.
func (h *Handler) GetStoreThemes(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetStoreThemesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	res := &messages.GetStoreThemesResponse{}
	res.StoreThemes = []types.StoreTheme{}
	if req.StoreOID == services.WorldShopOID {
		themes, err := h.svc.GetStoreThemes(r.Context(), h.platform(r))
		if err != nil {
			return err
		}
		res.StoreThemes = themes
	}
	return w.Write(res)
}

// EnterStore answers the opening of a store window, nothing is tracked.
func (h *Handler) EnterStore(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.EnterStoreRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.EnterStoreResponse{})
}

// ExitStore answers the closing of a store window, nothing is tracked.
func (h *Handler) ExitStore(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.ExitStoreRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.ExitStoreResponse{})
}

// PurchaseItems buys or sells store products, the store is free.
// The client reads every list of the response, so they are never null.
func (h *Handler) PurchaseItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.PurchaseItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	res := &messages.PurchaseItemsResponse{
		ReturnedPlayerItems:         []types.PlayerItem{},
		PurchasedPlayerItems:        []types.PlayerItem{},
		ReturnedPlayerBuildObjects:  []types.PlayerBuildObject{},
		PurchasedPlayerBuildObjects: []types.PlayerBuildObject{},
		PurchasedPlayerMazes:        []types.PlayerMaze{},
		PurchasedPlayerAvatars:      []types.PlayerAvatar{},
		ReturnedPlayerAvatars:       []types.PlayerAvatar{},
		ProcessedStoreItems:         []types.StoreItem{},
		AvatarAdoptionItems:         []types.PlayerItem{},
	}
	// A sale is a purchase of the buyback products with RequestedQuantity -1, the client never mixes it with a purchase
	// (WorldShopSVCS.SellProductItem, House.DoSellUpgrades).
	if len(req.Items) > 0 && req.Items[0].RequestedQuantity < 0 {
		items, err := h.svc.SellStoreItems(r.Context(), session.platform, session.player.OID, req.Items)
		if errors.Is(err, services.ErrStoreItemNotFound) {
			res.PurchaseResponse = purchaseInvalid
			return w.Write(res)
		}
		if err != nil {
			return err
		}
		res.ReturnedPlayerItems = items
		return w.Write(res)
	}
	items, mazes, err := h.svc.BuyStoreItems(r.Context(), session.platform, session.player.OID, req.Items)
	if errors.Is(err, services.ErrStoreItemNotFound) {
		res.PurchaseResponse = purchaseInvalid
		return w.Write(res)
	}
	if err != nil {
		return err
	}
	res.PurchasedPlayerItems = items
	res.PurchasedPlayerMazes = mazes
	return w.Write(res)
}

// GetStoreItemsForItems lets the player sell an item: it gives the sell card and the price of a house with its upgrades.
func (h *Handler) GetStoreItemsForItems(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetStoreItemsForItemsRequest{}
	if err := r.Read(req); err != nil {
		return err
	}
	session, err := h.session(r)
	if err != nil {
		return err
	}
	storeItem, err := h.svc.GetBuybackItem(r.Context(), session.platform, session.player.OID, req.ItemOID)
	if err != nil {
		return err
	}
	return w.Write(&messages.GetStoreItemsForItemsResponse{StoreItem: storeItem})
}

// purchaseInvalid is the INVALID_PURCHASES purchase response, NPCStoreSVCS shows "not enabled for purchase".
const purchaseInvalid = 163

// GetCurrencies lets the store windows show prices, they need the Coinz name of the price currency.
func (h *Handler) GetCurrencies(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetCurrenciesRequest{}); err != nil {
		return err
	}
	currencies, err := h.svc.GetCurrencies(r.Context(), h.platform(r))
	if err != nil {
		return err
	}
	return w.Write(&messages.GetCurrenciesResponse{Currencies: currencies})
}

// GetPlayerStats answers with no statistics, the player has no money and the store is free.
func (h *Handler) GetPlayerStats(w gsf.ResponseWriter, r *gsf.Request) error {
	if err := r.Read(&messages.GetPlayerStatsRequest{}); err != nil {
		return err
	}
	return w.Write(&messages.GetPlayerStatsResponse{PlayerStats: []types.PlayerStats{}})
}
