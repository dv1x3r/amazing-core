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

// syncKeepAlive is shorter than the 60 seconds the client waits for data on the SYNC connection.
const syncKeepAlive = 30 * time.Second

// Server is the game server, it owns the sessions and the world, the handlers get them from it.
type Server struct {
	logger   *slog.Logger
	server   *gsf.Server
	sessions *SessionManager
	world    *worldsync.World
	stop     chan struct{}
}

func NewServer(logger *slog.Logger, svc *services.Service) *Server {
	s := &Server{
		logger:   logger,
		sessions: NewSessionManager(),
		world:    worldsync.NewWorld(),
		stop:     make(chan struct{}),
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
	handleUser(usermessagetype.PLACE_MAZE_ITEM, handler.PlaceMazeItem)
	handleUser(usermessagetype.REMOVE_MAZE_ITEM, handler.RemoveMazeItem)
	handleUser(usermessagetype.REMOVE_AND_DETACH_MAZE_ITEM, handler.RemoveAndDetachMazeItem)
	handleUser(usermessagetype.DETACH_ITEMS, handler.DetachItems)
	handleUser(usermessagetype.ATTACH_SINGLE_ITEM, handler.AttachSingleItem)
	handleUser(usermessagetype.ENTER_BUILDING, handler.EnterBuilding)
	handleUser(usermessagetype.GET_ALL_BUILDINGS, handler.GetAllBuildings)

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
	handleUser(usermessagetype.GET_FRIENDSHIP_REQUEST_COUNTS, handler.GetFriendshipRequestCounts)
	handleUser(usermessagetype.LIST_LIMITS, handler.ListLimits)
	handleUser(usermessagetype.GET_HOME_INVITATIONS, handler.GetHomeInvitations)
	handleUser(usermessagetype.GET_NOTIFICATION_BY_PLAYER_ID, handler.GetNotificationByPlayerID)
	handleUser(usermessagetype.GET_CRISP_ACTIONS, handler.GetCrispActions)

	// ── NPC and Quests ───────────────────────────────────────────────────────────
	handleUser(usermessagetype.GET_PLAYER_NPCS, handler.GetPlayerNPCs)
	handleUser(usermessagetype.GET_RULE_COUNT_LIST, handler.GetRuleCountList)
	handleUser(usermessagetype.GET_RULE_INFO, handler.GetRuleInfo)
	handleUser(usermessagetype.GET_RULE_INSTANCE_INFO, handler.GetRuleInstanceInfo)
	handleUser(usermessagetype.GET_NPC_RELATIONSHIPS, handler.GetNPCRelationships)
	handleUser(usermessagetype.GET_NPC_RELATIONSHIP_LEVELS, handler.GetNPCRelationshipLevels)
	handleUser(usermessagetype.GET_NPCS_WITH_QUEST_OFFER, handler.GetNPCsWithQuestOffer)
	handleUser(usermessagetype.GET_PLAYER_QUESTS, handler.GetPlayerQuests)
	handleUser(usermessagetype.GET_PLAYER_MISSIONS, handler.GetPlayerMissions)
	handleUser(usermessagetype.GET_CMS_MISSIONS, handler.GetCMSMissions)
	handleUser(usermessagetype.QUEST_EVENT_IN_PROGRESS, handler.QuestEventInProgress)
	handleUser(usermessagetype.ADD_EVENT, handler.AddEvent)
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
	handleUser(usermessagetype.LIST_STORES, handler.ListStores)
	handleUser(usermessagetype.GET_ALL_STORE_ITEMS, handler.GetAllStoreItems)
	handleUser(usermessagetype.GET_STORE_THEMES, handler.GetStoreThemes)
	handleUser(usermessagetype.ENTER_STORE, handler.EnterStore)
	handleUser(usermessagetype.EXIT_STORE, handler.ExitStore)
	handleUser(usermessagetype.PURCHASE_ITEMS, handler.PurchaseItems)
	handleUser(usermessagetype.GET_STORE_ITEMS_FOR_ITEMS, handler.GetStoreItemsForItems)
	handleUser(usermessagetype.GET_CURRENCIES, handler.GetCurrencies)
	handleUser(usermessagetype.GET_PLAYER_STATS, handler.GetPlayerStats)
	handleUser(usermessagetype.GET_REQUIRED_EXPERIENCE, handler.GetRequiredExperience)
	handleUser(usermessagetype.GET_STATS_TYPE, handler.GetStatsType)

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
	go s.keepSyncAlive()
	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			s.logger.Error(err.Error())
		}
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	s.logger.Info("shutting down the game server")
	close(s.stop)
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Server) keepSyncAlive() {
	ticker := time.NewTicker(syncKeepAlive)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.world.KeepAlive()
		case <-s.stop:
			return
		}
	}
}
