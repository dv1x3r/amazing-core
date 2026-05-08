package item

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
	"github.com/huandu/go-sqlbuilder"
)

type Item struct {
	ID         int              `json:"id"`
	Name       w2.Field[string] `json:"name"`
	Container  w2.Dropdown      `json:"container"`
	Categories []w2.Dropdown    `json:"categories"`
	Slots      []w2.Dropdown    `json:"slots"`
}

func (s *Service) GetItemGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Item], error) {
	const op = "item.Service.GetItemGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Item]{
		From: "item as it",
		Select: []string{
			"it.id",
			"it.name",
			"it.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
		},
		WhereMapping: map[string]string{
			"id":        "it.id",
			"name":      "it.name",
			"container": "ac.gsfoid || ' - ' || ac.name",
		},
		OrderByMapping: map[string]string{
			"id":        "it.id",
			"name":      "it.name",
			"container": "ac.gsfoid",
		},
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = it.container_id")
		},
		Scan: func(rows *sql.Rows, record *Item) error {
			return rows.Scan(
				&record.ID,
				&record.Name,
				&record.Container.ID,
				&record.Container.Text,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateItem(ctx context.Context, req w2.SaveFormRequest[Item]) (int, error) {
	const op = "item.Service.CreateItem"
	// id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
	// 	Into:   "site_frame",
	// 	Cols:   []string{"type_value", "container_id"},
	// 	Values: []any{req.Record.TypeValue, req.Record.Container.ID},
	// })
	// if s.store.IsErrConstraintUnique(err) {
	// 	return 0, wrap.IfErr(op, ErrSiteFrameExists)
	// }
	// return id, wrap.IfErr(op, err)
	s.logger.Debug("item", "req", req)
	return 0, nil
}
