package avatar

import (
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/services/asset"
)

var (
	ErrAvatarExists = errors.New("avatar with the same name or container already exists")
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
