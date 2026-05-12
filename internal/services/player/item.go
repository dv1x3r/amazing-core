package player

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type PlayerItem struct {
	ID       int              `json:"id"`
	OID      w2.Field[string] `json:"oid"`
	OIDStr   string           `json:"oid_str"`
	Item     w2.Dropdown      `json:"item"`
	Quantity w2.Field[int]    `json:"quantity"`
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
			"pi.quantity",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("item as it", "it.id = pi.item_id")
			sb.Where(sb.EQ("pi.player_id", playerID))
			sb.OrderByAsc("pi.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerItem) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.Item.ID,
				&record.Item.Text,
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
		Cols: []string{"player_id", "gsfoid", "item_id", "quantity"},
		Values: []any{
			playerID,
			sqlbuilder.Buildf("(select coalesce(max(gsfoid) + 1, 1) from player_item)"),
			req.Record.Item.ID,
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
		Update:  "player_item",
		Cols:    []string{"gsfoid", "item_id", "quantity"},
		Values:  []any{req.Record.OID, req.Record.Item.ID, req.Record.Quantity},
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

func (s *Service) GetGSFAvatarItems(ctx context.Context, platform gsf.Platform, playerAvatarOID, playerOID types.OID) ([]types.PlayerItem, error) {
	const op = "player.Service.GetGSFAvatarItems"

	rows, err := s.store.DB().QueryContext(ctx, `
			select
				pi.gsfoid as player_item_gsfoid,
				pa.gsfoid as player_avatar_gsfoid,
				pl.gsfoid as player_gsfoid,
				avs.gsfoid as slot_gsfoid,
				pi.item_id
			from player_avatar_item as pai
			join player_item as pi on pi.id = pai.player_item_id
			join player_avatar as pa on pa.id = pai.player_avatar_id
			join player as pl on pl.id = pa.player_id
			left join avatar_slot as avs on avs.id = pai.avatar_slot_id
			where pa.gsfoid = ? and pl.gsfoid = ?;
		`, playerAvatarOID, playerOID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	playerItems := []types.PlayerItem{}
	for rows.Next() {
		var playerItem types.PlayerItem
		var itemID int
		if err := rows.Scan(
			&playerItem.OID,
			&playerItem.PlayerAvatarOID,
			&playerItem.PlayerOID,
			&playerItem.SlotOID,
			&itemID,
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
				pi.gsfoid as player_item_gsfoid,
				pao.gsfoid as player_avatar_outfit_gsfoid,
				pa.gsfoid as player_avatar_gsfoid,
				pl.gsfoid as player_gsfoid,
				avs.gsfoid as slot_gsfoid,
				pi.quantity as player_item_quantity,
				pi.item_id
			from player_avatar_outfit_item as paoi
			join player_item as pi on pi.id = paoi.player_item_id
			join player_avatar_outfit as pao on pao.id = paoi.player_avatar_outfit_id
			join player_avatar as pa on pa.id = pao.player_avatar_id
			join player as pl on pl.id = pa.player_id
			left join avatar_slot as avs on avs.id = paoi.avatar_slot_id
			where pao.gsfoid = ? and pl.gsfoid = ?;
		`, playerAvatarOutfitOID, playerOID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	playerItems := []types.PlayerItem{}
	for rows.Next() {
		var playerItem types.PlayerItem
		var itemID int
		if err := rows.Scan(
			&playerItem.OID,
			&playerItem.PlayerAvatarOutfitOID,
			&playerItem.PlayerAvatarOID,
			&playerItem.PlayerOID,
			&playerItem.SlotOID,
			&playerItem.Quantity,
			&itemID,
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
