package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

var ErrPlayerItemExists = errors.New("player item with the same oid already exists")

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
	const op = "api.Service.GetPlayerItemGrid"
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
	const op = "api.Service.CreatePlayerItem"
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
	const op = "api.Service.UpdatePlayerItem"
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
	const op = "api.Service.DeletePlayerItems"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "player_item",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
