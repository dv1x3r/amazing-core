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

var ErrZoneExists = errors.New("zone with the same container already exists")

type Zone struct {
	ID        int         `json:"id"`
	Container w2.Dropdown `json:"container"`
}

func (s *Service) GetZoneGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Zone], error) {
	const op = "api.Service.GetZoneGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Zone]{
		From: "zone as z",
		Select: []string{
			"z.id",
			"z.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
		},
		WhereMapping: map[string]string{
			"id":        "z.id",
			"container": "ac.gsfoid || ' - ' || ac.name",
		},
		OrderByMapping: map[string]string{
			"id":        "z.id",
			"container": "ac.gsfoid",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = z.container_id")
		},
		Scan: func(rows *sql.Rows, record *Zone) error {
			return rows.Scan(
				&record.ID,
				&record.Container.ID,
				&record.Container.Text,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateZone(ctx context.Context, req w2.SaveFormRequest[Zone]) (int, error) {
	const op = "api.Service.CreateZone"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "zone",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrZoneExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateZone(ctx context.Context, req w2.SaveFormRequest[Zone]) error {
	const op = "api.Service.UpdateZone"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "zone",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
		},
		Where: map[string]any{
			"id": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrZoneExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteZones(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteZones"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "zone",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
