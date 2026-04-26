package randname

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
)

var (
	ErrNameNotFound = errors.New("name not found")
	ErrNameExists   = errors.New("name with the same type and name already exists")
)

type Service struct {
	logger *slog.Logger
	store  db.Store
}

func NewService(logger *slog.Logger, store db.Store) *Service {
	return &Service{
		logger: logger,
		store:  store,
	}
}

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

	var names []string

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

func (s *Service) GetRandomNameForm(ctx context.Context, req w2.GetFormRequest) (w2.GetFormResponse[RandomName], error) {
	const op = "randname.Service.GetRandomNameForm"
	res, err := w2db.GetFormContext(ctx, s.store.DB(), req, w2db.GetFormOptions[RandomName]{
		From:    "random_name",
		IDField: "id",
		Select:  []string{"id", "part_type", "name"},
		Scan: func(row *sql.Row, record *RandomName) error {
			return row.Scan(&record.ID, &record.PartType, &record.Name)
		},
	})
	if errors.Is(err, sql.ErrNoRows) {
		return res, wrap.IfErr(op, ErrNameNotFound)
	}
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateRandomName(ctx context.Context, req w2.SaveFormRequest[RandomName]) (int, error) {
	const op = "randname.Service.CreateRandomName"
	id, err := w2db.InsertFormContext(ctx, s.store.DB(), req, w2db.InsertFormOptions{
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
	affected, err := w2db.UpdateFormContext(ctx, s.store.DB(), req, w2db.UpdateFormOptions{
		Update:  "random_name",
		IDField: "id",
		Cols:    []string{"part_type", "name"},
		Values:  []any{req.Record.PartType, req.Record.Name},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrNameExists)
	} else if affected == 0 && err == nil {
		return wrap.IfErr(op, ErrNameNotFound)
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
