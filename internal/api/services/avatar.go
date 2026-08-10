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

var ErrAvatarExists = errors.New("avatar with the same name or container already exists")

type Avatar struct {
	ID        int              `json:"id"`
	Container w2.Dropdown      `json:"container"`
	Name      w2.Field[string] `json:"name"`
}

func (s *Service) GetAvatarGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Avatar], error) {
	const op = "api.Service.GetAvatarGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Avatar]{
		From: "avatar as av",
		Select: []string{
			"av.id",
			"av.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
			"av.name",
		},
		Where: map[string]string{
			"id":        "av.id",
			"container": "ac.gsfoid || ' - ' || ac.name",
			"name":      "av.name",
		},
		OrderBy: map[string]string{
			"id":        "av.id",
			"container": "ac.gsfoid",
			"name":      "av.name",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = av.container_id")
		},
		Scan: func(rows *sql.Rows, record *Avatar) error {
			return rows.Scan(
				&record.ID,
				&record.Container.ID,
				&record.Container.Text,
				&record.Name,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateAvatar(ctx context.Context, req w2.SaveFormRequest[Avatar]) (int, error) {
	const op = "api.Service.CreateAvatar"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "avatar",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
			"name":         req.Record.Name,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrAvatarExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateAvatar(ctx context.Context, req w2.SaveFormRequest[Avatar]) error {
	const op = "api.Service.UpdateAvatar"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "avatar",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
			"name":         req.Record.Name,
		},
		Where: map[string]any{
			"id": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrAvatarExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteAvatars(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteAvatars"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "avatar",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
