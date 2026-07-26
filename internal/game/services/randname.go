package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
)

func (s *Service) GetNStringsByType(ctx context.Context, namePartType string, amount int) ([]string, error) {
	const op = "game.Service.GetNStringsByType"

	const query = "select name from random_name where part_type = ? order by random() limit ?;"
	rows, err := s.store.DB().QueryContext(ctx, query, namePartType, amount)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	names := make([]string, 0, amount)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, wrap.IfErr(op, err)
		}
		names = append(names, name)
	}

	return names, wrap.IfErr(op, rows.Err())
}
