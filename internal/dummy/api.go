package dummy

import (
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
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
