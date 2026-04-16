package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/dv1x3r/amazing-core/internal/game/middleware"
	"github.com/dv1x3r/amazing-core/internal/network/bitprotocol"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
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
		middleware.Logger(logger),
		middleware.Recover(),
	)

	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_CLIENT_VERSION_INFO), handler.GetClientVersionInfo)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_PUBLIC_ITEM_CATEGORIES), handler.GetPublicItemCategories)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_PUBLIC_ITEMS_BY_OIDS), handler.GetPublicItemsByOIDs)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_RANDOM_NAMES), handler.GetRandomNames)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.VALIDATE_NAME), handler.ValidateName)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.SELECT_PLAYER_NAME), handler.SelectPlayerName)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.CHECK_USERNAME), handler.CheckUsername)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.REGISTER_PLAYER), handler.RegisterPlayer)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.LOGIN), handler.Login)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_TIERS), handler.GetTiers)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_SITE_FRAME), handler.GetSiteFrame)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_CMS_ITEMCATEGORIES), handler.GetCMSItemCategories)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_OUTFIT_ITEMS), handler.GetOutfitItems)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_AVATARS), handler.GetAvatars)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_OUTFITS), handler.GetOutfits)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_ZONES), handler.GetZones)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.INIT_LOCATION), handler.InitLocation)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_MAZE_ITEMS), handler.GetMazeItems)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_CHAT_CHANNEL_TYPES), handler.GetChatChannelTypes)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_ANNOUNCEMENTS), handler.GetAnnouncements)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.ENTER_BUILDING), handler.EnterBuilding)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_ONLINE_STATUSES), handler.GetOnlineStatuses)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_PLAYER_NPCS), handler.GetPlayerNPCs)
	router.HandleFunc(int32(serviceclass.SYNC_SERVER), int32(syncmessagetype.LOGIN), handler.SyncLogin)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.LOGOUT), handler.Logout)

	server := &gsf.Server{
		Router: router,
		Codec:  bitprotocol.NewBitCodec(),
		Hooks: gsf.ServerHooks{
			OnConnect: func(remoteIP string) {
				logger.Info(fmt.Sprintf("tcp %s connected", remoteIP))
			},
			OnDisconnect: func(remoteIP string, reason string) {
				logger.Info(fmt.Sprintf("tcp %s disconnected: %s", remoteIP, reason))
			},
			OnUnhandled: func(remoteIP string, header *gsf.Header, data []byte) {
				logger.Warn(fmt.Sprintf("gsf %s unhandled: %+v", remoteIP, header),
					slog.String("service_class", header.ServiceClassText()),
					slog.String("message_type", header.MessageTypeText()),
					slog.Any("hex", fmt.Sprintf("%x", data)),
				)
			},
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
