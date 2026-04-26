package asset

import (
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
)

var (
	ErrContainerExists         = errors.New("container with this GSF OID already exists")
	ErrContainerInUse          = errors.New("container is still in use and cannot be removed")
	ErrContainerAssetExists    = errors.New("container already contains this primary asset")
	ErrContainerPackageExists  = errors.New("container already contains this package")
	ErrPackageCyclicDependency = errors.New("circular dependency detected (A → B → A)")
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

func (s *Service) DeliveryURL() string {
	return s.deliveryURL
}
