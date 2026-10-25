package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetZones(ctx context.Context, platform gsf.Platform) ([]types.Zone, error) {
	const op = "game.Service.GetZones"

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
		container, err := s.GetAssetContainer(ctx, platform, containerID)
		if err != nil {
			return zones, wrap.IfErr(op, err)
		}
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
