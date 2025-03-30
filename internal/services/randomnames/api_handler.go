package randomnames

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

func (h *APIHandler) GetForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseFormGetRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	record, err := h.service.GetByID(r.Context(), req.RecID)
	if errors.Is(err, ErrNameNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewFormGetResponse(record).Write(w)
}

func (h *APIHandler) PostForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseFormSaveRequest[RandomName](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if req.RecID == 0 {
		req.RecID, err = h.service.Insert(r.Context(), req.Record)
	} else {
		err = h.service.UpdateByID(r.Context(), req.RecID, req.Record)
	}

	if errors.Is(err, ErrNameNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if errors.Is(err, ErrNameExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewFormSaveResponse(req.RecID).Write(w)
}

func (h *APIHandler) GetRecords(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGridDataRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	records, total, err := h.service.GetList(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewGridDataResponse(records, total).Write(w)
}

func (h *APIHandler) PostRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGridRemoveRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := h.service.Delete(r.Context(), req.ID); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}
