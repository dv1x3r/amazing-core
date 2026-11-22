package routes

import (
	"errors"
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/api/services"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
)

func (h *Handler) GetStoreThemeGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.GetStoreThemeGrid(r.Context(), req)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) GetStoreThemeDropdown(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.GetStoreThemesDropdown(r.Context(), req)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostStoreThemeForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[services.StoreTheme](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if req.Record.ID == 0 {
		req.Record.ID, err = h.svc.CreateStoreTheme(r.Context(), req)
	} else {
		err = h.svc.UpdateStoreTheme(r.Context(), req)
	}
	if errors.Is(err, services.ErrStoreThemeExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if errors.Is(err, services.ErrStoreThemeCyclicDependency) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return err
	}
	return w2.NewSaveFormResponse(req.Record.ID).Write(w)
}

func (h *Handler) PostStoreThemeRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.DeleteStoreThemes(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}
