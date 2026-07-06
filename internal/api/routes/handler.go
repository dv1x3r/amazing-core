package routes

import (
	"github.com/dv1x3r/amazing-core/internal/lib/slogbus"
	"github.com/dv1x3r/amazing-core/internal/lib/webauth"
	"github.com/dv1x3r/amazing-core/internal/services"
)

type Handler struct {
	svc    services.Services
	auth   *webauth.Authenticator
	logBus *slogbus.Bus
}

func NewHandler(svc services.Services, auth *webauth.Authenticator, logBus *slogbus.Bus) *Handler {
	return &Handler{
		svc:    svc,
		auth:   auth,
		logBus: logBus,
	}
}

func (h *Handler) Authenticator() *webauth.Authenticator {
	return h.auth
}
