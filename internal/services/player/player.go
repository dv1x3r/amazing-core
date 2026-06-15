package player

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type PlayerList struct {
	ID         int    `json:"id"`
	OID        string `json:"oid"`
	OIDStr     string `json:"oid_str"`
	ActiveName string `json:"active_name"`
}

type PlayerDetails struct {
	ID                  int         `json:"id"`
	OID                 string      `json:"oid"`
	CreatedAt           w2.UnixTime `json:"created_at"`
	IsTutorialCompleted bool        `json:"is_tutorial_completed"`
	IsQA                bool        `json:"is_qa"`
	MaxOutfits          int         `json:"max_outfits"`
	ActiveAvatar        w2.Dropdown `json:"active_avatar"`
}

func (s *Service) GetPlayerListGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[PlayerList], error) {
	const op = "player.Service.GetPlayerListGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerList]{
		From: "player as p",
		Select: []string{
			"p.id",
			"p.gsfoid",
			"coalesce(pa.name, '[None]') as active_name",
		},
		WhereMapping: map[string]string{
			"id":          "p.id",
			"oid":         "p.gsfoid",
			"active_name": "pa.name",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin, `(
					select
						id,
						player_id,
						(gsfoid || ' - ' || name) as name,
						row_number() over (partition by player_id order by id) as rn
					from player_avatar
					where is_active = 1
				) as pa`,
				"pa.player_id = p.id and pa.rn = 1",
			)
			sb.OrderByAsc("p.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerList) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.ActiveName,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromString(record.OID).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetPlayerDetailsForm(ctx context.Context, req w2.GetFormRequest) (w2.GetFormResponse[PlayerDetails], error) {
	const op = "player.Service.GetPlayerDetailsForm"
	res, err := w2db.GetFormContext(ctx, s.store.DB(), req, w2db.GetFormOptions[PlayerDetails]{
		From:    "player as p",
		IDField: "p.id",
		Select: []string{
			"p.id",
			"p.gsfoid",
			"p.created_at",
			"p.is_tutorial_completed",
			"p.is_qa",
			"p.max_outfits",
			"pa.id as active_avatar_id",
			"pa.name as active_avatar_name",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin, `(
					select
						id,
						player_id,
						(gsfoid || ' - ' || name) as name,
						row_number() over (partition by player_id order by id) as rn
					from player_avatar
					where is_active = 1
				) as pa`,
				"pa.player_id = p.id and pa.rn = 1",
			)
		},
		Scan: func(row *sql.Row, record *PlayerDetails) error {
			return row.Scan(
				&record.ID,
				&record.OID,
				&record.CreatedAt,
				&record.IsTutorialCompleted,
				&record.IsQA,
				&record.MaxOutfits,
				&record.ActiveAvatar.ID,
				&record.ActiveAvatar.Text,
			)
		},
	})
	if errors.Is(err, sql.ErrNoRows) {
		return res, wrap.IfErr(op, ErrPlayerNotFound)
	}
	return res, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerDetails(ctx context.Context, req w2.SaveFormRequest[PlayerDetails]) error {
	const op = "player.Service.UpdatePlayerDetails"
	playerAvatarID := req.Record.ActiveAvatar.ID.V
	playerID := req.RecID
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		affected, err := w2db.UpdateContext(ctx, tx, w2db.UpdateOptions{
			Update:  "player",
			Cols:    []string{"gsfoid", "is_tutorial_completed", "is_qa", "max_outfits"},
			Values:  []any{req.Record.OID, req.Record.IsTutorialCompleted, req.Record.IsQA, req.Record.MaxOutfits},
			IDField: "id",
			IDValue: playerID,
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrPlayerExists
		} else if affected == 0 && err == nil {
			return ErrPlayerNotFound
		} else if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			update player_avatar
			set is_active = case when id = ? then 1 else 0 end
			where player_id = ?;
		`, playerAvatarID, playerID)
		return err
	})
	return wrap.IfErr(op, err)
}

func (s *Service) GetGSFPlayer(ctx context.Context, platform gsf.Platform, playerID int) (types.Player, error) {
	const op = "player.Service.GetGSFPlayer"

	row := s.store.DB().QueryRowContext(ctx, `
			select
				pl.gsfoid,
				pl.created_at,
				pl.is_tutorial_completed,
				pl.is_qa
			from player as pl
			where pl.id = ?;
		`, playerID)

	var player types.Player

	if err := row.Scan(
		&player.OID,
		&player.CreateDate,
		&player.IsTutorialCompleted,
		&player.IsQA,
	); err != nil {
		return player, wrap.IfErr(op, err)
	}

	activeAvatar, err := s.GetGSFActivePlayerAvatar(ctx, platform, playerID)
	if err != nil {
		return player, wrap.IfErr(op, err)
	}
	player.ActivePlayerAvatar = activeAvatar

	return player, nil
}
