package blob

import (
	"encoding/json"
	"errors"
	"fmt"
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

	res := w2.NewGridDataResponse(records, total)
	return res.Write(w)
}

func (h *APIHandler) PostUpload(w http.ResponseWriter, r *http.Request) error {
	const defaultMemory = 32 << 20 // 32 MB
	if err := r.ParseMultipartForm(defaultMemory); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	err := h.service.SaveFiles(r.Context(), r.MultipartForm.File["files[]"])
	if errors.Is(err, ErrFileTooLarge) {
		return wrap.WithHTTPStatus(err, http.StatusRequestEntityTooLarge)
	} else if errors.Is(err, ErrFileExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	return res.Write(w, http.StatusOK)
}

func (h *APIHandler) PostRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGridRemoveRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := h.service.DeleteFiles(r.Context(), req.ID); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	return res.Write(w, http.StatusOK)
}

const defaultCacheDir = "cache"

func (h *APIHandler) PostImport(w http.ResponseWriter, r *http.Request) error {
	result, err := h.service.ImportFromFolder(r.Context(), defaultCacheDir)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf("Import completed: %d imported, %d skipped", result.ImportedFiles, result.SkippedFiles)
	return res.Write(w, http.StatusOK)
}

func (h *APIHandler) PostExport(w http.ResponseWriter, r *http.Request) error {
	result, err := h.service.ExportFromFolder(r.Context(), defaultCacheDir, true)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf("Export completed: %d exported, %d skipped", result.ExportedFiles, result.SkippedFiles)
	return res.Write(w, http.StatusOK)
}

func (h *APIHandler) PostS3Sync(w http.ResponseWriter, r *http.Request) error {
	cfg := S3Config{}
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	result, err := h.service.SyncToS3(r.Context(), cfg)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf("Sync completed: %d uploaded, %d skipped", result.SyncedFiles, result.SkippedFiles)
	return res.Write(w, http.StatusOK)
}
