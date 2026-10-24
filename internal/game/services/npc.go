package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) IsNPC(ctx context.Context, oid types.OID) (bool, error) {
	const op = "game.Service.IsNPC"
	var exists bool
	err := s.store.DB().QueryRowContext(ctx, `
			select exists (
				select 1
				from npc as n
				join asset_container as ac on ac.id = n.container_id
				where ac.gsfoid = ?
			);
		`, oid).Scan(&exists)
	return exists, wrap.IfErr(op, err)
}

// GetZoneNPCs returns the NPCs placed in the zone.
// The client reads only zoneId, locked and the asset map, relationships come from GET_NPC_RELATIONSHIPS.
func (s *Service) GetZoneNPCs(ctx context.Context, platform gsf.Platform, zoneOID types.OID) ([]types.NPC, error) {
	const op = "game.Service.GetZoneNPCs"

	rows, err := s.store.DB().QueryContext(ctx, `
		select
			n.container_id,
			n.is_locked
		from npc as n
		join zone as z on z.id = n.zone_id
		join asset_container as zc on zc.id = z.container_id
		where zc.gsfoid = ?
		order by n.id;
	`, zoneOID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	containerIDs := []int{}
	npcs := []types.NPC{}
	for rows.Next() {
		var containerID int
		npc := types.NPC{ZoneOID: zoneOID, StartQuestOIDs: []types.OID{}}
		if err := rows.Scan(&containerID, &npc.Locked); err != nil {
			return nil, wrap.IfErr(op, err)
		}
		containerIDs = append(containerIDs, containerID)
		npcs = append(npcs, npc)
	}
	if err := rows.Err(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	for i, containerID := range containerIDs {
		container, err := s.GetAssetContainer(ctx, platform, containerID)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		npcs[i].AssetContainer = container
	}

	return npcs, nil
}
