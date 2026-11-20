package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/middleware"
	"github.com/dv1x3r/amazing-core/internal/game/services"
	"github.com/dv1x3r/amazing-core/internal/game/worldsync"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/bitprotocol"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/serviceclass"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/syncmessagetype"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/usermessagetype"
)

// idleTimeout closes a dead connection.
// After LOGIN the client sends HEARTBEAT and HEARTBEAT_NOTIFY every 60 seconds.
const idleTimeout = 10 * time.Minute

// Server is the game server, it owns the sessions and the world, the handlers get them from it.
type Server struct {
	logger   *slog.Logger
	server   *gsf.Server
	sessions *SessionManager
	world    *worldsync.World
}

func NewServer(logger *slog.Logger, svc *services.Service) *Server {
	s := &Server{
		logger:   logger,
		sessions: NewSessionManager(),
		world:    worldsync.NewWorld(),
	}
	handler := NewHandler(svc, s.sessions, s.world)

	router := gsf.NewRouter()

	router.Use(
		middleware.Recover(),
	)

	handleUser := func(msgType usermessagetype.UserMessageType, fn gsf.HandlerFunc) {
		router.HandleFunc(int32(serviceclass.USER_SERVER), int32(msgType), fn)
	}
	handleSync := func(msgType syncmessagetype.SyncMessageType, fn gsf.HandlerFunc) {
		router.HandleFunc(int32(serviceclass.SYNC_SERVER), int32(msgType), fn)
	}

	// ── General ──────────────────────────────────────────────────────────────────
	handleUser(usermessagetype.GET_CLIENT_VERSION_INFO, handler.GetClientVersionInfo)
	handleUser(usermessagetype.GET_SITE_FRAME, handler.GetSiteFrame)
	handleUser(usermessagetype.LOGIN, handler.Login)
	handleUser(usermessagetype.LOGOUT, handler.Logout)
	handleUser(usermessagetype.RELOGIN, handler.Relogin)
	handleUser(usermessagetype.HEARTBEAT, handler.Heartbeat)
	handleUser(usermessagetype.GET_TIERS, handler.GetTiers)

	// ── Registration ─────────────────────────────────────────────────────────────
	handleUser(usermessagetype.GET_RANDOM_NAMES, handler.GetRandomNames)
	handleUser(usermessagetype.VALIDATE_NAME, handler.ValidateName)
	handleUser(usermessagetype.SELECT_PLAYER_NAME, handler.SelectPlayerName)
	handleUser(usermessagetype.CHECK_USERNAME, handler.CheckUsername)
	handleUser(usermessagetype.REGISTER_PLAYER, handler.RegisterPlayer)

	// ── Inventory ────────────────────────────────────────────────────────────────
	handleUser(usermessagetype.GET_PUBLIC_ITEM_CATEGORIES, handler.GetPublicItemCategories)
	handleUser(usermessagetype.GET_CMS_ITEMCATEGORIES, handler.GetCMSItemCategories)
	handleUser(usermessagetype.GET_BUILD_OBJECTS, handler.GetBuildObjects)
	handleUser(usermessagetype.GET_INVENTORY_OBJECTS, handler.GetInventoryObjects)
	handleUser(usermessagetype.GET_PUBLIC_ITEMS_BY_OIDS, handler.GetPublicItemsByOIDs)

	// ── Avatars ──────────────────────────────────────────────────────────────────
	handleUser(usermessagetype.GET_AVATARS, handler.GetAvatars)
	handleUser(usermessagetype.GET_AVATAR_ITEMS, handler.GetAvatarItems)
	handleUser(usermessagetype.UPDATE_PLAYER_ACTIVE_AVATAR, handler.UpdatePlayerActiveAvatar)

	// ── Outfits ──────────────────────────────────────────────────────────────────
	handleUser(usermessagetype.GET_OUTFITS, handler.GetOutfits)
	handleUser(usermessagetype.ADD_OUTFIT, handler.AddOutfit)
	handleUser(usermessagetype.SET_CURRENT_OUTFIT, handler.SetCurrentOutfit)
	handleUser(usermessagetype.GET_OUTFIT_ITEMS, handler.GetOutfitItems)
	handleUser(usermessagetype.ADD_OUTFIT_ITEMS, handler.AddOutfitItems)
	handleUser(usermessagetype.REMOVE_OUTFIT_ITEMS, handler.RemoveOutfitItems)
	handleUser(usermessagetype.REPLACE_OUTFIT_ITEMS, handler.ReplaceOutfitItems)

	// ── Location ─────────────────────────────────────────────────────────────────
	handleUser(usermessagetype.GET_ZONES, handler.GetZones)
	handleUser(usermessagetype.INIT_LOCATION, handler.InitLocation)
	handleUser(usermessagetype.GET_MAZE_ITEMS, handler.GetMazeItems)
	handleUser(usermessagetype.ENTER_BUILDING, handler.EnterBuilding)

	// ── Sync ─────────────────────────────────────────────────────────────────────
	handleSync(syncmessagetype.LOGIN, handler.SyncLogin)
	handleSync(syncmessagetype.RELOGIN, handler.SyncLogin)
	handleSync(syncmessagetype.MOVE_PLAYER, handler.MovePlayer)
	handleSync(syncmessagetype.ACTION, handler.SyncAction)
	handleSync(syncmessagetype.HEARTBEAT_NOTIFY, handler.HeartbeatNotify)

	// ── Social ───────────────────────────────────────────────────────────────────
	handleUser(usermessagetype.GET_OTHER_PLAYER_DETAILS, handler.GetOtherPlayerDetails)
	handleUser(usermessagetype.GET_CHAT_CHANNEL_TYPES, handler.GetChatChannelTypes)
	handleUser(usermessagetype.SEND_MESSAGE, handler.SendMessage)
	handleUser(usermessagetype.GET_ANNOUNCEMENTS, handler.GetAnnouncements)
	handleUser(usermessagetype.GET_ONLINE_STATUSES, handler.GetOnlineStatuses)

	// ── NPC and Quests ───────────────────────────────────────────────────────────
	handleUser(usermessagetype.GET_PLAYER_NPCS, handler.GetPlayerNPCs)
	handleUser(usermessagetype.GET_RULE_COUNT_LIST, handler.GetRuleCountList)
	handleUser(usermessagetype.GET_RULE_INFO, handler.GetRuleInfo)
	handleUser(usermessagetype.GET_RULE_INSTANCE_INFO, handler.GetRuleInstanceInfo)
	handleUser(usermessagetype.GET_PLAYER_QUESTS, handler.GetPlayerQuests)
	handleUser(usermessagetype.GET_PLAYER_QUESTS_BY_QUEST_IDS, handler.GetPlayerQuestsByOIDs)
	handleUser(usermessagetype.GET_QUEST_FROM_PARENT, handler.GetQuestFromParent)
	handleUser(usermessagetype.CREATE_QUEST, handler.CreateQuest)
	handleUser(usermessagetype.ACCEPT_QUEST, handler.AcceptQuest)
	handleUser(usermessagetype.START_QUEST, handler.StartQuest)
	handleUser(usermessagetype.ADD_QUEST_ITEM, handler.AddQuestItem)
	handleUser(usermessagetype.COMPLETE_QUEST, handler.CompleteQuest)
	handleUser(usermessagetype.ABANDON_QUEST, handler.AbandonQuest)
	handleUser(usermessagetype.GET_CMS_NOTIFICATIONS, handler.GetCMSNotifications)
	handleUser(usermessagetype.GET_ASSETS_BY_OIDS, handler.GetAssetsByOIDs)

	// ── Store ────────────────────────────────────────────────────────────────────
	handleUser(usermessagetype.GET_STORE_ITEMS, handler.GetStoreItems)

	gsfLogger := newGSFLogger(logger, s.sessions)

	s.server = &gsf.Server{
		Router:      router,
		Codec:       bitprotocol.NewBitCodec(),
		IdleTimeout: idleTimeout,
		Hooks: gsf.ServerHooks{
			OnConnect: func(conn *gsf.Conn) {
				logger.Info(fmt.Sprintf("tcp %s connected", conn.RemoteIP()))
			},
			OnDisconnect: func(conn *gsf.Conn, reason string) {
				s.world.Disconnect(conn)
				s.sessions.Detach(conn)
				logger.Info(fmt.Sprintf("tcp %s disconnected: %s", conn.RemoteIP(), reason))
			},
			OnUnhandled: gsfLogger.OnUnhandled,
			OnRequest:   gsfLogger.OnRequest,
			OnNotify:    gsfLogger.OnNotify,
		},
	}

	return s
}

func (s *Server) ListenAndServe(address string) {
	s.server.Addr = address
	s.logger.Info("starting the game server on " + address)
	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			s.logger.Error(err.Error())
		}
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	s.logger.Info("shutting down the game server")
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(err.Error())
	}
}
