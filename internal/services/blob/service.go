package blob

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
)

var (
	ErrFileNotFound = errors.New("file is not found")
	ErrFileExists   = errors.New("file with the same name already exists")
)

type Service struct {
	logger      *slog.Logger
	store       db.Store
	deliveryURL string
}

func NewService(logger *slog.Logger, store db.Store, deliveryURL string) *Service {
	return &Service{
		logger:      logger,
		store:       store,
		deliveryURL: deliveryURL,
	}
}

func (s *Service) ImportFromFolder(ctx context.Context) (*ImportResult, error) {
	return ImportFromFolder(ctx, s.logger, s.store.DB(), "cache")
}

func (s *Service) ExportToFolder(ctx context.Context) (*ExportResult, error) {
	return ExportToFolder(ctx, s.logger, s.store.DB(), "cache", true)
}
