package api

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/services/asset"
	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
	"github.com/dv1x3r/amazing-core/internal/services/dummy"
	"github.com/dv1x3r/amazing-core/internal/services/randname"
	"github.com/dv1x3r/amazing-core/internal/services/siteframe"

	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2file"
)

//go:embed *.gotmpl
var templatesFS embed.FS

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFS(templatesFS, "*.gotmpl"))
}

type Handler struct {
	authService      *auth.Service
	assetService     *asset.Service
	blobService      *blob.Service
	siteFrameService *siteframe.Service
	dummyService     *dummy.Service
	randnameService  *randname.Service
}

func NewHandler(
	authService *auth.Service,
	assetService *asset.Service,
	blobService *blob.Service,
	siteFrameService *siteframe.Service,
	dummyService *dummy.Service,
	randnameService *randname.Service,
) *Handler {
	return &Handler{
		authService:      authService,
		assetService:     assetService,
		blobService:      blobService,
		siteFrameService: siteFrameService,
		dummyService:     dummyService,
		randnameService:  randnameService,
	}
}

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) error {
	username, ok := h.authService.GetSessionUsername(w, r)
	if ok {
		if err := h.authService.RefreshSession(w, r); err != nil {
			return err
		}
	}
	data := map[string]any{"username": username}
	return tmpl.ExecuteTemplate(w, "admin.gotmpl", data)
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

func (h *Handler) GetBlob(w http.ResponseWriter, r *http.Request) error {
	data, err := h.blobService.GetBlobFile(r.Context(), r.PathValue("cdnid"))
	if errors.Is(err, blob.ErrFileNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(data)
	return err
}

// ── Asset ────────────────────────────────────────────────────────────────────

func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.assetService.GetAssetsDropdown(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) GetAssetGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.assetService.GetAssetGrid(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostAssetGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[asset.Asset](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.assetService.UpdateAssets(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
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

// ── Asset ENums ──────────────────────────────────────────────────────────────

func (h *Handler) GetAssetFileType(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.assetService.GetFileTypesDropdown(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) GetAssetType(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.assetService.GetAssetTypesDropdown(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) GetAssetGroup(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.assetService.GetAssetGroupsDropdown(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

// ── Containers ───────────────────────────────────────────────────────────────

func (h *Handler) GetContainer(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.assetService.GetContainersDropdown(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) GetContainerGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.assetService.GetContainerGrid(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostContainerGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[asset.Container](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.assetService.UpdateContainers(r.Context(), req)
	if errors.Is(err, asset.ErrContainerExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.assetService.DeleteContainers(r.Context(), req)
	if errors.Is(err, asset.ErrContainerInUse) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[asset.Container](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	_, err = h.assetService.CreateContainer(r.Context(), req)
	if errors.Is(err, asset.ErrContainerExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

// ── Container Assets ─────────────────────────────────────────────────────────

func (h *Handler) GetContainerAssetGrid(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.assetService.GetContainerAssetGrid(r.Context(), req, id)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostContainerAssetForm(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseSaveFormRequest[asset.ContainerAsset](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.assetService.AddContainerAsset(r.Context(), req, id)
	if errors.Is(err, asset.ErrContainerAssetExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

func (h *Handler) PostContainerAssetGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[asset.ContainerAsset](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.assetService.UpdateContainerAssets(r.Context(), req)
	if errors.Is(err, asset.ErrContainerAssetExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerAssetRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err = h.assetService.DeleteContainerAssets(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerAssetReorder(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseReorderGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err = h.assetService.ReorderContainerAssets(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

// ── Container Packages ───────────────────────────────────────────────────────

func (h *Handler) GetContainerPackageGrid(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.assetService.GetContainerPackageGrid(r.Context(), req, id)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostContainerPackageForm(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseSaveFormRequest[asset.ContainerPackage](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.assetService.AddContainerPackage(r.Context(), req, id)
	if errors.Is(err, asset.ErrContainerPackageExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if errors.Is(err, asset.ErrPackageCyclicDependency) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

func (h *Handler) PostContainerPackageGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[asset.ContainerPackage](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.assetService.UpdateContainerPackages(r.Context(), req)
	if errors.Is(err, asset.ErrContainerPackageExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if errors.Is(err, asset.ErrPackageCyclicDependency) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerPackageRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err = h.assetService.DeleteContainerPackages(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerPackageReorder(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseReorderGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err = h.assetService.ReorderContainerPackages(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

// ── Site Frame ───────────────────────────────────────────────────────────────

func (h *Handler) GetSiteFrameGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.siteFrameService.GetSiteFrameGrid(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostSiteFrameGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[siteframe.SiteFrame](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.siteFrameService.UpdateSiteFrames(r.Context(), req)
	if errors.Is(err, siteframe.ErrSiteFrameExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostSiteFrameRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.siteFrameService.DeleteSiteFrames(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostSiteFrameForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[siteframe.SiteFrame](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	_, err = h.siteFrameService.CreateSiteFrame(r.Context(), req)
	if errors.Is(err, siteframe.ErrSiteFrameExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

// ── Dummy Parameters ─────────────────────────────────────────────────────────

func (h *Handler) GetDummyGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.dummyService.GetDummyParametersGrid(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostDummyGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[dummy.Param](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.dummyService.UpdateDummyParameters(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

// ── Random Names ─────────────────────────────────────────────────────────────

func (h *Handler) GetRandnameGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.randnameService.GetRandomNameGrid(r.Context(), req)
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
	if err := h.randnameService.DeleteRandomNames(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) GetRandnameForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetFormRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.randnameService.GetRandomNameForm(r.Context(), req)
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
		req.RecID, err = h.randnameService.CreateRandomName(r.Context(), req)
		if errors.Is(err, randname.ErrNameExists) {
			return wrap.WithHTTPStatus(err, http.StatusConflict)
		}
	} else {
		err = h.randnameService.UpdateRandomName(r.Context(), req)
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

// ── Blob ─────────────────────────────────────────────────────────────────────

func (h *Handler) GetBlobGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.blobService.GetBlobGrid(r.Context(), req)
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
