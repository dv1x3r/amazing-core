package zone

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type Zone struct {
	ID        int         `json:"id"`
	Container w2.Dropdown `json:"container"`
}

func (s *Service) GetZoneGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Zone], error) {
	const op = "zone.Service.GetZoneGrid"
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
	const op = "zone.Service.CreateZone"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into:   "zone",
		Cols:   []string{"container_id"},
		Values: []any{req.Record.Container.ID},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrZoneExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateZone(ctx context.Context, req w2.SaveFormRequest[Zone]) error {
	const op = "zone.Service.UpdateZone"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update:  "zone",
		Cols:    []string{"container_id"},
		Values:  []any{req.Record.Container.ID},
		IDField: "id",
		IDValue: req.Record.ID,
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrZoneExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteZones(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "zone.Service.DeleteZones"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "zone",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) GetGSFZones(ctx context.Context, platform gsf.Platform) ([]types.Zone, error) {
	const op = "zone.Service.GetGSFZones"

	rows, err := s.store.DB().QueryContext(ctx, "select container_id from zone;")
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	zones := []types.Zone{}
	for rows.Next() {
		var zone types.Zone
		var containerID int
		if err := rows.Scan(&containerID); err != nil {
			return zones, wrap.IfErr(op, err)
		}
		container, err := s.assets.GetGSFAssetContainer(ctx, platform, containerID)
		if err != nil {
			return zones, wrap.IfErr(op, err)
		}
		zone.AssetContainer = container
		zones = append(zones, zone)
	}

	return zones, nil
}
