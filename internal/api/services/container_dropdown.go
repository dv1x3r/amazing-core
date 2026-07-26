package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
)

func (s *Service) GetContainersDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetContainersDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "asset_container",
		IDField:      "id",
		TextField:    "concat(gsfoid, ' - ', name, ' (' || ptag || ')')",
		OrderByField: "gsfoid",
	})
	return res, wrap.IfErr(op, err)
}
