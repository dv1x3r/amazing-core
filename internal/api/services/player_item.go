package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/oidtype"
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

	Maze     string        `json:"maze"`
	Position string        `json:"position"`
	Rotation string        `json:"rotation"`
	Parent   string        `json:"parent"`
	Ordinal  w2.Field[int] `json:"ordinal"`
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
			"coalesce(iif(pm.player_item_id is null, 'Yard', 'Interior of ' || hit.name) || ' (' || pm.gsfoid || ')', '') as maze",
			"iif(pi.player_maze_id is null, '', pi.position_x || ', ' || pi.position_y || ', ' || pi.position_z) as position",
			"iif(pi.player_maze_id is null, '', coalesce(pi.rotation, '')) as rotation",
			"coalesce(ppi.gsfoid || ' - ' || pit.name, '') as parent",
			"pi.ordinal",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("item as it", "it.id = pi.item_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "player_maze as pm", "pm.id = pi.player_maze_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "player_item as hpi", "hpi.id = pm.player_item_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "item as hit", "hit.id = hpi.item_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "player_item as ppi", "ppi.id = pi.parent_player_item_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "item as pit", "pit.id = ppi.item_id")
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
				&record.Maze,
				&record.Position,
				&record.Rotation,
				&record.Parent,
				&record.Ordinal,
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
	var id int
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		var err error
		id, err = w2db.InsertContext(ctx, tx, w2db.InsertOptions{
			Into: "player_item",
			Values: map[string]any{
				"player_id":               playerID,
				"gsfoid":                  oidOrNext(req.Record.OID, "player_item", oidtype.UNTYPED),
				"item_id":                 req.Record.Item.ID,
				"player_avatar_id":        req.Record.PlayerAvatar.ID,
				"player_avatar_outfit_id": req.Record.PlayerOutfit.ID,
				"avatar_slot_id":          req.Record.AvatarSlot.ID,
				"quantity":                req.Record.Quantity,
			},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrPlayerItemExists
		} else if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
				insert into player_maze (gsfoid, player_id, player_item_id, container_id)
				select (
					select coalesce(max(pm.gsfoid) + 1, t.min_oid)
					from oid_type as t
					left join player_maze as pm on pm.gsfoid between t.min_oid and t.max_oid
					where t.id = ?
				), pi.player_id, pi.id, i.interior_container_id
				from player_item as pi
				join item as i on i.id = pi.item_id
				where pi.id = ? and i.interior_container_id is not null;
			`, oidtype.PLAYER_MAZE, id)
		return err
	})
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerItem(ctx context.Context, req w2.SaveFormRequest[PlayerItem]) error {
	const op = "api.Service.UpdatePlayerItem"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "player_item",
		Values: map[string]any{
			"gsfoid":                  oidOrNext(req.Record.OID, "player_item", oidtype.UNTYPED),
			"item_id":                 req.Record.Item.ID,
			"player_avatar_id":        req.Record.PlayerAvatar.ID,
			"player_avatar_outfit_id": req.Record.PlayerOutfit.ID,
			"avatar_slot_id":          req.Record.AvatarSlot.ID,
			"quantity":                req.Record.Quantity,
		},
		Where: map[string]any{
			"id": req.Record.ID,
		},
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
