package routes

import (
	"encoding/json"
	"net/http"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/services/asset"
	"github.com/dv1x3r/w2go/w2"
)

func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Asset.GetAssetsDropdown(r.Context(), req)
	if err != nil {
		return err
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
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostAssetGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[asset.Asset](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.Asset.UpdateAssets(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostAssetRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err := h.svc.Asset.DeleteAssets(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) GetAssetFileType(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.Asset.GetFileTypesDropdown(r.Context(), req)
	if err != nil {
		return err
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
		return err
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
		return err
	}
	return res.Write(w)
}

func (h *Handler) GetAssetIndex(w http.ResponseWriter, r *http.Request) error {
	res, err := h.svc.Asset.GetAssetIndex(r.Context())
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="index.json"`)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(res)
}

func (h *Handler) PostAssetImport(w http.ResponseWriter, r *http.Request) error {
	err := h.svc.Asset.ImportAssets(r.Context())
	if err != nil {
		return err
	}
	res := w2.NewSuccessResponse()
	res.Message = "Import completed!"
	return res.Write(w, http.StatusOK)
}
