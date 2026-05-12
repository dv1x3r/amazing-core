package player

import (
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/services/avatar"
	"github.com/dv1x3r/amazing-core/internal/services/item"
)

var (
	ErrPlayerNotFound     = errors.New("player not found")
	ErrPlayerExists       = errors.New("player with the same oid already exists")
	ErrPlayerAvatarExists = errors.New("player avatar with the same name or zing already exists")
	ErrPlayerOutfitExists = errors.New("player outfit with the same avatar and outfit number already exists")
	ErrPlayerItemExists   = errors.New("player item with the same oid already exists")
	ErrPlayerItemAttached = errors.New("player item is already attached")
)

type Service struct {
	logger  *slog.Logger
	store   db.Store
	avatars *avatar.Service
	items   *item.Service
}

func NewService(logger *slog.Logger, store db.Store, avatars *avatar.Service, items *item.Service) *Service {
	return &Service{
		logger:  logger,
		store:   store,
		avatars: avatars,
		items:   items,
	}
}
