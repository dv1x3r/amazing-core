package player

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type PlayerOutfit struct {
	ID           int             `json:"id"`
	OID          w2.Field[int64] `json:"oid"`
	OIDStr       string          `json:"oid_str"`
	PlayerAvatar w2.Dropdown     `json:"player_avatar"`
	OutfitNo     w2.Field[int]   `json:"outfit_no"`
}

func (s *Service) GetPlayerOutfitGrid(ctx context.Context, req w2.GetGridRequest, id int) (w2.GetGridResponse[PlayerOutfit], error) {
	const op = "player.Service.GetPlayerOutfitGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerOutfit]{
		From: "player_avatar_outfit as po",
		Select: []string{
			"po.id",
			"po.gsfoid",
			"po.player_avatar_id",
			"pa.name as player_avatar_name",
			"po.outfit_no",
		},
		WhereMapping: map[string]string{
			"id":            "po.id",
			"oid":           "po.gsfoid",
			"player_avatar": "pa.name",
			"outfit_no":     "po.outfit_no",
		},
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("player_avatar as pa", "pa.id = po.player_avatar_id")
			sb.Where(sb.EQ("pa.player_id", id))
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
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

func (s *Service) CreatePlayerOutfit(ctx context.Context, req w2.SaveFormRequest[PlayerOutfit]) error {
	const op = "player.Service.CreatePlayerOutfit"
	_, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "player_avatar_outfit",
		Cols: []string{"player_avatar_id", "outfit_no", "gsfoid"},
		Values: []any{
			req.Record.PlayerAvatar.ID,
			req.Record.OutfitNo,
			sqlbuilder.Buildf("(select coalesce(max(gsfoid) + 1, 1) from player_avatar_outfit)"),
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrPlayerOutfitExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerOutfits(ctx context.Context, req w2.SaveGridRequest[PlayerOutfit]) error {
	const op = "player.Service.UpdatePlayerOutfits"
	_, err := w2db.SaveGridContext(ctx, s.store.DB(), req, w2db.SaveGridOptions[PlayerOutfit]{
		BuildOptions: func(change PlayerOutfit) w2db.UpdateOptions {
			return w2db.UpdateOptions{
				Update:  "player_avatar_outfit",
				Cols:    []string{"gsfoid", "player_avatar_id", "outfit_no"},
				Values:  []any{change.OID, change.PlayerAvatar.ID, change.OutfitNo},
				IDField: "id",
				IDValue: change.ID,
			}
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrPlayerOutfitExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeletePlayerOutfits(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "player.Service.DeletePlayerOutfits"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "player_avatar_outfit",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
