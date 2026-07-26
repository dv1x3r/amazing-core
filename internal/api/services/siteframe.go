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

var (
	ErrSiteFrameExists = errors.New("site frame with the same type value already exists")
)

type SiteFrame struct {
	ID        int           `json:"id"`
	Container w2.Dropdown   `json:"container"`
	TypeValue w2.Field[int] `json:"type_value"`
}

func (s *Service) GetSiteFrameGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[SiteFrame], error) {
	const op = "api.Service.GetSiteFrameGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[SiteFrame]{
		From: "site_frame as sf",
		Select: []string{
			"sf.id",
			"sf.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
			"sf.type_value",
		},
		WhereMapping: map[string]string{
			"id":         "sf.id",
			"container":  "ac.gsfoid || ' - ' || ac.name",
			"type_value": "sf.type_value",
		},
		OrderByMapping: map[string]string{
			"id":         "sf.id",
			"container":  "ac.gsfoid",
			"type_value": "sf.type_value",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = sf.container_id")
		},
		Scan: func(rows *sql.Rows, record *SiteFrame) error {
			return rows.Scan(
				&record.ID,
				&record.Container.ID,
				&record.Container.Text,
				&record.TypeValue,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateSiteFrame(ctx context.Context, req w2.SaveFormRequest[SiteFrame]) (int, error) {
	const op = "api.Service.CreateSiteFrame"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into:   "site_frame",
		Cols:   []string{"container_id", "type_value"},
		Values: []any{req.Record.Container.ID, req.Record.TypeValue},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrSiteFrameExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateSiteFrame(ctx context.Context, req w2.SaveFormRequest[SiteFrame]) error {
	const op = "api.Service.UpdateSiteFrame"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update:  "site_frame",
		Cols:    []string{"container_id", "type_value"},
		Values:  []any{req.Record.Container.ID, req.Record.TypeValue},
		IDField: "id",
		IDValue: req.Record.ID,
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrSiteFrameExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteSiteFrames(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteSiteFrames"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "site_frame",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
