package services

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2db"
)

func (s *Service) GetZones(ctx context.Context, platform gsf.Platform) ([]types.Zone, error) {
	const op = "game.Service.GetZones"

	containerIDs, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[int]{
		Query: "select container_id from zone;",
		Scan:  func(rows *sql.Rows, containerID *int) error { return rows.Scan(containerID) },
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	zones := make([]types.Zone, 0, len(containerIDs))
	for _, containerID := range containerIDs {
		container, err := s.GetAssetContainer(ctx, platform, containerID)
		if err != nil {
			return zones, wrap.IfErr(op, err)
		}
		var zone types.Zone
		zone.AssetContainer = container
		zones = append(zones, zone)
	}

	return zones, nil
}

// GetZone returns the zone by the OID of its container.
func (s *Service) GetZone(ctx context.Context, platform gsf.Platform, zoneOID types.OID) (types.Zone, error) {
	const op = "game.Service.GetZone"

	var zone types.Zone
	var containerID int
	err := s.store.DB().QueryRowContext(ctx, `
			select z.container_id
			from zone as z
			join asset_container as ac on ac.id = z.container_id
			where ac.gsfoid = ?;
		`, zoneOID).Scan(&containerID)
	if err != nil {
		return zone, wrap.IfErr(op, err)
	}

	container, err := s.GetAssetContainer(ctx, platform, containerID)
	if err != nil {
		return zone, wrap.IfErr(op, err)
	}
	zone.AssetContainer = container

	return zone, nil
}
