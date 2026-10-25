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

var ErrQuestExists = errors.New("quest with the same container already exists")

type Quest struct {
	ID        int              `json:"id"`
	Container w2.Dropdown      `json:"container"`
	Ordinal   int              `json:"ordinal"`
	Name      w2.Field[string] `json:"name"`
	XP        w2.Field[int]    `json:"xp"`
	Coins     w2.Field[int]    `json:"coins"`
	NPCPoints w2.Field[int]    `json:"npc_points"`
	Cooldown  w2.Field[int]    `json:"cooldown"`
}

func (s *Service) GetQuestGrid(ctx context.Context, req w2.GetGridRequest, npcID int) (w2.GetGridResponse[Quest], error) {
	const op = "api.Service.GetQuestGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Quest]{
		From: "quest as q",
		Select: []string{
			"q.id",
			"q.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
			"q.ordinal",
			"q.name",
			"q.xp",
			"q.coins",
			"q.npc_points",
			"q.cooldown",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = q.container_id")
			sb.Where(sb.EQ("q.npc_id", npcID))
			sb.OrderByAsc("q.ordinal")
		},
		Scan: func(rows *sql.Rows, record *Quest) error {
			return rows.Scan(
				&record.ID,
				&record.Container.ID,
				&record.Container.Text,
				&record.Ordinal,
				&record.Name,
				&record.XP,
				&record.Coins,
				&record.NPCPoints,
				&record.Cooldown,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateQuest(ctx context.Context, req w2.SaveFormRequest[Quest], npcID int) (int, error) {
	const op = "api.Service.CreateQuest"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "quest",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
			"npc_id":       npcID,
			"ordinal":      sqlbuilder.Buildf("(select coalesce(max(ordinal) + 1, 1) from quest where npc_id = %v)", npcID),
			"name":         req.Record.Name,
			"xp":           req.Record.XP.NotNull(),
			"coins":        req.Record.Coins.NotNull(),
			"npc_points":   req.Record.NPCPoints.NotNull(),
			"cooldown":     req.Record.Cooldown.NotNull(),
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrQuestExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateQuest(ctx context.Context, req w2.SaveFormRequest[Quest]) error {
	const op = "api.Service.UpdateQuest"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "quest",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
			"name":         req.Record.Name,
			"xp":           req.Record.XP.NotNull(),
			"coins":        req.Record.Coins.NotNull(),
			"npc_points":   req.Record.NPCPoints.NotNull(),
			"cooldown":     req.Record.Cooldown.NotNull(),
		},
		Where: map[string]any{
			"id": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrQuestExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteQuests(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteQuests"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "quest",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) ReorderQuests(ctx context.Context, req w2.ReorderGridRequest) error {
	const op = "api.Service.ReorderQuests"
	_, err := w2db.ReorderGridContext(ctx, s.store.DB(), req, w2db.ReorderGridOptions{
		Update:     "quest",
		IDField:    "id",
		SetField:   "ordinal",
		GroupField: "npc_id",
	})
	return wrap.IfErr(op, err)
}
