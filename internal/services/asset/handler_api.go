package asset

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2file"
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

func (h *APIHandler) PostSave(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[AssetItem](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := h.service.SaveGrid(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	return res.Write(w, http.StatusOK)
}

func (h *APIHandler) PostRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := h.service.DeleteByGrid(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	return res.Write(w, http.StatusOK)
}

func (h *APIHandler) GetFileTypeDropdown(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.service.GetFileTypeDropdownRecords(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *APIHandler) GetAssetTypeDropdown(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.service.GetAssetTypeDropdownRecords(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *APIHandler) GetAssetGroupDropdown(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.service.GetAssetGroupDropdownRecords(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *APIHandler) PostCacheJSON(w http.ResponseWriter, r *http.Request) error {
	headers, err := w2file.ParseMultipartFiles(r)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if len(headers) != 1 {
		return wrap.WithHTTPStatus(errors.New("invalid file payload"), http.StatusBadRequest)
	}

	file, err := headers[0].Open()
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	defer file.Close()

	var items []CacheItem
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	result, err := ImportCacheItems(h.service.logger, h.service.store.DB(), items, true)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf("Import completed: %d assets, %d metadata", result.ImportedAssets, result.ImportedMetadata)
	return res.Write(w, http.StatusOK)
}

func (h *APIHandler) GetAssetsSQL(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.DumpTable(r.Context())
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(data)
	return err
}

func (h *APIHandler) PostAssetsSQL(w http.ResponseWriter, r *http.Request) error {
	headers, err := w2file.ParseMultipartFiles(r)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if len(headers) != 1 {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	file, err := headers[0].Open()
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	err = h.service.store.RecreateTableFromQuery(h.service.logger, string(data), "asset", true)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	res.Message = "Import completed"
	return res.Write(w, http.StatusOK)
}
