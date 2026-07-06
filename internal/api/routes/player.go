package routes

import (
	"errors"
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/services/player"
	"github.com/dv1x3r/w2go/w2"
)

func (h *Handler) GetPlayerGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Player.GetPlayerListGrid(r.Context(), req)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) GetPlayerForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetFormRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Player.GetPlayerDetailsForm(r.Context(), req)
	if errors.Is(err, player.ErrPlayerNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostPlayerForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[player.PlayerDetails](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.Player.UpdatePlayerDetails(r.Context(), req)
	if errors.Is(err, player.ErrPlayerExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if errors.Is(err, player.ErrPlayerNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return err
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}
