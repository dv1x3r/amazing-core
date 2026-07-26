package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dv1x3r/amazing-core/internal/api/services"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
)

func (h *Handler) GetContainer(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetDropdownRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.GetContainersDropdown(r.Context(), req)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) GetContainerGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.GetContainerGrid(r.Context(), req)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostContainerForm(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveFormRequest[services.Container](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if req.Record.ID == 0 {
		req.Record.ID, err = h.svc.CreateContainer(r.Context(), req)
	} else {
		err = h.svc.UpdateContainer(r.Context(), req)
	}
	if errors.Is(err, services.ErrContainerExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return err
	}
	return w2.NewSaveFormResponse(req.Record.ID).Write(w)
}

func (h *Handler) PostContainerRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.DeleteContainers(r.Context(), req)
	if errors.Is(err, services.ErrContainerInUse) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) GetContainerAssetGrid(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.GetContainerAssetGrid(r.Context(), req, id)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostContainerAssetForm(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseSaveFormRequest[services.ContainerAsset](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req.Record.ID, err = h.svc.AddContainerAsset(r.Context(), req, id)
	if errors.Is(err, services.ErrContainerAssetExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return err
	}
	return w2.NewSaveFormResponse(req.Record.ID).Write(w)
}

func (h *Handler) PostContainerAssetGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[services.ContainerAsset](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.UpdateContainerAssets(r.Context(), req)
	if errors.Is(err, services.ErrContainerAssetExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerAssetRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err = h.svc.DeleteContainerAssets(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerAssetReorder(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseReorderGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err = h.svc.ReorderContainerAssets(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) GetContainerPackageGrid(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseGetGridRequest(r.URL.Query().Get("request"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	res, err := h.svc.GetContainerPackageGrid(r.Context(), req, id)
	if err != nil {
		return err
	}
	return res.Write(w)
}

func (h *Handler) PostContainerPackageForm(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req, err := w2.ParseSaveFormRequest[services.ContainerPackage](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	req.Record.ID, err = h.svc.AddContainerPackage(r.Context(), req, id)
	if errors.Is(err, services.ErrContainerPackageExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if errors.Is(err, services.ErrPackageCyclicDependency) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return err
	}
	return w2.NewSaveFormResponse(req.Record.ID).Write(w)
}

func (h *Handler) PostContainerPackageGrid(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseSaveGridRequest[services.ContainerPackage](r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	err = h.svc.UpdateContainerPackages(r.Context(), req)
	if errors.Is(err, services.ErrContainerPackageExists) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if errors.Is(err, services.ErrPackageCyclicDependency) {
		return wrap.WithHTTPStatus(err, http.StatusConflict)
	} else if err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerPackageRemove(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseRemoveGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err = h.svc.DeleteContainerPackages(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}

func (h *Handler) PostContainerPackageReorder(w http.ResponseWriter, r *http.Request) error {
	req, err := w2.ParseReorderGridRequest(r.Body)
	if err != nil {
		return wrap.WithHTTPStatus(err, http.StatusBadRequest)
	}
	if err = h.svc.ReorderContainerPackages(r.Context(), req); err != nil {
		return err
	}
	return w2.NewSuccessResponse().Write(w, http.StatusOK)
}
