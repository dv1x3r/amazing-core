package services

import (
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/python"
)

type Service struct {
	logger      *slog.Logger
	store       db.Store
	deliveryURL string
	python      *python.Runner
}

func NewService(logger *slog.Logger, store db.Store, cfg config.Config, pythonRunner *python.Runner) *Service {
	return &Service{
		logger:      logger,
		store:       store,
		deliveryURL: cfg.Settings.AssetDeliveryURL,
		python:      pythonRunner,
	}
}
