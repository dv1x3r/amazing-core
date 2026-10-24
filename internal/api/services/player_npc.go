package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

var ErrPlayerNPCExists = errors.New("player already has this npc")

type PlayerNPC struct {
	ID       int            `json:"id"`
	NPC      w2.Dropdown    `json:"npc"`
	Points   w2.Field[int]  `json:"points"`
	IsFriend w2.Field[bool] `json:"is_friend"`
}

func (s *Service) GetPlayerNPCGrid(ctx context.Context, req w2.GetGridRequest, playerID int) (w2.GetGridResponse[PlayerNPC], error) {
	const op = "api.Service.GetPlayerNPCGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerNPC]{
		From: "player_npc as pn",
		Select: []string{
			"pn.id",
			"pn.npc_id",
			"n.name as npc_name",
			"pn.points",
			"pn.is_friend",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("npc as n", "n.id = pn.npc_id")
			sb.Where(sb.EQ("pn.player_id", playerID))
			sb.OrderByAsc("pn.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerNPC) error {
			return rows.Scan(
				&record.ID,
				&record.NPC.ID,
				&record.NPC.Text,
				&record.Points,
				&record.IsFriend,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreatePlayerNPC(ctx context.Context, req w2.SaveFormRequest[PlayerNPC], playerID int) (int, error) {
	const op = "api.Service.CreatePlayerNPC"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "player_npc",
		Values: map[string]any{
			"player_id": playerID,
			"npc_id":    req.Record.NPC.ID,
			"points":    req.Record.Points,
			"is_friend": req.Record.IsFriend,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrPlayerNPCExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerNPC(ctx context.Context, req w2.SaveFormRequest[PlayerNPC]) error {
	const op = "api.Service.UpdatePlayerNPC"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "player_npc",
		Values: map[string]any{
			"npc_id":    req.Record.NPC.ID,
			"points":    req.Record.Points,
			"is_friend": req.Record.IsFriend,
		},
		Where: map[string]any{
			"id": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrPlayerNPCExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeletePlayerNPCs(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeletePlayerNPCs"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "player_npc",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
