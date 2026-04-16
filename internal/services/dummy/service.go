package dummy

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type Service struct {
	store db.Store
}

func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

type Param struct {
	RowID int    `json:"rowid"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *Service) GetValue(ctx context.Context, key string) (string, error) {
	const op = "dummy.Service.GetValue"
	var value string
	row := s.store.DB().QueryRowContext(ctx, "select value from dummy_config where key = ?;", key)
	return value, wrap.IfErr(op, row.Scan(&value))
}

func (s *Service) GetDummyParametersGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Param], error) {
	const op = "dummy.Service.GetDummyParametersGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Param]{
		From:   "dummy_config",
		Select: []string{"rowid", "key", "value"},
		WhereMapping: map[string]string{
			"key":   "key",
			"value": "value",
		},
		Flavor: sqlbuilder.SQLite,
		Scan: func(rows *sql.Rows) (Param, error) {
			var record Param
			return record, rows.Scan(&record.RowID, &record.Key, &record.Value)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) UpdateDummyParameters(ctx context.Context, req w2.SaveGridRequest[Param]) error {
	const op = "dummy.Service.UpdateDummyParameters"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.SaveGridContext(ctx, tx, req, w2db.SaveGridOptions[Param]{
			Flavor: sqlbuilder.SQLite,
			BuildUpdate: func(change Param) *sqlbuilder.UpdateBuilder {
				ub := sqlbuilder.Update("dummy_config")
				ub.Set(ub.Assign("value", change.Value))
				ub.Where(ub.EQ("rowid", change.RowID))
				return ub
			},
		})
		return err
	})
	return wrap.IfErr(op, err)
}
