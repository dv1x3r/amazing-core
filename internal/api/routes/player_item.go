package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/services/player"
	"github.com/dv1x3r/w2go/w2"
)

func (h *Handler) GetPlayerItem(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Player.GetPlayerItemsDropdown(r.Context(), req, id)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) GetPlayerItemGrid(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Player.GetPlayerItemGrid(r.Context(), req, id)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostPlayerItemForm(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseSaveFormRequest[player.PlayerItem](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if req.Record.ID == 0 {
		req.Record.ID, err = h.svc.Player.CreatePlayerItem(r.Context(), req, id)
	} else {
		err = h.svc.Player.UpdatePlayerItem(r.Context(), req)
	}
	if errors.Is(err, player.ErrPlayerItemExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return err
	}
	return w2.NewSaveFormResponse(req.Record.ID).Write(w)
}

func (h *Handler) PostPlayerItemRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.Player.DeletePlayerItems(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}
