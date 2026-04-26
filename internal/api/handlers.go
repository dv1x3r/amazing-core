package api

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/services"
	"github.com/dv1x3r/amazing-core/internal/services/asset"
	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/amazing-core/internal/services/avatar"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
	"github.com/dv1x3r/amazing-core/internal/services/dummy"
	"github.com/dv1x3r/amazing-core/internal/services/player"
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
	svc services.Set
}

func NewHandler(svc services.Set) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) error {
	username, ok := h.svc.Auth.GetSessionUsername(w, r)
	if ok {
		if err := h.svc.Auth.RefreshSession(w, r); err != nil {
			return err
		}
	}
	version := config.Get().Version
	data := map[string]any{"username": username, "version": version}
	return tmpl.ExecuteTemplate(w, "admin.gotmpl", data)
}

func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) error {
	form, err := w2.ParseSaveFormRequest[auth.AdminLoginForm](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if _, err := h.svc.Auth.AuthenticateSession(w, r, form.Record); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostLogout(w http.ResponseWriter, r *http.Request) error {
	if err := h.svc.Auth.DeauthenticateSession(w, r); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) GetBlob(w http.ResponseWriter, r *http.Request) error {
	data, err := h.svc.Blob.GetBlobFile(r.Context(), r.PathValue("cdnid"))
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
	res, err := h.svc.Asset.GetAssetsDropdown(r.Context(), req)
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
	res, err := h.svc.Asset.GetAssetGrid(r.Context(), req)
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
	if err := h.svc.Asset.UpdateAssets(r.Context(), req); err != nil {
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

	result, err := h.svc.Asset.ImportCacheItems(r.Context(), items)
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
	res, err := h.svc.Asset.GetFileTypesDropdown(r.Context(), req)
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
	res, err := h.svc.Asset.GetAssetTypesDropdown(r.Context(), req)
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
	res, err := h.svc.Asset.GetAssetGroupsDropdown(r.Context(), req)
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
	res, err := h.svc.Asset.GetContainersDropdown(r.Context(), req)
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
	res, err := h.svc.Asset.GetContainerGrid(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostContainerForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[asset.Container](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	_, err = h.svc.Asset.CreateContainer(r.Context(), req)
	if errors.Is(err, asset.ErrContainerExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

func (h *Handler) PostContainerGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[asset.Container](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.Asset.UpdateContainers(r.Context(), req)
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
	err = h.svc.Asset.DeleteContainers(r.Context(), req)
	if errors.Is(err, asset.ErrContainerInUse) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
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
	res, err := h.svc.Asset.GetContainerAssetGrid(r.Context(), req, id)
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
	err = h.svc.Asset.AddContainerAsset(r.Context(), req, id)
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
	err = h.svc.Asset.UpdateContainerAssets(r.Context(), req)
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
	if err = h.svc.Asset.DeleteContainerAssets(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerAssetReorder(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseReorderGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err = h.svc.Asset.ReorderContainerAssets(r.Context(), req); err != nil {
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
	res, err := h.svc.Asset.GetContainerPackageGrid(r.Context(), req, id)
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
	err = h.svc.Asset.AddContainerPackage(r.Context(), req, id)
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
	err = h.svc.Asset.UpdateContainerPackages(r.Context(), req)
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
	if err = h.svc.Asset.DeleteContainerPackages(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerPackageReorder(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseReorderGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err = h.svc.Asset.ReorderContainerPackages(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

// ── Avatars ──────────────────────────────────────────────────────────────────

func (h *Handler) GetAvatar(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Avatar.GetAvatarsDropdown(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) GetAvatarGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Avatar.GetAvatarGrid(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostAvatarForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[avatar.Avatar](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	_, err = h.svc.Avatar.CreateAvatar(r.Context(), req)
	if errors.Is(err, avatar.ErrAvatarExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

func (h *Handler) PostAvatarGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[avatar.Avatar](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.Avatar.UpdateAvatars(r.Context(), req)
	if errors.Is(err, avatar.ErrAvatarExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostAvatarRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.Avatar.DeleteAvatars(r.Context(), req); err != nil {
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
	res, err := h.svc.SiteFrame.GetSiteFrameGrid(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostSiteFrameForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[siteframe.SiteFrame](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	_, err = h.svc.SiteFrame.CreateSiteFrame(r.Context(), req)
	if errors.Is(err, siteframe.ErrSiteFrameExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

func (h *Handler) PostSiteFrameGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[siteframe.SiteFrame](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.SiteFrame.UpdateSiteFrames(r.Context(), req)
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
	if err := h.svc.SiteFrame.DeleteSiteFrames(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

// ── Dummy Parameters ─────────────────────────────────────────────────────────

func (h *Handler) GetDummyGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Dummy.GetDummyParametersGrid(r.Context(), req)
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
	if err := h.svc.Dummy.UpdateDummyParameters(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

// ── Players ──────────────────────────────────────────────────────────────────

func (h *Handler) GetPlayerGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Player.GetPlayerListGrid(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) GetPlayerForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetFormRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Player.GetPlayerDetailsForm(r.Context(), req)
	if errors.Is(err, player.ErrPlayerNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostPlayerForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[player.PlayerDetails](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.Player.UpdatePlayerDetails(r.Context(), req)
	if errors.Is(err, player.ErrPlayerExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if errors.Is(err, player.ErrPlayerNotFound) {
		return wrap.WithHTTPStatus(err, http.StatusNotFound)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

// ── Player Avatars ───────────────────────────────────────────────────────────

func (h *Handler) GetPlayerAvatar(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Player.GetPlayerAvatarsDropdown(r.Context(), req, id)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) GetPlayerAvatarGrid(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Player.GetPlayerAvatarGrid(r.Context(), req, id)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostPlayerAvatarForm(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseSaveFormRequest[player.PlayerAvatar](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.Player.CreatePlayerAvatar(r.Context(), req, id)
	if errors.Is(err, player.ErrPlayerAvatarExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if errors.Is(err, asset.ErrPackageCyclicDependency) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

func (h *Handler) PostPlayerAvatarGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[player.PlayerAvatar](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.Player.UpdatePlayerAvatars(r.Context(), req)
	if errors.Is(err, player.ErrPlayerAvatarExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostPlayerAvatarRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.Player.DeletePlayerAvatars(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

// ── Player Outfits ───────────────────────────────────────────────────────────

func (h *Handler) GetPlayerOutfitGrid(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Player.GetPlayerOutfitGrid(r.Context(), req, id)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) PostPlayerOutfitForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[player.PlayerOutfit](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.Player.CreatePlayerOutfit(r.Context(), req)
	if errors.Is(err, player.ErrPlayerOutfitExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSaveFormResponse(req.RecID).Write(w)
}

func (h *Handler) PostPlayerOutfitGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[player.PlayerOutfit](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.Player.UpdatePlayerOutfits(r.Context(), req)
	if errors.Is(err, player.ErrPlayerOutfitExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostPlayerOutfitRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.Player.DeletePlayerOutfits(r.Context(), req); err != nil {
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
	res, err := h.svc.RandName.GetRandomNameGrid(r.Context(), req)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return res.Write(w)
}

func (h *Handler) GetRandnameForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetFormRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.RandName.GetRandomNameForm(r.Context(), req)
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
		req.RecID, err = h.svc.RandName.CreateRandomName(r.Context(), req)
		if errors.Is(err, randname.ErrNameExists) {
			return wrap.WithHTTPStatus(err, http.StatusConflict)
		}
	} else {
		err = h.svc.RandName.UpdateRandomName(r.Context(), req)
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

func (h *Handler) PostRandnameRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.RandName.DeleteRandomNames(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

// ── Blob ─────────────────────────────────────────────────────────────────────

func (h *Handler) GetBlobGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Blob.GetBlobGrid(r.Context(), req)
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
	if err := h.svc.Blob.DeleteFiles(r.Context(), req); err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostBlobUpload(w http.ResponseWriter, r *http.Request) error {
	headers, err := w2file.ParseMultipartFiles(r)
	if err != nil {
		return err
	}
	err = h.svc.Blob.SaveFiles(r.Context(), headers)
	if errors.Is(err, blob.ErrFileExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostBlobImport(w http.ResponseWriter, r *http.Request) error {
	result, err := h.svc.Blob.ImportFromFolder(r.Context())
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf("Import completed: %d imported, %d skipped", result.ImportedFiles, result.SkippedFiles)
	return res.Write(w, http.StatusOK)
}

func (h *Handler) PostBlobExport(w http.ResponseWriter, r *http.Request) error {
	result, err := h.svc.Blob.ExportToFolder(r.Context())
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	res := w2.NewSuccessResponse()
	res.Message = fmt.Sprintf("Export completed: %d exported, %d skipped", result.ExportedFiles, result.SkippedFiles)
	return res.Write(w, http.StatusOK)
}
