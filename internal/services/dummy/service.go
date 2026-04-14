package dummy

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
)

type DummyConfig map[string]string

type Service struct {
	store db.Store
}

func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetValue(ctx context.Context, key string) (string, error) {
	const op = "dummy.Service.GetValue"
	var value string
	row := s.store.DB().QueryRowContext(ctx, "select value from dummy_config where key = ?;", key)
	return value, wrap.IfErr(op, row.Scan(&value))
}

func (s *Service) GetForm(ctx context.Context) (w2.GetFormResponse[DummyConfig], error) {
	const op = "dummy.Service.GetForm"
	rows, err := s.store.DB().Query(`select key, value from dummy_config;`)
	if err != nil {
		return w2.GetFormResponse[DummyConfig]{}, wrap.IfErr(op, err)
	}
	defer rows.Close()

	cfg := DummyConfig{}

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return w2.GetFormResponse[DummyConfig]{}, wrap.IfErr(op, err)
		}
		cfg[key] = value
	}

	if err := rows.Err(); err != nil {
		return w2.GetFormResponse[DummyConfig]{}, wrap.IfErr(op, err)
	}

	return w2.NewGetFormResponse(cfg), nil
}

func (s *Service) UpdateForm(ctx context.Context, cfg DummyConfig) error {
	const op = "dummy.Service.UpdateForm"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	for key, value := range cfg {
		_, err := tx.ExecContext(ctx, "update dummy_config set value = ? where key = ?;", value, key)
		if err != nil {
			return wrap.IfErr(op, err)
		}
	}

	return wrap.IfErr(op, tx.Commit())
}
