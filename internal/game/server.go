package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/dv1x3r/amazing-core/internal/game/middleware"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/bitprotocol"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/serviceclass"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/syncmessagetype"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/usermessagetype"
)

type Server struct {
	logger *slog.Logger
	server *gsf.Server
}

func NewServer(
	logger *slog.Logger,
	handler *Handler,
) *Server {
	router := gsf.NewRouter()

	router.Use(
		middleware.Recover(),
	)

	// ── General ──────────────────────────────────────────────────────────────────
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_CLIENT_VERSION_INFO), handler.GetClientVersionInfo)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_SITE_FRAME), handler.GetSiteFrame)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.LOGIN), handler.Login)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.LOGOUT), handler.Logout)

	// ── Registration ─────────────────────────────────────────────────────────────
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_RANDOM_NAMES), handler.GetRandomNames)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.VALIDATE_NAME), handler.ValidateName)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.SELECT_PLAYER_NAME), handler.SelectPlayerName)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.CHECK_USERNAME), handler.CheckUsername)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.REGISTER_PLAYER), handler.RegisterPlayer)

	// ── Inventory ────────────────────────────────────────────────────────────────
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_PUBLIC_ITEM_CATEGORIES), handler.GetPublicItemCategories)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_CMS_ITEMCATEGORIES), handler.GetCMSItemCategories)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_BUILD_OBJECTS), handler.GetBuildObjects)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_INVENTORY_OBJECTS), handler.GetInventoryObjects)

	// ── Avatars ──────────────────────────────────────────────────────────────────
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_AVATARS), handler.GetAvatars)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_AVATAR_ITEMS), handler.GetAvatarItems)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.UPDATE_PLAYER_ACTIVE_AVATAR), handler.UpdatePlayerActiveAvatar)

	// ──  Outfits ─────────────────────────────────────────────────────────────────
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_OUTFITS), handler.GetOutfits)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.ADD_OUTFIT), handler.AddOutfit)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.SET_CURRENT_OUTFIT), handler.SetCurrentOutfit)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_OUTFIT_ITEMS), handler.GetOutfitItems)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.ADD_OUTFIT_ITEMS), handler.AddOutfitItems)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.REMOVE_OUTFIT_ITEMS), handler.RemoveOutfitItems)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.REPLACE_OUTFIT_ITEMS), handler.ReplaceOutfitItems)

	// ── Login ────────────────────────────────────────────────────────────────────
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_TIERS), handler.GetTiers)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_PUBLIC_ITEMS_BY_OIDS), handler.GetPublicItemsByOIDs)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_ZONES), handler.GetZones)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.INIT_LOCATION), handler.InitLocation)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.HEARTBEAT), handler.Heartbeat)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_OTHER_PLAYER_DETAILS), handler.GetOtherPlayerDetails)
	router.HandleFunc(int32(serviceclass.SYNC_SERVER), int32(syncmessagetype.LOGIN), handler.SyncLogin)
	router.HandleFunc(int32(serviceclass.SYNC_SERVER), int32(syncmessagetype.MOVE_PLAYER), handler.MovePlayer)
	router.HandleFunc(int32(serviceclass.SYNC_SERVER), int32(syncmessagetype.HEARTBEAT_NOTIFY), handler.HeartbeatNotify)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_MAZE_ITEMS), handler.GetMazeItems)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_CHAT_CHANNEL_TYPES), handler.GetChatChannelTypes)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.SEND_MESSAGE), handler.SendMessage)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_ANNOUNCEMENTS), handler.GetAnnouncements)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.ENTER_BUILDING), handler.EnterBuilding)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_ONLINE_STATUSES), handler.GetOnlineStatuses)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_PLAYER_NPCS), handler.GetPlayerNPCs)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_STORE_ITEMS), handler.GetStoreItems)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_PLAYER_QUESTS), handler.GetPlayerQuests)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_QUEST_FROM_PARENT), handler.GetQuestFromParent)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.CREATE_QUEST), handler.CreateQuest)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.ACCEPT_QUEST), handler.AcceptQuest)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.START_QUEST), handler.StartQuest)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.COMPLETE_QUEST), handler.CompleteQuest)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_PLAYER_QUESTS_BY_QUEST_IDS), handler.GetPlayerQuestsByOIDs)

	gsfLogger := newGSFLogger(logger)

	server := &gsf.Server{
		Router: router,
		Codec:  bitprotocol.NewBitCodec(),
		Hooks: gsf.ServerHooks{
			OnConnect: func(session *gsf.Session) {
				logger.Info(fmt.Sprintf("tcp %s connected", session.RemoteIP()))
			},
			OnDisconnect: func(session *gsf.Session, reason string) {
				if err := handler.syncHub.Leave(session); err != nil {
					logger.Error(err.Error())
				}
				logger.Info(fmt.Sprintf("tcp %s disconnected: %s", session.RemoteIP(), reason))
			},
			OnUnhandled: gsfLogger.OnUnhandled,
			OnRequest:   gsfLogger.OnRequest,
			OnNotify:    gsfLogger.OnNotify,
		},
	}

	return &Server{
		logger: logger,
		server: server,
	}
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
