package api

import (
	"errors"
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/services/randname"
	"github.com/dv1x3r/w2go/w2"
)

func (h *Handler) GetRandnameGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.RandName.GetRandomNameGrid(r.Context(), req)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) GetRandnameForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetFormRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.RandName.GetRandomNameForm(r.Context(), req)
	if errors.Is(err, randname.ErrNameNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostRandnameForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[randname.RandomName](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if req.RecID == 0 {
		req.RecID, err = h.svc.RandName.CreateRandomName(r.Context(), req)
		if errors.Is(err, randname.ErrNameExists) {
			return wrap.WithHTTPStatus(err, http.StatusConflict)
		} else if err != nil {
			return err
		}
	} else {
		err = h.svc.RandName.UpdateRandomName(r.Context(), req)
		if errors.Is(err, randname.ErrNameExists) {
			return wrap.WithHTTPStatus(err, http.StatusConflict)
		} else if err != nil {
			return err
		}
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

func (h *Handler) PostRandnameRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.RandName.DeleteRandomNames(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}
