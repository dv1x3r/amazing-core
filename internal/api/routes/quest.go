package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dv1x3r/amazing-core/internal/api/services"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
)

func (h *Handler) GetNPCQuestGrid(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.GetQuestGrid(r.Context(), req, id)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostNPCQuestForm(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseSaveFormRequest[services.Quest](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if req.Record.ID == 0 {
		req.Record.ID, err = h.svc.CreateQuest(r.Context(), req, id)
	} else {
		err = h.svc.UpdateQuest(r.Context(), req)
	}
	if errors.Is(err, services.ErrQuestExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return err
	}
	return w2.NewSaveFormResponse(req.Record.ID).Write(w)
}

func (h *Handler) PostNPCQuestRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.DeleteQuests(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostNPCQuestReorder(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseReorderGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.ReorderQuests(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) GetQuestDropdown(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.GetQuestsDropdown(r.Context(), req)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) GetQuestStateDropdown(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.GetQuestStatesDropdown(r.Context(), req)
	if err != nil {
		return err
	}
	return res.Write(w)
}
