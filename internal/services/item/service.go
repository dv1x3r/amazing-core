package item

import (
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
)

var (
	ErrItemNotFound             = errors.New("item not found")
	ErrItemExists               = errors.New("item with the same name or container already exists")
	ErrCategoryExists           = errors.New("category with the same name already exists")
	ErrCategoryCyclicDependency = errors.New("circular dependency detected (A → B → A)")
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
