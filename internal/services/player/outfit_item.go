package player

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type PlayerOutfitItem struct {
	ID                 int         `json:"id"`
	PlayerAvatarOutfit w2.Dropdown `json:"player_avatar_outfit"`
	PlayerItem         w2.Dropdown `json:"player_item"`
	AvatarSlot         w2.Dropdown `json:"avatar_slot"`
}

func (s *Service) GetPlayerOutfitItemGrid(ctx context.Context, req w2.GetGridRequest, playerID int) (w2.GetGridResponse[PlayerOutfitItem], error) {
	const op = "player.Service.GetPlayerOutfitItemGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerOutfitItem]{
		From: "player_avatar_outfit_item as paoi",
		Select: []string{
			"paoi.id",
			"paoi.player_avatar_outfit_id",
			"(pa.name || ' / Outfit ' || po.outfit_no) as player_avatar_outfit_name",
			"paoi.player_item_id",
			"(pi.gsfoid || ' - ' || it.name || ' x' || pi.quantity) as player_item_name",
			"paoi.avatar_slot_id",
			"avs.slot as avatar_slot",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("player_avatar_outfit as po", "po.id = paoi.player_avatar_outfit_id")
			sb.Join("player_avatar as pa", "pa.id = po.player_avatar_id")
			sb.Join("player_item as pi", "pi.id = paoi.player_item_id")
			sb.Join("item as it", "it.id = pi.item_id")
			sb.Join("avatar_slot as avs", "avs.id = paoi.avatar_slot_id")
			sb.Where(sb.EQ("pa.player_id", playerID))
			sb.OrderByAsc("paoi.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerOutfitItem) error {
			return rows.Scan(
				&record.ID,
				&record.PlayerAvatarOutfit.ID,
				&record.PlayerAvatarOutfit.Text,
				&record.PlayerItem.ID,
				&record.PlayerItem.Text,
				&record.AvatarSlot.ID,
				&record.AvatarSlot.Text,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreatePlayerOutfitItem(ctx context.Context, req w2.SaveFormRequest[PlayerOutfitItem]) (int, error) {
	const op = "player.Service.CreatePlayerOutfitItem"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into:   "player_avatar_outfit_item",
		Cols:   []string{"player_avatar_outfit_id", "player_item_id", "avatar_slot_id"},
		Values: []any{req.Record.PlayerAvatarOutfit.ID, req.Record.PlayerItem.ID, req.Record.AvatarSlot.ID},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrPlayerItemAttached)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerOutfitItem(ctx context.Context, req w2.SaveFormRequest[PlayerOutfitItem]) error {
	const op = "player.Service.UpdatePlayerOutfitItem"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update:  "player_avatar_outfit_item",
		Cols:    []string{"player_avatar_outfit_id", "player_item_id", "avatar_slot_id"},
		Values:  []any{req.Record.PlayerAvatarOutfit.ID, req.Record.PlayerItem.ID, req.Record.AvatarSlot.ID},
		IDField: "id",
		IDValue: req.Record.ID,
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrPlayerItemAttached)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeletePlayerOutfitItems(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "player.Service.DeletePlayerOutfitItems"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "player_avatar_outfit_item",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
