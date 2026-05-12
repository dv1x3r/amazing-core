package randname

import (
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
)

var (
	ErrNameNotFound = errors.New("name not found")
	ErrNameExists   = errors.New("name with the same type and name already exists")
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
