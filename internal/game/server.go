package game

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/game/dummy"
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/gsf/bitprotocol"
	"github.com/dv1x3r/amazing-core/internal/game/middleware"
	"github.com/dv1x3r/amazing-core/internal/game/types/serviceclass"
	"github.com/dv1x3r/amazing-core/internal/game/types/syncmessagetype"
	"github.com/dv1x3r/amazing-core/internal/game/types/usermessagetype"
	"github.com/dv1x3r/amazing-core/internal/lib/logger"
	"github.com/dv1x3r/amazing-core/internal/services/randomnames"
)

type Server struct {
	server *gsf.Server
}

func NewServer(
	randomNamesService *randomnames.Service,
) *Server {
	router := gsf.NewRouter()

	router.Use(
		middleware.Logger(logger.Get()),
		middleware.Recover(),
	)

	randomNamesHandler := randomnames.NewGSFHandler(randomNamesService)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_RANDOM_NAMES), randomNamesHandler.GetRandomNames)

	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_CLIENT_VERSION_INFO), dummy.GetClientVersionInfo)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_PUBLIC_ITEM_CATEGORIES), dummy.GetPublicItemCategories)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_PUBLIC_ITEMS_BY_OIDS), dummy.GetPublicItemsByOIDs)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.VALIDATE_NAME), dummy.ValidateName)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.SELECT_PLAYER_NAME), dummy.SelectPlayerName)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.CHECK_USERNAME), dummy.CheckUsername)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.REGISTER_PLAYER), dummy.RegisterPlayer)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.LOGIN), dummy.Login)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_TIERS), dummy.GetTiers)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_SITE_FRAME), dummy.GetSiteFrame)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_OUTFIT_ITEMS), dummy.GetOutfitItems)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_AVATARS), dummy.GetAvatars)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_OUTFITS), dummy.GetOutfits)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_ZONES), dummy.GetZones)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.INIT_LOCATION), dummy.InitLocation)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_MAZE_ITEMS), dummy.GetMazeItems)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_CHAT_CHANNEL_TYPES), dummy.GetChatChannelTypes)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_ANNOUNCEMENTS), dummy.GetAnnouncements)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.ENTER_BUILDING), dummy.EnterBuilding)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_ONLINE_STATUSES), dummy.GetOnlineStatuses)
	router.HandleFunc(int32(serviceclass.USER_SERVER), int32(usermessagetype.GET_PLAYER_NPCS), dummy.GetPlayerNPCs)
	router.HandleFunc(int32(serviceclass.SYNC_SERVER), int32(syncmessagetype.LOGIN), dummy.SyncLogin)

	server := &gsf.Server{
		Router: router,
		Codec:  bitprotocol.NewBitCodec(),
	}

	return &Server{server: server}
}

func (s *Server) Start(address string) error {
	s.server.Addr = address
	logger.Get().Info("starting the game server on " + address)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	logger.Get().Info("shutting down the game server")
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Get().Error("[gsf]" + err.Error())
	}
}
