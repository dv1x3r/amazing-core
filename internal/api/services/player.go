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

var (
	ErrPlayerNotFound = errors.New("player not found")
	ErrPlayerExists   = errors.New("player with the same oid already exists")
)

type PlayerList struct {
	ID       int    `json:"id"`
	OID      string `json:"oid"`
	OIDStr   string `json:"oid_str"`
	Username string `json:"username"`
}

type PlayerDetails struct {
	ID                  int         `json:"id"`
	OID                 string      `json:"oid"`
	Username            string      `json:"username"`
	CreatedAt           w2.UnixTime `json:"created_at"`
	IsTutorialCompleted bool        `json:"is_tutorial_completed"`
	IsChatAllowed       bool        `json:"is_chat_allowed"`
	IsQA                bool        `json:"is_qa"`
	MaxOutfits          int         `json:"max_outfits"`
	ActiveAvatar        w2.Dropdown `json:"active_avatar"`
}

func (s *Service) GetPlayerListGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[PlayerList], error) {
	const op = "api.Service.GetPlayerListGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerList]{
		From: "player as p",
		Select: []string{
			"p.id",
			"p.gsfoid",
			"p.username",
		},
		WhereMapping: map[string]string{
			"id":       "p.id",
			"oid":      "p.gsfoid",
			"username": "p.username",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.OrderByAsc("p.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerList) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.Username,
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
	const op = "api.Service.GetPlayerDetailsForm"
	res, err := w2db.GetFormContext(ctx, s.store.DB(), req, w2db.GetFormOptions[PlayerDetails]{
		From:    "player as p",
		IDField: "p.id",
		Select: []string{
			"p.id",
			"p.gsfoid",
			"p.username",
			"p.created_at",
			"p.is_tutorial_completed",
			"p.is_chat_allowed",
			"p.is_qa",
			"p.max_outfits",
			"pa.id as active_avatar_id",
			"pa.name as active_avatar_name",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
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
				&record.Username,
				&record.CreatedAt,
				&record.IsTutorialCompleted,
				&record.IsChatAllowed,
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
	const op = "api.Service.UpdatePlayerDetails"
	playerAvatarID := req.Record.ActiveAvatar.ID.V
	playerID := req.RecID
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		affected, err := w2db.UpdateContext(ctx, tx, w2db.UpdateOptions{
			Update: "player",
			Values: map[string]any{
				"gsfoid":                req.Record.OID,
				"username":              req.Record.Username,
				"is_tutorial_completed": req.Record.IsTutorialCompleted,
				"is_chat_allowed":       req.Record.IsChatAllowed,
				"is_qa":                 req.Record.IsQA,
				"max_outfits":           req.Record.MaxOutfits,
			},
			Where: map[string]any{
				"id": playerID,
			},
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
