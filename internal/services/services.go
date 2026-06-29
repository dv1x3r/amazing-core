package services

import (
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/python"

	"github.com/dv1x3r/amazing-core/internal/services/asset"
	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/amazing-core/internal/services/avatar"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
	"github.com/dv1x3r/amazing-core/internal/services/item"
	"github.com/dv1x3r/amazing-core/internal/services/player"
	"github.com/dv1x3r/amazing-core/internal/services/randname"
	"github.com/dv1x3r/amazing-core/internal/services/siteframe"
	"github.com/dv1x3r/amazing-core/internal/services/zone"
)

type Services struct {
	Auth      *auth.Service
	Asset     *asset.Service
	Avatar    *avatar.Service
	Blob      *blob.Service
	Item      *item.Service
	Player    *player.Service
	RandName  *randname.Service
	SiteFrame *siteframe.Service
	Zone      *zone.Service
}

func New(logger *slog.Logger, store db.Store, cfg config.Config) Services {
	assets := asset.NewService(logger, store, cfg.Settings.AssetDeliveryURL)
	avatars := avatar.NewService(logger, store, assets)
	items := item.NewService(logger, store, assets)
	pythonRunner := python.NewRunner(logger, cfg.Storage.PythonVenv)
	return Services{
		Asset:     assets,
		Auth:      auth.NewService(cfg.Secure.Session.Secure, cfg.Secure.Session.Key, cfg.Secure.Auth.Username, cfg.Secure.Auth.Password),
		Avatar:    avatars,
		Blob:      blob.NewService(logger, store, pythonRunner),
		Item:      items,
		Player:    player.NewService(logger, store, avatars, items),
		RandName:  randname.NewService(logger, store),
		SiteFrame: siteframe.NewService(logger, store, assets),
		Zone:      zone.NewService(logger, store, assets),
	}
}
