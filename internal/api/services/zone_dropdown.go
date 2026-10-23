package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

func (s *Service) GetZonesDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetZonesDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "zone as z",
		IDField:      "z.id",
		TextField:    "(ac.gsfoid || ' - ' || ac.name)",
		OrderByField: "ac.gsfoid",
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = z.container_id")
		},
	})
	return res, wrap.IfErr(op, err)
}
