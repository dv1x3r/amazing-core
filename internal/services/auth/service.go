package auth

import (
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
)

type Service struct {
	logger *slog.Logger
	store  db.Store
}

func NewService(logger *slog.Logger, store db.Store) *Service {
	return &Service{
		logger: logger,
		store:  store,
	}
}
