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
