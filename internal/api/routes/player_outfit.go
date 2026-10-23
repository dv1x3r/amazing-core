package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dv1x3r/amazing-core/internal/api/services"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
)

func (h *Handler) GetPlayerOutfitGrid(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.GetPlayerOutfitGrid(r.Context(), req, id)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) GetPlayerOutfitDropdown(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.GetPlayerOutfitsDropdown(r.Context(), req, id)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostPlayerOutfitForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[services.PlayerOutfit](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if req.Record.ID == 0 {
		req.Record.ID, err = h.svc.CreatePlayerOutfit(r.Context(), req)
	} else {
		err = h.svc.UpdatePlayerOutfit(r.Context(), req)
	}
	if errors.Is(err, services.ErrPlayerOutfitExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return err
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

func (h *Handler) PostPlayerOutfitRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.DeletePlayerOutfits(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}
