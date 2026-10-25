package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/oidtype"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/queststate"
)

// GetNPCQuest returns the quest the NPC offers to the player, or nil if it offers nothing.
// The NPC offers the quest after the last completed one by ordinal and starts over after the last one,
// but nothing until the cooldown of the last completed quest is over.
func (s *Service) GetNPCQuest(ctx context.Context, platform gsf.Platform, playerOID, npcOID types.OID) (*types.Quest, error) {
	const op = "game.Service.GetNPCQuest"

	var containerID int
	var quest types.Quest
	err := s.store.DB().QueryRowContext(ctx, `
			with
			npc_quest as (
				select q.id, q.container_id, q.ordinal, q.cooldown
				from quest as q
				join npc as n on n.id = q.npc_id
				join asset_container as ac on ac.id = n.container_id
				where ac.gsfoid = ?
			),
			last_quest as (
				select nq.ordinal, pq.completed_at + nq.cooldown as available_at
				from player_quest as pq
				join player as pl on pl.id = pq.player_id
				join npc_quest as nq on nq.id = pq.quest_id
				where pl.gsfoid = ? and pq.completed_at is not null
				order by pq.completed_at desc, pq.id desc
				limit 1
			)
			select nq.container_id, nq.ordinal
			from npc_quest as nq
			left join last_quest as lq on true
			where coalesce(lq.available_at, 0) <= unixepoch()
			order by nq.ordinal <= coalesce(lq.ordinal, 0), nq.ordinal
			limit 1;
		`, npcOID, playerOID).Scan(&containerID, &quest.Ordinal)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	container, err := s.GetAssetContainer(ctx, platform, containerID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	quest.AssetContainer = container

	return &quest, nil
}

// CreatePlayerQuest creates an AVAILABLE player quest and returns its OID.
func (s *Service) CreatePlayerQuest(ctx context.Context, playerOID, questOID types.OID) (types.OID, error) {
	const op = "game.Service.CreatePlayerQuest"
	var playerQuestOID types.OID
	err := s.store.DB().QueryRowContext(ctx, `
			insert into player_quest (gsfoid, player_id, quest_id, quest_state_id)
			select
				(
					select coalesce(max(pq.gsfoid) + 1, t.min_oid)
					from oid_type as t
					left join player_quest as pq on pq.gsfoid between t.min_oid and t.max_oid
					where t.id = ?
				),
				pl.id,
				q.id,
				?
			from player as pl
			join quest as q
			join asset_container as ac on ac.id = q.container_id
			where pl.gsfoid = ? and ac.gsfoid = ?
			returning gsfoid;
		`, oidtype.PLAYER_QUEST, queststate.AVAILABLE, playerOID, questOID).Scan(&playerQuestOID)
	return playerQuestOID, wrap.IfErr(op, err)
}

// SetPlayerQuestState sets the state of the player quest.
func (s *Service) SetPlayerQuestState(ctx context.Context, playerOID, playerQuestOID types.OID, state queststate.QuestState) error {
	const op = "game.Service.SetPlayerQuestState"
	result, err := s.store.DB().ExecContext(ctx, `
			update player_quest set quest_state_id = ?
			where gsfoid = ? and player_id = (select id from player where gsfoid = ?);
		`, state, playerQuestOID, playerOID)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return wrap.IfErr(op, err)
	}
	if rowsAffected == 0 {
		return wrap.IfErr(op, fmt.Errorf("player quest %v not found", playerQuestOID))
	}
	return nil
}
