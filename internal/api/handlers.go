package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/services/asset"
	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
	"github.com/dv1x3r/amazing-core/internal/services/dummy"
	"github.com/dv1x3r/amazing-core/internal/services/randname"

	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2file"
)

type Handler struct {
	authService     *auth.Service
	assetService    *asset.Service
	blobService     *blob.Service
	dummyService    *dummy.Service
	randnameService *randname.Service
}

func NewHandler(
	authService *auth.Service,
	dummyService *dummy.Service,
	blobService *blob.Service,
	assetService *asset.Service,
	randnameService *randname.Service,
) *Handler {
	return &Handler{
		authService:     authService,
		assetService:    assetService,
		blobService:     blobService,
		dummyService:    dummyService,
		randnameService: randnameService,
	}
}

// ── Admin ────────────────────────────────────────────────────────────────────

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	username, ok := h.authService.GetSessionUsername(w, r)
	if ok {
		if err := h.authService.RefreshSession(w, r); err != nil {
			w2.NewErrorResponse(err.Error()).Write(w, http.StatusInternalServerError)
			return
		}
	}

	data := map[string]any{"username": username}
	if err := tmpl.ExecuteTemplate(w, "admin.gotmpl", data); err != nil {
		w2.NewErrorResponse(err.Error()).Write(w, http.StatusInternalServerError)
	}
}

func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) error {
	form, err := w2.ParseSaveFormRequest[auth.AdminLoginForm](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if _, err := h.authService.AuthenticateSession(w, r, form.Record); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostLogout(w http.ResponseWriter, r *http.Request) error {
	if err := h.authService.DeauthenticateSession(w, r); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

// ── Asset ────────────────────────────────────────────────────────────────────

func (h *Handler) GetAssetRecords(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.assetService.GetGridRecords(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *Handler) PostAssetSave(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[asset.Asset](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := h.assetService.SaveGrid(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostAssetRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := h.assetService.DeleteByGrid(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) GetAssetFileTypeDropdown(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.assetService.GetFileTypeDropdownRecords(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *Handler) GetAssetTypeDropdown(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.assetService.GetAssetTypeDropdownRecords(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *Handler) GetAssetGroupDropdown(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.assetService.GetAssetGroupDropdownRecords(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *Handler) PostAssetCacheJSON(w http.ResponseWriter, r *http.Request) error {
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

	var items []asset.CacheItem
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	result, err := h.assetService.ImportCacheItems(r.Context(), items)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf("Import completed: %d assets, %d metadata", result.ImportedAssets, result.ImportedMetadata)
	return res.Write(w, http.StatusOK)
}

// ── Blob ─────────────────────────────────────────────────────────────────────

func (h *Handler) GetBlob(w http.ResponseWriter, r *http.Request) error {
	data, err := h.blobService.FetchFileBlob(r.Context(), r.PathValue("cdnid"))
	if errors.Is(err, blob.ErrFileNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(data)
	return err
}

func (h *Handler) GetBlobRecords(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.blobService.FetchFilesList(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *Handler) PostBlobRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := h.blobService.DeleteFiles(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostBlobUpload(w http.ResponseWriter, r *http.Request) error {
	headers, err := w2file.ParseMultipartFiles(r)
	if err != nil {
		return err
	}

	err = h.blobService.SaveFiles(r.Context(), headers)
	if errors.Is(err, blob.ErrFileExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostBlobImport(w http.ResponseWriter, r *http.Request) error {
	result, err := h.blobService.ImportFromFolder(r.Context())
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf("Import completed: %d imported, %d skipped", result.ImportedFiles, result.SkippedFiles)
	return res.Write(w, http.StatusOK)
}

func (h *Handler) PostBlobExport(w http.ResponseWriter, r *http.Request) error {
	result, err := h.blobService.ExportToFolder(r.Context())
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf("Export completed: %d exported, %d skipped", result.ExportedFiles, result.SkippedFiles)
	return res.Write(w, http.StatusOK)
}

func (h *Handler) PostBlobS3Sync(w http.ResponseWriter, r *http.Request) error {
	cfg := blob.S3Config{}
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	result, err := h.blobService.SyncToS3(r.Context(), cfg)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf("Sync completed: %d uploaded, %d skipped", result.SyncedFiles, result.SkippedFiles)
	return res.Write(w, http.StatusOK)
}

// ── Dummy ────────────────────────────────────────────────────────────────────

func (h *Handler) GetDummyForm(w http.ResponseWriter, r *http.Request) error {
	res, err := h.dummyService.GetForm(r.Context())
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostDummyForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[dummy.DummyConfig](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := h.dummyService.UpdateForm(r.Context(), req.Record); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

// ── Random Names ─────────────────────────────────────────────────────────────

func (h *Handler) GetRandnameRecords(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.randnameService.GetGridRecords(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *Handler) PostRandnameRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := h.randnameService.DeleteByGridIDs(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) GetRandnameForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetFormRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	res, err := h.randnameService.GetByFormID(r.Context(), req)
	if errors.Is(err, randname.ErrNameNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return res.Write(w)
}

func (h *Handler) PostRandnameForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[randname.RandomName](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if req.RecID == 0 {
		req.RecID, err = h.randnameService.InsertFormRecord(r.Context(), req)
		if errors.Is(err, randname.ErrNameExists) {
			return wrap.WithHTTPStatus(err, http.StatusConflict)
		}
	} else {
		err = h.randnameService.UpdateFormRecord(r.Context(), req)
		if errors.Is(err, randname.ErrNameExists) {
			return wrap.WithHTTPStatus(err, http.StatusConflict)
		} else if errors.Is(err, randname.ErrNameNotFound) {
			return wrap.WithHTTPStatus(err, http.StatusNotFound)
		}
	}

	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return w2.NewSaveFormResponse(req.RecID).Write(w)
}
