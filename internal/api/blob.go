package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
	"github.com/dv1x3r/w2go/w2"
)

func (h *Handler) GetBlob(w http.ResponseWriter, r *http.Request) error {
	data, err := h.svc.Blob.GetBlobData(r.Context(), r.PathValue("cdnid"))
	if errors.Is(err, blob.ErrFileNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(data)
	return err
}

func (h *Handler) GetBlobGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Blob.GetBlobGrid(r.Context(), req)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostBlobRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.Blob.DeleteFiles(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostBlobImport(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[blob.ImportOptions](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	result, err := h.svc.Blob.ImportFromFolder(r.Context(), req.Record)
	if err != nil {
		return err
	}
	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf(
		"Import completed: %d imported, %d skipped, %d metadata generated",
		result.ImportedFiles, result.SkippedFiles, result.GeneratedMetadata,
	)
	return res.Write(w, http.StatusOK)
}

func (h *Handler) PostBlobExtract(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[blob.ExtractOptions](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	result, err := h.svc.Blob.ExtractFiles(r.Context(), req.Record)
	if err != nil {
		return err
	}
	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf("Extraction completed: %d files extracted", result.ExtractedFiles)
	return res.Write(w, http.StatusOK)
}
