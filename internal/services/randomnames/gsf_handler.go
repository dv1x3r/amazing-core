package randomnames

import (
	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/messages"
)

type GSFHandler struct {
	service *Service
}

func NewGSFHandler(service *Service) *GSFHandler {
	return &GSFHandler{service: service}
}

/*
GetRandomNames returns the array with random family names.

Scenarios:
  - New player registration:
  - "second_name" for random zing name.
  - "Family_1", "Family_2" and "Family_3" for random family name parts.

Returns:
  - The array with random strings based on requested amount and type.
*/
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
