package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetItemCategories(ctx context.Context, publicOnly bool) ([]types.ItemCategory, error) {
	const op = "game.Service.GetItemCategories"

	rows, err := s.store.DB().QueryContext(ctx, `
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
		`, publicOnly)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	categories := []types.ItemCategory{}
	for rows.Next() {
		var category types.ItemCategory
		if err := rows.Scan(
			&category.OID,
			&category.ParentOID,
			&category.Name,
			&category.IsOutdoor,
			&category.IsWalkover,
			&category.ShowInDock,
		); err != nil {
			return categories, wrap.IfErr(op, err)
		}
		categories = append(categories, category)
	}

	return categories, wrap.IfErr(op, rows.Err())
}

func (s *Service) GetItemCategoriesByItemID(ctx context.Context, itemID int) ([]types.ItemCategory, error) {
	const op = "game.Service.GetItemCategoriesByItemID"

	rows, err := s.store.DB().QueryContext(ctx, `
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
		`, itemID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	categories := []types.ItemCategory{}
	for rows.Next() {
		var category types.ItemCategory
		if err := rows.Scan(
			&category.OID,
			&category.ParentOID,
			&category.Name,
			&category.IsOutdoor,
			&category.IsWalkover,
			&category.ShowInDock,
		); err != nil {
			return categories, wrap.IfErr(op, err)
		}
		categories = append(categories, category)
	}

	return categories, wrap.IfErr(op, rows.Err())
}
