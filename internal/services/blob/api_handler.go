package blob

import (
	"errors"
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
)

const defaultMemory = 32 << 20 // 32 MB

type APIHandler struct {
	service *Service
}

func NewAPIHandler(service *Service) *APIHandler {
	return &APIHandler{service: service}
}

func (h *APIHandler) GetBlob(w http.ResponseWriter, r *http.Request) error {
	blob, err := h.service.FetchFileBlob(r.Context(), r.PathValue("cdnid"))
	if errors.Is(err, ErrFileNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(blob)
	return err
}

func (h *APIHandler) GetRecords(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGridDataRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	records, total, err := h.service.FetchFilesList(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewGridDataResponse(records, total).Write(w)
}

func (h *APIHandler) PostUpload(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseMultipartForm(defaultMemory); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	err := h.service.SaveFiles(r.Context(), r.MultipartForm.File["files[]"])
	if errors.Is(err, ErrFileExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *APIHandler) PostRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGridRemoveRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := h.service.DeleteFiles(r.Context(), req.ID); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}
