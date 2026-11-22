package services

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2db"
)

func (s *Service) GetItemCategories(ctx context.Context, publicOnly bool) ([]types.ItemCategory, error) {
	const op = "game.Service.GetItemCategories"
	categories, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[types.ItemCategory]{
		Query: `
			select
				ic.gsfoid,
				icp.gsfoid as parent_gsfoid,
				ic.name,
				ic.is_outdoor,
				ic.is_walkover,
				ic.show_in_dock
			from item_category as ic
			left join item_category as icp on icp.id = ic.parent_id
			where (? = 0 or ic.is_public = 1);
		`,
		Args: []any{publicOnly},
		Scan: func(rows *sql.Rows, category *types.ItemCategory) error {
			return rows.Scan(
				&category.OID,
				&category.ParentOID,
				&category.Name,
				&category.IsOutdoor,
				&category.IsWalkover,
				&category.ShowInDock,
			)
		},
	})
	if categories == nil {
		categories = []types.ItemCategory{}
	}
	return categories, wrap.IfErr(op, err)
}

func (s *Service) GetItemCategoriesByItemID(ctx context.Context, itemID int) ([]types.ItemCategory, error) {
	const op = "game.Service.GetItemCategoriesByItemID"
	categories, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[types.ItemCategory]{
		Query: `
			select
				ic.gsfoid,
				icp.gsfoid as parent_gsfoid,
				ic.name,
				ic.is_outdoor,
				ic.is_walkover,
				ic.show_in_dock
			from item_category_map as icm
			join item_category as ic on ic.id = icm.category_id
			left join item_category as icp on icp.id = ic.parent_id
			where icm.item_id = ?;
		`,
		Args: []any{itemID},
		Scan: func(rows *sql.Rows, category *types.ItemCategory) error {
			return rows.Scan(
				&category.OID,
				&category.ParentOID,
				&category.Name,
				&category.IsOutdoor,
				&category.IsWalkover,
				&category.ShowInDock,
			)
		},
	})
	if categories == nil {
		categories = []types.ItemCategory{}
	}
	return categories, wrap.IfErr(op, err)
}
