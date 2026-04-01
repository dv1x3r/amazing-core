package dummy

import (
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
	res, err := h.service.GetForm(r.Context())
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *APIHandler) PostForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[DummyConfig](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	err = h.service.UpdateForm(r.Context(), req.Record)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSaveFormResponse(req.RecID)
	return res.Write(w)
}
