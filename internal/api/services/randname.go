package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
)

var ErrNameExists = errors.New("name with the same type and name already exists")

type RandomName struct {
	ID       int    `json:"id"`
	PartType string `json:"part_type"`
	Name     string `json:"name"`
}

func (s *Service) GetRandomNameGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[RandomName], error) {
	const op = "api.Service.GetRandomNameGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[RandomName]{
		From:   "random_name",
		Select: []string{"id", "part_type", "name"},
		Where: map[string]string{
			"id":        "id",
			"part_type": "part_type",
			"name":      "name",
		},
		OrderBy: map[string]string{
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
	const op = "api.Service.CreateRandomName"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "random_name",
		Values: map[string]any{
			"part_type": req.Record.PartType,
			"name":      req.Record.Name,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrNameExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateRandomName(ctx context.Context, req w2.SaveFormRequest[RandomName]) error {
	const op = "api.Service.UpdateRandomName"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "random_name",
		Values: map[string]any{
			"part_type": req.Record.PartType,
			"name":      req.Record.Name,
		},
		Where: map[string]any{
			"id": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrNameExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteRandomNames(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteRandomNames"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "random_name",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
