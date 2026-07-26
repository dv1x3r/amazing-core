package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetInventoryObjects(ctx context.Context, platform gsf.Platform, playerOID types.OID) ([]types.PlayerItem, error) {
	const op = "game.Service.GetInventoryObjects"

	rows, err := s.store.DB().QueryContext(ctx, `
			select
				pi.item_id,
				pi.gsfoid as player_item_gsfoid,
				pl.gsfoid as player_gsfoid,
				pi.quantity as quantity
			from player_item as pi
			join player as pl on pl.id = pi.player_id
			where pl.gsfoid = ? and pi.player_avatar_outfit_id is null;
		`, playerOID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	playerItems := []types.PlayerItem{}
	for rows.Next() {
		var itemID int
		var playerItem types.PlayerItem
		if err := rows.Scan(
			&itemID,
			&playerItem.OID,
			&playerItem.PlayerOID,
			&playerItem.Quantity,
		); err != nil {
			return playerItems, wrap.IfErr(op, err)
		}
		item, err := s.GetItem(ctx, platform, itemID)
		if err != nil {
			return playerItems, wrap.IfErr(op, err)
		}
		playerItem.Item = item
		playerItems = append(playerItems, playerItem)
	}

	return playerItems, wrap.IfErr(op, rows.Err())
}

func (s *Service) GetOutfitItems(ctx context.Context, platform gsf.Platform, playerAvatarOutfitOID, playerOID types.OID) ([]types.PlayerItem, error) {
	const op = "game.Service.GetOutfitItems"

	rows, err := s.store.DB().QueryContext(ctx, `
			select
				pi.item_id,
				pi.gsfoid as player_item_gsfoid,
				pl.gsfoid as player_gsfoid,
				po.gsfoid as outfit_gsfoid,
				avs.gsfoid as slot_gsfoid,
				pi.quantity as quantity
			from player_item as pi
			join player as pl on pl.id = pi.player_id
			join player_avatar_outfit as po on po.id = pi.player_avatar_outfit_id
			join avatar_slot as avs on avs.id = pi.avatar_slot_id
			where po.gsfoid = ? and pl.gsfoid = ?;
		`, playerAvatarOutfitOID, playerOID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	playerItems := []types.PlayerItem{}
	for rows.Next() {
		var itemID int
		var playerItem types.PlayerItem
		if err := rows.Scan(
			&itemID,
			&playerItem.OID,
			&playerItem.PlayerOID,
			&playerItem.PlayerAvatarOutfitOID,
			&playerItem.SlotOID,
			&playerItem.Quantity,
		); err != nil {
			return playerItems, wrap.IfErr(op, err)
		}
		item, err := s.GetItem(ctx, platform, itemID)
		if err != nil {
			return playerItems, wrap.IfErr(op, err)
		}
		playerItem.Item = item
		playerItems = append(playerItems, playerItem)
	}

	return playerItems, wrap.IfErr(op, rows.Err())
}

func (s *Service) AddOutfitItems(ctx context.Context, outfitOID types.OID, inventoryOIDs, slotOIDs []types.OID) error {
	const op = "game.Service.AddOutfitItems"

	if len(inventoryOIDs) == 0 || len(slotOIDs) == 0 || len(inventoryOIDs) != len(slotOIDs) {
		return wrap.IfErr(op, fmt.Errorf("invalid OIDs length"))
	}

	values := make([]string, len(inventoryOIDs))
	args := make([]any, 0, len(inventoryOIDs)*2+1)
	for i := range inventoryOIDs {
		values[i] = "(?, ?)"
		args = append(args, inventoryOIDs[i], slotOIDs[i])
	}
	args = append(args, outfitOID)

	query := fmt.Sprintf(`
		with add_items(inventory_oid, slot_oid) as (values %s)
		update player_item
		set
			player_avatar_outfit_id = (select id from player_avatar_outfit where gsfoid = ?),
			avatar_slot_id = new.avatar_slot_id
		from (
			select
				pi.id as player_item_id,
				avs.id as avatar_slot_id
			from add_items
				join player_item as pi on pi.gsfoid = add_items.inventory_oid
				join avatar_slot as avs on avs.gsfoid = add_items.slot_oid
		) new
		where player_item.id = new.player_item_id;
	`, strings.Join(values, ", "))

	if _, err := s.store.DB().ExecContext(ctx, query, args...); err != nil {
		return wrap.IfErr(op, err)
	}

	return nil
}

func (s *Service) RemoveOutfitItems(ctx context.Context, outfitOID types.OID, inventoryOIDs []types.OID) error {
	const op = "game.Service.RemoveOutfitItems"

	if len(inventoryOIDs) == 0 {
		return wrap.IfErr(op, fmt.Errorf("invalid OIDs length"))
	}

	values := make([]string, len(inventoryOIDs))
	args := make([]any, 0, len(inventoryOIDs)+1)
	for i := range inventoryOIDs {
		values[i] = "(?)"
		args = append(args, inventoryOIDs[i])
	}
	args = append(args, outfitOID)

	query := fmt.Sprintf(`
		with remove_items(inventory_oid) as (values %s)
		update player_item
		set
			player_avatar_outfit_id = null,
			avatar_slot_id = null
		from (
			select pi.id as player_item_id
			from remove_items
				join player_item as pi on pi.gsfoid = remove_items.inventory_oid
		) remove
		where player_item.id = remove.player_item_id;
	`, strings.Join(values, ", "))

	if _, err := s.store.DB().ExecContext(ctx, query, args...); err != nil {
		return wrap.IfErr(op, err)
	}

	return nil
}

func (s *Service) ReplaceOutfitItems(ctx context.Context, oldOIDs, newOIDs []types.OID) error {
	const op = "game.Service.ReplaceOutfitItems"

	if len(oldOIDs) == 0 || len(newOIDs) == 0 || len(oldOIDs) != len(newOIDs) {
		return wrap.IfErr(op, fmt.Errorf("invalid OIDs length"))
	}

	values := make([]string, len(oldOIDs))
	args := make([]any, 0, len(oldOIDs)*2+1)
	for i := range oldOIDs {
		values[i] = "(?, ?)"
		args = append(args, oldOIDs[i], newOIDs[i])
	}

	query := fmt.Sprintf(`
		with pairs(old_oid, new_oid) as (values %s)
		update player_item
		set
			player_avatar_outfit_id = swap.player_avatar_outfit_id,
			avatar_slot_id = swap.avatar_slot_id
		from (
			select
				old_item.id as id,
				new_item.player_avatar_outfit_id,
				new_item.avatar_slot_id
			from pairs
				join player_item as old_item on old_item.gsfoid = pairs.old_oid
				join player_item as new_item on new_item.gsfoid = pairs.new_oid
			union all
			select
				new_item.id as id,
				old_item.player_avatar_outfit_id,
				old_item.avatar_slot_id
			from pairs
				join player_item as old_item on old_item.gsfoid = pairs.old_oid
				join player_item as new_item on new_item.gsfoid = pairs.new_oid
		) swap
		where player_item.id = swap.id;
	`, strings.Join(values, ", "))

	if _, err := s.store.DB().ExecContext(ctx, query, args...); err != nil {
		return wrap.IfErr(op, err)
	}

	return nil
}
