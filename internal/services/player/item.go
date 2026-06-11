package player

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type PlayerItem struct {
	ID           int              `json:"id"`
	OID          w2.Field[string] `json:"oid"`
	OIDStr       string           `json:"oid_str"`
	Item         w2.Dropdown      `json:"item"`
	PlayerAvatar w2.Dropdown      `json:"player_avatar"`
	PlayerOutfit w2.Dropdown      `json:"player_outfit"`
	AvatarSlot   w2.Dropdown      `json:"avatar_slot"`
	Quantity     w2.Field[int]    `json:"quantity"`
}

func (s *Service) GetPlayerItemGrid(ctx context.Context, req w2.GetGridRequest, playerID int) (w2.GetGridResponse[PlayerItem], error) {
	const op = "player.Service.GetPlayerItemGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerItem]{
		From: "player_item as pi",
		Select: []string{
			"pi.id",
			"pi.gsfoid",
			"pi.item_id",
			"it.name as item_name",
			"pi.player_avatar_id",
			"(pa.gsfoid || ' - ' || pa.name) as player_avatar_name",
			"pi.player_avatar_outfit_id",
			"(poa.name || ' / Outfit #' || po.outfit_no) as player_avatar_outfit_name",
			"pi.avatar_slot_id",
			"avs.slot as avatar_slot_name",
			"pi.quantity",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("item as it", "it.id = pi.item_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "player_avatar as pa", "pa.id = pi.player_avatar_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "player_avatar_outfit as po", "po.id = pi.player_avatar_outfit_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "player_avatar as poa", "poa.id = po.player_avatar_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "avatar_slot as avs", "avs.id = pi.avatar_slot_id")
			sb.Where(sb.EQ("pi.player_id", playerID))
			sb.OrderByAsc("pi.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerItem) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.Item.ID,
				&record.Item.Text,
				&record.PlayerAvatar.ID,
				&record.PlayerAvatar.Text,
				&record.PlayerOutfit.ID,
				&record.PlayerOutfit.Text,
				&record.AvatarSlot.ID,
				&record.AvatarSlot.Text,
				&record.Quantity,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromString(record.OID.V).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreatePlayerItem(ctx context.Context, req w2.SaveFormRequest[PlayerItem], playerID int) (int, error) {
	const op = "player.Service.CreatePlayerItem"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "player_item",
		Cols: []string{
			"player_id",
			"gsfoid",
			"item_id",
			"player_avatar_id",
			"player_avatar_outfit_id",
			"avatar_slot_id",
			"quantity",
		},
		Values: []any{
			playerID,
			sqlbuilder.Buildf("(select coalesce(max(gsfoid) + 1, 1) from player_item)"),
			req.Record.Item.ID,
			req.Record.PlayerAvatar.ID,
			req.Record.PlayerOutfit.ID,
			req.Record.AvatarSlot.ID,
			req.Record.Quantity,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrPlayerItemExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerItem(ctx context.Context, req w2.SaveFormRequest[PlayerItem]) error {
	const op = "player.Service.UpdatePlayerItem"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "player_item",
		Cols: []string{
			"gsfoid",
			"item_id",
			"player_avatar_id",
			"player_avatar_outfit_id",
			"avatar_slot_id",
			"quantity",
		},
		Values: []any{
			req.Record.OID,
			req.Record.Item.ID,
			req.Record.PlayerAvatar.ID,
			req.Record.PlayerOutfit.ID,
			req.Record.AvatarSlot.ID,
			req.Record.Quantity,
		},
		IDField: "id",
		IDValue: req.Record.ID,
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrPlayerItemExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeletePlayerItems(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "player.Service.DeletePlayerItems"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "player_item",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) GetGSFInventoryObjects(ctx context.Context, platform gsf.Platform) ([]types.PlayerItem, error) {
	const op = "player.Service.GetGSFInventoryObjects"

	rows, err := s.store.DB().QueryContext(ctx, `
			select
				pi.item_id,
				pi.gsfoid as player_item_gsfoid,
				pl.gsfoid as player_gsfoid,
				pi.quantity as quantity
			from player_item as pi
			join player as pl on pl.id = pi.player_id
			where pi.player_avatar_outfit_id is null;
		`)
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
		item, err := s.items.GetGSFItem(ctx, platform, itemID)
		if err != nil {
			return playerItems, wrap.IfErr(op, err)
		}
		playerItem.Item = item
		playerItems = append(playerItems, playerItem)
	}

	return playerItems, wrap.IfErr(op, rows.Err())
}

func (s *Service) GetGSFOutfitItems(ctx context.Context, platform gsf.Platform, playerAvatarOutfitOID, playerOID types.OID) ([]types.PlayerItem, error) {
	const op = "player.Service.GetGSFOutfitItems"

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
		item, err := s.items.GetGSFItem(ctx, platform, itemID)
		if err != nil {
			return playerItems, wrap.IfErr(op, err)
		}
		playerItem.Item = item
		playerItems = append(playerItems, playerItem)
	}

	return playerItems, wrap.IfErr(op, rows.Err())
}

func (s *Service) ReplaceGSFOutfitItems(ctx context.Context, oldOIDs, newOIDs []types.OID) error {
	const op = "player.Service.ReplaceGSFOutfitItems"

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
