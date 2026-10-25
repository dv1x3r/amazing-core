package services

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/oidtype"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type PlayerQuest struct {
	ID          int              `json:"id"`
	OID         w2.Field[string] `json:"oid"`
	OIDStr      string           `json:"oid_str"`
	Quest       w2.Dropdown      `json:"quest"`
	State       w2.Dropdown      `json:"state"`
	CreatedAt   w2.UnixTime      `json:"created_at"`
	CompletedAt *w2.UnixTime     `json:"completed_at"`
}

func (s *Service) GetPlayerQuestGrid(ctx context.Context, req w2.GetGridRequest, playerID int) (w2.GetGridResponse[PlayerQuest], error) {
	const op = "api.Service.GetPlayerQuestGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerQuest]{
		From: "player_quest as pq",
		Select: []string{
			"pq.id",
			"pq.gsfoid",
			"pq.quest_id",
			"(n.name || ' - ' || q.name) as quest",
			"pq.quest_state_id",
			"qs.name as state",
			"pq.created_at",
			"pq.completed_at",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("quest as q", "q.id = pq.quest_id")
			sb.Join("npc as n", "n.id = q.npc_id")
			sb.Join("quest_state as qs", "qs.id = pq.quest_state_id")
			sb.Where(sb.EQ("pq.player_id", playerID))
			sb.OrderByDesc("pq.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerQuest) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.Quest.ID,
				&record.Quest.Text,
				&record.State.ID,
				&record.State.Text,
				&record.CreatedAt,
				&record.CompletedAt,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromString(record.OID.V).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreatePlayerQuest(ctx context.Context, req w2.SaveFormRequest[PlayerQuest], playerID int) (int, error) {
	const op = "api.Service.CreatePlayerQuest"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "player_quest",
		Values: map[string]any{
			"gsfoid":         oidOrNext(req.Record.OID, "player_quest", oidtype.PLAYER_QUEST),
			"player_id":      playerID,
			"quest_id":       req.Record.Quest.ID,
			"quest_state_id": req.Record.State.ID,
		},
	})
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerQuest(ctx context.Context, req w2.SaveFormRequest[PlayerQuest]) error {
	const op = "api.Service.UpdatePlayerQuest"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "player_quest",
		Values: map[string]any{
			"gsfoid":         oidOrNext(req.Record.OID, "player_quest", oidtype.PLAYER_QUEST),
			"quest_id":       req.Record.Quest.ID,
			"quest_state_id": req.Record.State.ID,
		},
		Where: map[string]any{
			"id": req.Record.ID,
		},
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DeletePlayerQuests(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeletePlayerQuests"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "player_quest",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
