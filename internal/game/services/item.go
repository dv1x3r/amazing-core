package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetItem(ctx context.Context, platform gsf.Platform, itemID int) (types.Item, error) {
	const op = "game.Service.GetItem"

	row := s.store.DB().QueryRowContext(ctx, `
			select
				container_id,
				name
			from item
			where id = ?;
		`, itemID)

	var containerID int
	var item types.Item

	if err := row.Scan(
		&containerID,
		&item.Name,
	); err != nil {
		return item, wrap.IfErr(op, err)
	}

	container, err := s.GetAssetContainer(ctx, platform, containerID)
	if err != nil {
		return item, wrap.IfErr(op, err)
	}
	item.AssetContainer = container

	categories, err := s.GetItemCategoriesByItemID(ctx, itemID)
	if err != nil {
		return item, wrap.IfErr(op, err)
	}
	item.ItemCategories = categories

	slots, err := s.GetItemAcceptableSlotsByItemID(ctx, itemID)
	if err != nil {
		return item, wrap.IfErr(op, err)
	}
	item.AcceptableSlotOIDs = slots

	return item, nil
}
