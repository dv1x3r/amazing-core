package services

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type QuestStep struct {
	ID        int              `json:"id"`
	Container w2.Dropdown      `json:"container"`
	Ordinal   int              `json:"ordinal"`
	Name      w2.Field[string] `json:"name"`
}

func (s *Service) GetQuestStepGrid(ctx context.Context, req w2.GetGridRequest, questID int) (w2.GetGridResponse[QuestStep], error) {
	const op = "api.Service.GetQuestStepGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[QuestStep]{
		From: "quest as q",
		Select: []string{
			"q.id",
			"q.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
			"q.ordinal",
			"q.name",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = q.container_id")
			sb.Where(sb.EQ("q.parent_id", questID))
			sb.OrderByAsc("q.ordinal")
		},
		Scan: func(rows *sql.Rows, record *QuestStep) error {
			return rows.Scan(
				&record.ID,
				&record.Container.ID,
				&record.Container.Text,
				&record.Ordinal,
				&record.Name,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateQuestStep(ctx context.Context, req w2.SaveFormRequest[QuestStep], questID int) (int, error) {
	const op = "api.Service.CreateQuestStep"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "quest",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
			"npc_id":       sqlbuilder.Buildf("(select npc_id from quest where id = %v)", questID),
			"parent_id":    questID,
			"ordinal":      sqlbuilder.Buildf("(select coalesce(max(ordinal) + 1, 1) from quest where npc_id = (select npc_id from quest where id = %v))", questID),
			"name":         req.Record.Name,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrQuestExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateQuestStep(ctx context.Context, req w2.SaveFormRequest[QuestStep]) error {
	const op = "api.Service.UpdateQuestStep"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "quest",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
			"name":         req.Record.Name,
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
