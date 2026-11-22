package services

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2db"
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

	type row struct {
		containerID int
		locked      bool
	}
	npcRows, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[row]{
		Query: `
			select
				n.container_id,
				n.is_locked
			from npc as n
			join zone as z on z.id = n.zone_id
			join asset_container as zc on zc.id = z.container_id
			where zc.gsfoid = ?
			order by n.id;
		`,
		Args: []any{zoneOID},
		Scan: func(rows *sql.Rows, r *row) error { return rows.Scan(&r.containerID, &r.locked) },
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	npcs := make([]types.NPC, 0, len(npcRows))
	for _, r := range npcRows {
		container, err := s.GetAssetContainer(ctx, platform, r.containerID)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		npc := types.NPC{ZoneOID: zoneOID, StartQuestOIDs: []types.OID{}}
		npc.AssetContainer = container
		npc.Locked = r.locked
		npcs = append(npcs, npc)
	}

	return npcs, nil
}

// GetNPCsWithQuestOffer returns the NPCs that offer a quest to the player now (GetNPCQuest), nil when none.
func (s *Service) GetNPCsWithQuestOffer(ctx context.Context, platform gsf.Platform, playerOID types.OID) ([]types.NPC, error) {
	const op = "game.Service.GetNPCsWithQuestOffer"

	type row struct {
		containerID int
		npcOID      types.OID
		zoneOID     types.OID
		locked      bool
	}
	npcRows, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[row]{
		Query: `
			select n.container_id, ac.gsfoid, zc.gsfoid, n.is_locked
			from npc as n
			join asset_container as ac on ac.id = n.container_id
			join zone as z on z.id = n.zone_id
			join asset_container as zc on zc.id = z.container_id
			order by n.id;
		`,
		Scan: func(rows *sql.Rows, r *row) error { return rows.Scan(&r.containerID, &r.npcOID, &r.zoneOID, &r.locked) },
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	var npcs []types.NPC
	for _, r := range npcRows {
		quest, err := s.GetNPCQuest(ctx, platform, playerOID, r.npcOID)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		if quest == nil {
			continue
		}
		container, err := s.GetAssetContainer(ctx, platform, r.containerID)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		npc := types.NPC{ZoneOID: r.zoneOID, StartQuestOIDs: []types.OID{}}
		npc.AssetContainer = container
		npc.Locked = r.locked
		npcs = append(npcs, npc)
	}
	return npcs, nil
}

// GetNPCRelationships returns the relationships of the player with the NPCs of the zone, with every NPC when zoneOID is nil.
// The level and the meter range are the client default of a new relationship (RewardSequenceAsyncManager),
// the points of the levels are unknown.
func (s *Service) GetNPCRelationships(ctx context.Context, playerOID types.OID, zoneOID *types.OID) ([]types.NPCRelationship, error) {
	const op = "game.Service.GetNPCRelationships"
	var zone any
	if zoneOID != nil {
		zone = zoneOID.Int64()
	}
	relationships, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[types.NPCRelationship]{
		Query: `
			select pn.created_at, pl.gsfoid, ac.gsfoid, pn.points, pn.is_friend
			from player_npc as pn
			join player as pl on pl.id = pn.player_id
			join npc as n on n.id = pn.npc_id
			join asset_container as ac on ac.id = n.container_id
			join zone as z on z.id = n.zone_id
			join asset_container as zc on zc.id = z.container_id
			where pl.gsfoid = ? and (? is null or zc.gsfoid = ?)
			order by pn.id;
		`,
		Args: []any{playerOID, zone, zone},
		Scan: func(rows *sql.Rows, r *types.NPCRelationship) error {
			r.Level, r.MinPoints, r.MaxPoints = 1, 0, 39
			return rows.Scan(&r.CreateDate, &r.PlayerOID, &r.NPCOID, &r.Points, &r.Friend)
		},
	})
	if relationships == nil {
		relationships = []types.NPCRelationship{}
	}
	return relationships, wrap.IfErr(op, err)
}
