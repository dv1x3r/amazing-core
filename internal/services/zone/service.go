package zone

import (
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/services/asset"
)

var (
	ErrZoneExists = errors.New("zone with the same container already exists")
)

type Service struct {
	logger *slog.Logger
	store  db.Store
	assets *asset.Service
}

func NewService(logger *slog.Logger, store db.Store, assets *asset.Service) *Service {
	return &Service{
		logger: logger,
		store:  store,
		assets: assets,
	}
}
