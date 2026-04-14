package randname

import (
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/messages"
)

type GSFHandler struct {
	service *Service
}

func NewGSFHandler(service *Service) *GSFHandler {
	return &GSFHandler{service: service}
}

func (h *GSFHandler) GetRandomNames(w gsf.ResponseWriter, r *gsf.Request) error {
	req := &messages.GetRandomNamesRequest{}
	if err := r.Read(req); err != nil {
		return err
	}

	names, err := h.service.GetNStringsByType(r.Context(), req.NamePartType, int(req.Amount))
	if err != nil {
		return err
	}

	res := &messages.GetRandomNamesResponse{}
	res.Names = names

	return w.Write(res)
}
