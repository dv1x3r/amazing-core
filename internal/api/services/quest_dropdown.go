package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

func (s *Service) GetQuestsDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetQuestsDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "quest as q",
		IDField:      "q.id",
		TextField:    "(n.name || ' - ' || q.name)",
		OrderByField: "n.name, q.ordinal",
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("npc as n", "n.id = q.npc_id")
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetQuestStatesDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetQuestStatesDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "quest_state",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "id",
	})
	return res, wrap.IfErr(op, err)
}
