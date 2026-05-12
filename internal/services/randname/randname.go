package randname

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
)

type RandomName struct {
	ID       int    `json:"id"`
	PartType string `json:"part_type"`
	Name     string `json:"name"`
}

func (s *Service) GetNStringsByType(ctx context.Context, namePartType string, amount int) ([]string, error) {
	const op = "randname.Service.GetNStringsByType"

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

func (s *Service) GetRandomNameGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[RandomName], error) {
	const op = "randname.Service.GetRandomNameGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[RandomName]{
		From:   "random_name",
		Select: []string{"id", "part_type", "name"},
		WhereMapping: map[string]string{
			"id":        "id",
			"part_type": "part_type",
			"name":      "name",
		},
		OrderByMapping: map[string]string{
			"id":        "id",
			"part_type": "part_type",
			"name":      "name",
		},
		Scan: func(rows *sql.Rows, record *RandomName) error {
			return rows.Scan(&record.ID, &record.PartType, &record.Name)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateRandomName(ctx context.Context, req w2.SaveFormRequest[RandomName]) (int, error) {
	const op = "randname.Service.CreateRandomName"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into:   "random_name",
		Cols:   []string{"part_type", "name"},
		Values: []any{req.Record.PartType, req.Record.Name},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrNameExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateRandomName(ctx context.Context, req w2.SaveFormRequest[RandomName]) error {
	const op = "randname.Service.UpdateRandomName"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update:  "random_name",
		Cols:    []string{"part_type", "name"},
		Values:  []any{req.Record.PartType, req.Record.Name},
		IDField: "id",
		IDValue: req.Record.ID,
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrNameExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteRandomNames(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "randname.Service.DeleteRandomNames"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "random_name",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
