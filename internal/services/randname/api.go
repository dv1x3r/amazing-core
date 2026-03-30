package randname

import (
	"errors"
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

func (h *APIHandler) GetRecords(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.service.GetGridRecords(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *APIHandler) PostRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := h.service.DeleteByGridIDs(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	return res.Write(w, http.StatusOK)
}

func (h *APIHandler) GetForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetFormRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.service.GetByFormID(r.Context(), req)
	if errors.Is(err, ErrNameNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *APIHandler) PostForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[RandomName](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if req.RecID == 0 {
		req.RecID, err = h.service.InsertFormRecord(r.Context(), req)
		if errors.Is(err, ErrNameExists) {
			return wrap.WithHTTPStatus(err, http.StatusConflict)
		}
	} else {
		err = h.service.UpdateFormRecord(r.Context(), req)
		if errors.Is(err, ErrNameExists) {
			return wrap.WithHTTPStatus(err, http.StatusConflict)
		} else if errors.Is(err, ErrNameNotFound) {
			return wrap.WithHTTPStatus(err, http.StatusNotFound)
		}
	}

	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSaveFormResponse(req.RecID)
	return res.Write(w)
}
