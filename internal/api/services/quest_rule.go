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

var ErrQuestRuleExists = errors.New("quest rule with the same item already exists")

type QuestRule struct {
	ID     int           `json:"id"`
	Item   w2.Dropdown   `json:"item"`
	Amount w2.Field[int] `json:"amount"`
}

func (s *Service) GetQuestRuleGrid(ctx context.Context, req w2.GetGridRequest, questID int) (w2.GetGridResponse[QuestRule], error) {
	const op = "api.Service.GetQuestRuleGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[QuestRule]{
		From: "quest_rule as qr",
		Select: []string{
			"qr.id",
			"qr.item_container_id",
			"(ac.gsfoid || ' - ' || ac.name) as item",
			"qr.amount",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = qr.item_container_id")
			sb.Where(sb.EQ("qr.quest_id", questID))
			sb.OrderByAsc("qr.id")
		},
		Scan: func(rows *sql.Rows, record *QuestRule) error {
			return rows.Scan(
				&record.ID,
				&record.Item.ID,
				&record.Item.Text,
				&record.Amount,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateQuestRule(ctx context.Context, req w2.SaveFormRequest[QuestRule], questID int) (int, error) {
	const op = "api.Service.CreateQuestRule"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "quest_rule",
		Values: map[string]any{
			"quest_id":          questID,
			"item_container_id": req.Record.Item.ID,
			"amount":            req.Record.Amount,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrQuestRuleExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateQuestRule(ctx context.Context, req w2.SaveFormRequest[QuestRule]) error {
	const op = "api.Service.UpdateQuestRule"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "quest_rule",
		Values: map[string]any{
			"item_container_id": req.Record.Item.ID,
			"amount":            req.Record.Amount,
		},
		Where: map[string]any{
			"id": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrQuestRuleExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteQuestRules(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteQuestRules"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "quest_rule",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
