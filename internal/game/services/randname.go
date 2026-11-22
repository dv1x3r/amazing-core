package services

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2db"
)

func (s *Service) GetNStringsByType(ctx context.Context, namePartType string, amount int) ([]string, error) {
	const op = "game.Service.GetNStringsByType"
	names, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[string]{
		Query: "select name from random_name where part_type = ? order by random() limit ?;",
		Args:  []any{namePartType, amount},
		Scan:  func(rows *sql.Rows, name *string) error { return rows.Scan(name) },
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	if names == nil {
		names = []string{}
	}
	return names, nil
}
