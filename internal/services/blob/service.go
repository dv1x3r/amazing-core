package blob

import (
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/python"
)

var (
	ErrFileNotFound = errors.New("file is not found")
	ErrFileExists   = errors.New("file with the same name already exists")
)

type Service struct {
	logger *slog.Logger
	store  db.Store
	python *python.Runner
}

func NewService(logger *slog.Logger, store db.Store, pythonRunner *python.Runner) *Service {
	return &Service{
		logger: logger,
		store:  store,
		python: pythonRunner,
	}
}
