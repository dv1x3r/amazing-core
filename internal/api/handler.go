package api

import (
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/lib/slogbus"
	"github.com/dv1x3r/amazing-core/internal/lib/webauth"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/services"
	"github.com/dv1x3r/w2go/w2"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

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

func errorHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		// client cancelled, do not write anything
		if r.Context().Err() != nil {
			if lw, ok := w.(interface{ SetError(error) }); ok {
				lw.SetError(r.Context().Err())
			}
			return
		}

		if lw, ok := w.(interface{ SetError(error) }); ok {
			lw.SetError(err)
		}

		status := wrap.HTTPStatus(err)
		message := err.Error()

		// do not expose details of 500 errors if not in debug
		if status >= 500 && config.Get().Logger.Level != "debug" {
			message = http.StatusText(status)
		}

		res := w2.NewErrorResponse(message)
		res.Write(w, status)
	}
}
