package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2db"
)

func (s *Service) GetInventoryObjects(ctx context.Context, platform gsf.Platform, playerOID types.OID) ([]types.PlayerItem, error) {
	const op = "game.Service.GetInventoryObjects"

	type row struct {
		itemID     int
		playerItem types.PlayerItem
	}
	itemRows, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[row]{
		Query: `
			select
				pi.item_id,
				pi.gsfoid as player_item_gsfoid,
				pl.gsfoid as player_gsfoid,
				pi.quantity as quantity
			from player_item as pi
			join player as pl on pl.id = pi.player_id
			where pl.gsfoid = ? and pi.player_avatar_outfit_id is null
				and pi.player_maze_id is null and pi.parent_player_item_id is null;
		`,
		Args: []any{playerOID},
		Scan: func(rows *sql.Rows, r *row) error {
			return rows.Scan(
				&r.itemID,
				&r.playerItem.OID,
				&r.playerItem.PlayerOID,
				&r.playerItem.Quantity,
			)
		},
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	playerItems := make([]types.PlayerItem, 0, len(itemRows))
	for _, r := range itemRows {
		item, err := s.GetItem(ctx, platform, r.itemID)
		if err != nil {
			return playerItems, wrap.IfErr(op, err)
		}
		playerItem := r.playerItem
		playerItem.Item = item
		playerItems = append(playerItems, playerItem)
	}

	return playerItems, nil
}

func (s *Service) GetOutfitItems(ctx context.Context, platform gsf.Platform, playerAvatarOutfitOID, playerOID types.OID) ([]types.PlayerItem, error) {
	const op = "game.Service.GetOutfitItems"

	type row struct {
		itemID     int
		playerItem types.PlayerItem
	}
	itemRows, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[row]{
		Query: `
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
		`,
		Args: []any{playerAvatarOutfitOID, playerOID},
		Scan: func(rows *sql.Rows, r *row) error {
			return rows.Scan(
				&r.itemID,
				&r.playerItem.OID,
				&r.playerItem.PlayerOID,
				&r.playerItem.PlayerAvatarOutfitOID,
				&r.playerItem.SlotOID,
				&r.playerItem.Quantity,
			)
		},
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	playerItems := make([]types.PlayerItem, 0, len(itemRows))
	for _, r := range itemRows {
		item, err := s.GetItem(ctx, platform, r.itemID)
		if err != nil {
			return playerItems, wrap.IfErr(op, err)
		}
		playerItem := r.playerItem
		playerItem.Item = item
		playerItems = append(playerItems, playerItem)
	}

	return playerItems, nil
}

// AddOutfitItems puts the items into the outfit, only the items and the outfit of the player.
func (s *Service) AddOutfitItems(ctx context.Context, playerOID, outfitOID types.OID, inventoryOIDs, slotOIDs []types.OID) error {
	const op = "game.Service.AddOutfitItems"

	if len(inventoryOIDs) == 0 || len(inventoryOIDs) != len(slotOIDs) {
		return wrap.IfErr(op, fmt.Errorf("invalid OIDs length"))
	}

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	for i, inventoryOID := range inventoryOIDs {
		_, err := tx.ExecContext(ctx, `
				update player_item
				set player_avatar_outfit_id = po.id, avatar_slot_id = avs.id
				from player_avatar_outfit as po
				join player_avatar as pa on pa.id = po.player_avatar_id
				join avatar_slot as avs on avs.gsfoid = ?
				where po.gsfoid = ? and pa.player_id = player_item.player_id
					and player_item.gsfoid = ? and player_item.player_id = (select id from player where gsfoid = ?);
			`, slotOIDs[i], outfitOID, inventoryOID, playerOID)
		if err != nil {
			return wrap.IfErr(op, err)
		}
	}

	return wrap.IfErr(op, tx.Commit())
}

// RemoveOutfitItems takes the items of the player out of their outfits.
func (s *Service) RemoveOutfitItems(ctx context.Context, playerOID types.OID, inventoryOIDs []types.OID) error {
	const op = "game.Service.RemoveOutfitItems"

	if len(inventoryOIDs) == 0 {
		return wrap.IfErr(op, fmt.Errorf("invalid OIDs length"))
	}

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	for _, inventoryOID := range inventoryOIDs {
		_, err := tx.ExecContext(ctx, `
				update player_item
				set player_avatar_outfit_id = null, avatar_slot_id = null
				where gsfoid = ? and player_id = (select id from player where gsfoid = ?);
			`, inventoryOID, playerOID)
		if err != nil {
			return wrap.IfErr(op, err)
		}
	}

	return wrap.IfErr(op, tx.Commit())
}

// ReplaceOutfitItems swaps the outfits and slots of item pairs, both items of a pair must be of the player.
func (s *Service) ReplaceOutfitItems(ctx context.Context, playerOID types.OID, oldOIDs, newOIDs []types.OID) error {
	const op = "game.Service.ReplaceOutfitItems"

	if len(oldOIDs) == 0 || len(oldOIDs) != len(newOIDs) {
		return wrap.IfErr(op, fmt.Errorf("invalid OIDs length"))
	}

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	// Both rows of a pair read the values of the other before the update.
	for i, oldOID := range oldOIDs {
		_, err := tx.ExecContext(ctx, `
				update player_item
				set player_avatar_outfit_id = other.player_avatar_outfit_id, avatar_slot_id = other.avatar_slot_id
				from player_item as other, player as pl
				where pl.gsfoid = ? and player_item.player_id = pl.id and other.player_id = pl.id
					and ((player_item.gsfoid = ? and other.gsfoid = ?) or (player_item.gsfoid = ? and other.gsfoid = ?));
			`, playerOID, oldOID, newOIDs[i], newOIDs[i], oldOID)
		if err != nil {
			return wrap.IfErr(op, err)
		}
	}

	return wrap.IfErr(op, tx.Commit())
}
