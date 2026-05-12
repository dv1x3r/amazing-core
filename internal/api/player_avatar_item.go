package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/services/player"
	"github.com/dv1x3r/w2go/w2"
)

func (h *Handler) GetPlayerAvatarItemGrid(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Player.GetPlayerAvatarItemGrid(r.Context(), req, id)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostPlayerAvatarItemForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[player.PlayerAvatarItem](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if req.Record.ID == 0 {
		req.Record.ID, err = h.svc.Player.CreatePlayerAvatarItem(r.Context(), req)
	} else {
		err = h.svc.Player.UpdatePlayerAvatarItem(r.Context(), req)
	}
	if errors.Is(err, player.ErrPlayerItemAttached) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return err
	}
	return w2.NewSaveFormResponse(req.Record.ID).Write(w)
}

func (h *Handler) PostPlayerAvatarItemRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.Player.DeletePlayerAvatarItems(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}
