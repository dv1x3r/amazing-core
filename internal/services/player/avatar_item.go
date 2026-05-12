package player

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type PlayerAvatarItem struct {
	ID           int         `json:"id"`
	PlayerAvatar w2.Dropdown `json:"player_avatar"`
	PlayerItem   w2.Dropdown `json:"player_item"`
	AvatarSlot   w2.Dropdown `json:"avatar_slot"`
}

func (s *Service) GetPlayerAvatarItemGrid(ctx context.Context, req w2.GetGridRequest, playerID int) (w2.GetGridResponse[PlayerAvatarItem], error) {
	const op = "player.Service.GetPlayerAvatarItemGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerAvatarItem]{
		From: "player_avatar_item as pai",
		Select: []string{
			"pai.id",
			"pai.player_avatar_id",
			"(pa.gsfoid || ' - ' || pa.name) as player_avatar_name",
			"pai.player_item_id",
			"(pi.gsfoid || ' - ' || it.name || ' x' || pi.quantity) as player_item_name",
			"pai.avatar_slot_id",
			"avs.slot as avatar_slot",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("player_avatar as pa", "pa.id = pai.player_avatar_id")
			sb.Join("player_item as pi", "pi.id = pai.player_item_id")
			sb.Join("item as it", "it.id = pi.item_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "avatar_slot as avs", "avs.id = pai.avatar_slot_id")
			sb.Where(sb.EQ("pa.player_id", playerID))
			sb.OrderByAsc("pai.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerAvatarItem) error {
			return rows.Scan(
				&record.ID,
				&record.PlayerAvatar.ID,
				&record.PlayerAvatar.Text,
				&record.PlayerItem.ID,
				&record.PlayerItem.Text,
				&record.AvatarSlot.ID,
				&record.AvatarSlot.Text,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreatePlayerAvatarItem(ctx context.Context, req w2.SaveFormRequest[PlayerAvatarItem]) (int, error) {
	const op = "player.Service.CreatePlayerAvatarItem"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into:   "player_avatar_item",
		Cols:   []string{"player_avatar_id", "player_item_id", "avatar_slot_id"},
		Values: []any{req.Record.PlayerAvatar.ID, req.Record.PlayerItem.ID, req.Record.AvatarSlot.ID},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrPlayerItemAttached)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerAvatarItem(ctx context.Context, req w2.SaveFormRequest[PlayerAvatarItem]) error {
	const op = "player.Service.UpdatePlayerAvatarItem"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update:  "player_avatar_item",
		Cols:    []string{"player_avatar_id", "player_item_id", "avatar_slot_id"},
		Values:  []any{req.Record.PlayerAvatar.ID, req.Record.PlayerItem.ID, req.Record.AvatarSlot.ID},
		IDField: "id",
		IDValue: req.Record.ID,
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrPlayerItemAttached)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeletePlayerAvatarItems(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "player.Service.DeletePlayerAvatarItems"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "player_avatar_item",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
