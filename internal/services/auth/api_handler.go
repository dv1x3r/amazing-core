package auth

import (
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
)

type APIHandler struct {
	service *Service
}

func NewAPIHandler(service *Service) *APIHandler {
	return &APIHandler{service: service}
}

func (h *APIHandler) PostLogin(w http.ResponseWriter, r *http.Request) error {
	form, err := w2.ParseFormSaveRequest[AdminLoginForm](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if _, err := h.service.AuthenticateSession(w, r, form.Record); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *APIHandler) PostLogout(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeauthenticateSession(w, r); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}
