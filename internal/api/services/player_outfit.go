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

var ErrPlayerOutfitExists = errors.New("player outfit with the same avatar and outfit number already exists")

type PlayerOutfit struct {
	ID           int             `json:"id"`
	OID          w2.Field[int64] `json:"oid"`
	OIDStr       string          `json:"oid_str"`
	PlayerAvatar w2.Dropdown     `json:"player_avatar"`
	OutfitNo     w2.Field[int]   `json:"outfit_no"`
}

func (s *Service) GetPlayerOutfitGrid(ctx context.Context, req w2.GetGridRequest, playerID int) (w2.GetGridResponse[PlayerOutfit], error) {
	const op = "api.Service.GetPlayerOutfitGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerOutfit]{
		From: "player_avatar_outfit as po",
		Select: []string{
			"po.id",
			"po.gsfoid",
			"po.player_avatar_id",
			"(pa.gsfoid || ' - ' || pa.name) as player_avatar_name",
			"po.outfit_no",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("player_avatar as pa", "pa.id = po.player_avatar_id")
			sb.Where(sb.EQ("pa.player_id", playerID))
			sb.OrderByAsc("po.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerOutfit) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.PlayerAvatar.ID,
				&record.PlayerAvatar.Text,
				&record.OutfitNo,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromInt64(record.OID.V).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreatePlayerOutfit(ctx context.Context, req w2.SaveFormRequest[PlayerOutfit]) (int, error) {
	const op = "api.Service.CreatePlayerOutfit"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "player_avatar_outfit",
		Values: map[string]any{
			"gsfoid":           oidOrNext(req.Record.OID, "player_avatar_outfit", oidtype.UNTYPED),
			"player_avatar_id": req.Record.PlayerAvatar.ID,
			"outfit_no":        req.Record.OutfitNo,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrPlayerOutfitExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerOutfit(ctx context.Context, req w2.SaveFormRequest[PlayerOutfit]) error {
	const op = "api.Service.UpdatePlayerOutfit"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "player_avatar_outfit",
		Values: map[string]any{
			"gsfoid":           oidOrNext(req.Record.OID, "player_avatar_outfit", oidtype.UNTYPED),
			"player_avatar_id": req.Record.PlayerAvatar.ID,
			"outfit_no":        req.Record.OutfitNo,
		},
		Where: map[string]any{
			"id": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrPlayerOutfitExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeletePlayerOutfits(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeletePlayerOutfits"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "player_avatar_outfit",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
