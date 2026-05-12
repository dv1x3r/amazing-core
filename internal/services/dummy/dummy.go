package dummy

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
)

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
		Scan: func(rows *sql.Rows, record *Param) error {
			return rows.Scan(&record.RowID, &record.Key, &record.Value)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) UpdateDummyParameters(ctx context.Context, req w2.SaveGridRequest[Param]) error {
	const op = "dummy.Service.UpdateDummyParameters"
	_, err := w2db.SaveGridContext(ctx, s.store.DB(), req, w2db.SaveGridOptions[Param]{
		BuildOptions: func(change Param) w2db.UpdateOptions {
			return w2db.UpdateOptions{
				Update:  "dummy_config",
				Cols:    []string{"value"},
				Values:  []any{change.Value},
				IDField: "rowid",
				IDValue: change.RowID,
			}
		},
	})
	return wrap.IfErr(op, err)
}
