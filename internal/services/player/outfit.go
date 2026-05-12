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

func (s *Service) GetPlayerOutfitGrid(ctx context.Context, req w2.GetGridRequest, playerID int) (w2.GetGridResponse[PlayerOutfit], error) {
	const op = "player.Service.GetPlayerOutfitGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerOutfit]{
		From: "player_avatar_outfit as po",
		Select: []string{
			"po.id",
			"po.gsfoid",
			"po.player_avatar_id",
			"(pa.gsfoid || ' - ' || pa.name) as player_avatar_name",
			"po.outfit_no",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
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
	const op = "player.Service.CreatePlayerOutfit"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "player_avatar_outfit",
		Cols: []string{"gsfoid", "player_avatar_id", "outfit_no"},
		Values: []any{
			sqlbuilder.Buildf("(select coalesce(max(gsfoid) + 1, 1) from player_avatar_outfit)"),
			req.Record.PlayerAvatar.ID,
			req.Record.OutfitNo,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrPlayerOutfitExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerOutfit(ctx context.Context, req w2.SaveFormRequest[PlayerOutfit]) error {
	const op = "player.Service.UpdatePlayerOutfit"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update:  "player_avatar_outfit",
		Cols:    []string{"gsfoid", "player_avatar_id", "outfit_no"},
		Values:  []any{req.Record.OID, req.Record.PlayerAvatar.ID, req.Record.OutfitNo},
		IDField: "id",
		IDValue: req.Record.ID,
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrPlayerOutfitExists)
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

func (s *Service) GetGSFOutfits(ctx context.Context, playerAvatarOID, playerOID types.OID) ([]types.PlayerAvatarOutfit, error) {
	const op = "player.Service.GetGSFOutfits"

	rows, err := s.store.DB().QueryContext(ctx, `
			select
				po.gsfoid,
				pl.gsfoid as player_gsfoid,
				pa.gsfoid as player_avatar_gsfoid,
				po.outfit_no
			from player_avatar_outfit as po
			join player_avatar as pa on pa.id = po.player_avatar_id
			join player as pl on pl.id = pa.player_id
			where pa.gsfoid = ? and pl.gsfoid = ?
			order by po.outfit_no;
		`, playerAvatarOID, playerOID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	outfits := []types.PlayerAvatarOutfit{}
	for rows.Next() {
		var outfit types.PlayerAvatarOutfit
		if err := rows.Scan(
			&outfit.OID,
			&outfit.PlayerOID,
			&outfit.PlayerAvatarOID,
			&outfit.OutfitNo,
		); err != nil {
			return outfits, wrap.IfErr(op, err)
		}
		outfits = append(outfits, outfit)
	}

	return outfits, wrap.IfErr(op, rows.Err())
}
