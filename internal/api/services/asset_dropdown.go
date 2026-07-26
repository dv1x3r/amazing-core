package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

func (s *Service) GetAssetsDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetAssetsDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:    "asset as a",
		IDField: "a.id",
		TextField: `
			concat_ws(' - ',
				coalesce(at.name, '[NO TYPE]'),
				a.gsfoid,
				coalesce(a.res_name, '[NULL]'),
				a.bundle_version
			)`,
		OrderByField: "at.name, a.gsfoid",
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_type as at", "at.id = a.asset_type_id")
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetFileTypesDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetFileTypesDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "file_type",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "name",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetAssetTypesDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetAssetTypesDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "asset_type",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "name",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetAssetGroupsDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetAssetGroupsDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "asset_group",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "name",
	})
	return res, wrap.IfErr(op, err)
}
