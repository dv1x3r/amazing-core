package avatar

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
	"github.com/dv1x3r/w2go/w2sql"

	"github.com/huandu/go-sqlbuilder"
)

var (
	ErrAvatarExists = errors.New("avatar with the same name already exists")
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

type Avatar struct {
	ID         int              `json:"id"`
	Name       w2.Field[string] `json:"name"`
	MaxOutfits w2.Field[int]    `json:"max_outfits"`
	Container  w2.Dropdown      `json:"container"`
}

func (s *Service) GetAvatarsDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "avatar.Service.GetAvatarsDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "avatar",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "name",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetAvatarGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Avatar], error) {
	const op = "avatar.Service.GetAvatarGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Avatar]{
		From: "avatar as av",
		Select: []string{
			"av.id",
			"av.name",
			"av.max_outfits",
			"av.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
		},
		WhereMapping: map[string]string{
			"id":          "av.id",
			"name":        "av.name",
			"max_outfits": "av.max_outfits",
			"container":   "ac.gsfoid || ' - ' || ac.name",
		},
		OrderByMapping: map[string]string{
			"id":          "av.id",
			"name":        "av.name",
			"max_outfits": "av.max_outfits",
			"container":   "ac.gsfoid",
		},
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = av.container_id")
		},
		Scan: func(rows *sql.Rows, record *Avatar) error {
			return rows.Scan(
				&record.ID,
				&record.Name,
				&record.MaxOutfits,
				&record.Container.ID,
				&record.Container.Text,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateAvatar(ctx context.Context, req w2.SaveFormRequest[Avatar]) (int, error) {
	const op = "avatar.Service.CreateAvatar"
	id, err := w2db.InsertFormContext(ctx, s.store.DB(), req, w2db.InsertFormOptions{
		Into:   "avatar",
		Cols:   []string{"name", "max_outfits", "container_id"},
		Values: []any{req.Record.Name, req.Record.MaxOutfits, req.Record.Container.ID},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrAvatarExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateAvatars(ctx context.Context, req w2.SaveGridRequest[Avatar]) error {
	const op = "avatar.Service.UpdateAvatars"
	_, err := w2db.SaveGridContext(ctx, s.store.DB(), req, w2db.SaveGridOptions[Avatar]{
		BuildUpdate: func(change Avatar) *sqlbuilder.UpdateBuilder {
			ub := sqlbuilder.Update("avatar")
			w2sql.Set(ub, change.Name, "name")
			w2sql.Set(ub, change.MaxOutfits, "max_outfits")
			w2sql.Set(ub, change.Container.ID, "container_id")
			ub.Where(ub.EQ("id", change.ID))
			return ub
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrAvatarExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteAvatars(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "avatar.Service.DeleteAvatars"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "avatar",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
