package services

import (
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/lib/db"

	"github.com/dv1x3r/amazing-core/internal/services/asset"
	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/amazing-core/internal/services/avatar"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
	"github.com/dv1x3r/amazing-core/internal/services/dummy"
	"github.com/dv1x3r/amazing-core/internal/services/item"
	"github.com/dv1x3r/amazing-core/internal/services/player"
	"github.com/dv1x3r/amazing-core/internal/services/randname"
	"github.com/dv1x3r/amazing-core/internal/services/siteframe"
)

type Services struct {
	Auth      *auth.Service
	Asset     *asset.Service
	Avatar    *avatar.Service
	Blob      *blob.Service
	Dummy     *dummy.Service
	Item      *item.Service
	Player    *player.Service
	RandName  *randname.Service
	SiteFrame *siteframe.Service
}

func New(logger *slog.Logger, coreStore db.Store, blobStore db.Store, cfg config.Config) Services {
	assets := asset.NewService(logger, coreStore, cfg.Settings.AssetDeliveryURL)
	avatars := avatar.NewService(logger, coreStore, assets)
	items := item.NewService(logger, coreStore, assets)
	return Services{
		Asset:     assets,
		Auth:      auth.NewService(cfg.Secure.Session.Secure, cfg.Secure.Session.Key, cfg.Secure.Auth.Username, cfg.Secure.Auth.Password),
		Avatar:    avatars,
		Blob:      blob.NewService(logger, blobStore, cfg.Settings.AssetDeliveryURL),
		Dummy:     dummy.NewService(logger, coreStore),
		Item:      items,
		Player:    player.NewService(logger, coreStore, avatars, items),
		RandName:  randname.NewService(logger, coreStore),
		SiteFrame: siteframe.NewService(logger, coreStore, assets),
	}
}
