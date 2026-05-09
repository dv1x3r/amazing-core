package api

import (
	"fmt"
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/w2go/w2"
)

func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) error {
	username, ok := h.svc.Auth.GetSessionUsername(w, r)
	if ok {
		if err := h.svc.Auth.RefreshSession(w, r); err != nil {
			return err
		}
	}
	cfg := config.Get()
	data := map[string]any{
		"username": username,
		"version":  cfg.Version,
		"explorer": cfg.Storage.Explorer,
	}
	return tmpl.ExecuteTemplate(w, "dashboard.gohtml", data)
}

func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) error {
	form, err := w2.ParseSaveFormRequest[auth.AdminLoginForm](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if valid, err := h.svc.Auth.AuthenticateSession(w, r, form.Record); err != nil {
		return err
	} else if !valid {
		return wrap.WithHTTPStatus(fmt.Errorf("Invalid username or password"), http.StatusUnauthorized)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostLogout(w http.ResponseWriter, r *http.Request) error {
	if err := h.svc.Auth.DeauthenticateSession(w, r); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}
