package player

import (
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/services/asset"
	"github.com/dv1x3r/amazing-core/internal/services/avatar"
)

var (
	ErrPlayerNotFound     = errors.New("player not found")
	ErrPlayerExists       = errors.New("player with the same oid already exists")
	ErrPlayerAvatarExists = errors.New("player avatar with the same name or zing already exists")
	ErrPlayerOutfitExists = errors.New("player outfit with the same avatar and outfit number already exists")
)

type Service struct {
	logger  *slog.Logger
	store   db.Store
	assets  *asset.Service
	avatars *avatar.Service
}

func NewService(logger *slog.Logger, store db.Store, assets *asset.Service, avatars *avatar.Service) *Service {
	return &Service{
		logger:  logger,
		store:   store,
		assets:  assets,
		avatars: avatars,
	}
}
