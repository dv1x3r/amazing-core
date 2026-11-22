package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/oidtype"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/queststate"
	"github.com/dv1x3r/w2go/w2db"
)

// ErrQuestNotFound is a quest or player quest the server does not have,
// the client also sends the quests of the systems that are not implemented, like fishing and mob quests.
var ErrQuestNotFound = errors.New("quest not found")

// GetNPCQuest returns the quest the NPC offers to the player with its steps, or nil.
// Story quests go first by ordinal, each until the player completes it. Then the jobs go in a loop.
// Nothing is offered during the cooldown of the last completed quest.
func (s *Service) GetNPCQuest(ctx context.Context, platform gsf.Platform, playerOID, npcOID types.OID) (*types.Quest, error) {
	const op = "game.Service.GetNPCQuest"

	var questID int
	var containerID int
	var quest types.Quest
	err := s.store.DB().QueryRowContext(ctx, `
			with
			npc_quest as (
				select q.id, q.container_id, q.ordinal, q.cooldown, q.repeatable
				from quest as q
				join npc as n on n.id = q.npc_id
				join asset_container as ac on ac.id = n.container_id
				where ac.gsfoid = ? and q.parent_id is null
			),
			closed_quest as (
				select pq.id, pq.quest_id, pq.completed_at
				from player_quest as pq
				join player as pl on pl.id = pq.player_id
				where pl.gsfoid = ? and pq.completed_at is not null
			),
			last_quest as (
				select nq.ordinal, nq.repeatable, cq.completed_at + nq.cooldown as available_at
				from closed_quest as cq
				join npc_quest as nq on nq.id = cq.quest_id
				order by cq.completed_at desc, cq.id desc
				limit 1
			)
			select nq.id, nq.container_id, nq.ordinal
			from npc_quest as nq
			left join last_quest as lq on true
			where coalesce(lq.available_at, 0) <= unixepoch()
				and (nq.repeatable or nq.id not in (select quest_id from closed_quest))
			order by nq.repeatable, nq.repeatable and coalesce(lq.repeatable, 0) and nq.ordinal <= lq.ordinal, nq.ordinal
			limit 1;
		`, npcOID, playerOID).Scan(&questID, &containerID, &quest.Ordinal)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	quest.AssetContainer, err = s.GetAssetContainer(ctx, platform, containerID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	// The client takes the NPC of a quest from ParentOID.
	quest.ParentOID = npcOID

	// The client plays the children of a quest as its steps.
	quest.Children, err = s.getQuestSteps(ctx, platform, questID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	return &quest, nil
}

func (s *Service) getQuestSteps(ctx context.Context, platform gsf.Platform, questID int) ([]types.Quest, error) {
	type row struct {
		containerID int
		step        types.Quest
	}
	stepRows, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[row]{
		Query: "select container_id, ordinal from quest where parent_id = ? order by ordinal;",
		Args:  []any{questID},
		Scan:  func(rows *sql.Rows, r *row) error { return rows.Scan(&r.containerID, &r.step.Ordinal) },
	})
	if err != nil {
		return nil, err
	}

	var steps []types.Quest
	for _, r := range stepRows {
		step := r.step
		if step.AssetContainer, err = s.GetAssetContainer(ctx, platform, r.containerID); err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}
	return steps, nil
}

// CreatePlayerQuest creates the player quest and the player quests of its steps.
// It returns the player quest OID and the step player quest OIDs by step quest OID.
func (s *Service) CreatePlayerQuest(ctx context.Context, playerOID, questOID types.OID) (types.OID, map[string]types.OID, error) {
	const op = "game.Service.CreatePlayerQuest"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return types.OID{}, nil, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	var playerID, questID int
	err = tx.QueryRowContext(ctx, `
			select pl.id, q.id
			from player as pl
			join quest as q
			join asset_container as ac on ac.id = q.container_id
			where pl.gsfoid = ? and ac.gsfoid = ? and q.parent_id is null;
		`, playerOID, questOID).Scan(&playerID, &questID)
	if errors.Is(err, sql.ErrNoRows) {
		return types.OID{}, nil, wrap.IfErr(op, fmt.Errorf("quest %v of player %v: %w", questOID, playerOID, ErrQuestNotFound))
	}
	if err != nil {
		return types.OID{}, nil, wrap.IfErr(op, err)
	}

	// The NPC offers a quest again when the player left it open, e.g. after a disconnect,
	// the earlier open player quest of it and its steps are abandoned.
	_, err = tx.ExecContext(ctx, `
			update player_quest set quest_state_id = ?
			where player_id = ? and quest_state_id in (?, ?, ?, ?)
				and (quest_id = ? or parent_id in (select id from player_quest where player_id = ? and quest_id = ?));
		`, queststate.CLOSED_ABANDONED, playerID,
		queststate.AVAILABLE, queststate.ACCEPTED, queststate.COMPLETED, queststate.STARTED,
		questID, playerID, questID)
	if err != nil {
		return types.OID{}, nil, wrap.IfErr(op, err)
	}

	playerQuestID, playerQuestOID, err := insertPlayerQuest(ctx, tx, playerID, questID, nil)
	if err != nil {
		return types.OID{}, nil, wrap.IfErr(op, err)
	}

	type step struct {
		questID int
		oid     types.OID
	}
	steps, err := w2db.SelectContext(ctx, tx, w2db.SelectOptions[step]{
		Query: `
			select q.id, ac.gsfoid
			from quest as q
			join asset_container as ac on ac.id = q.container_id
			where q.parent_id = ?
			order by q.ordinal;
		`,
		Args: []any{questID},
		Scan: func(rows *sql.Rows, st *step) error { return rows.Scan(&st.questID, &st.oid) },
	})
	if err != nil {
		return types.OID{}, nil, wrap.IfErr(op, err)
	}

	// The client looks up a step by its quest OID as a decimal string.
	subQuests := make(map[string]types.OID, len(steps))
	for _, st := range steps {
		_, stepOID, err := insertPlayerQuest(ctx, tx, playerID, st.questID, &playerQuestID)
		if err != nil {
			return types.OID{}, nil, wrap.IfErr(op, err)
		}
		subQuests[strconv.FormatInt(st.oid.Int64(), 10)] = stepOID
	}

	return playerQuestOID, subQuests, wrap.IfErr(op, tx.Commit())
}

func insertPlayerQuest(ctx context.Context, tx *sql.Tx, playerID, questID int, parentID *int) (int, types.OID, error) {
	var id int
	var oid types.OID
	err := tx.QueryRowContext(ctx, `
			insert into player_quest (gsfoid, player_id, quest_id, quest_state_id, parent_id)
			values (
				(
					select coalesce(max(pq.gsfoid) + 1, t.min_oid)
					from oid_type as t
					left join player_quest as pq on pq.gsfoid between t.min_oid and t.max_oid
					where t.id = ?
				),
				?, ?, ?, ?
			)
			returning id, gsfoid;
		`, oidtype.PLAYER_QUEST, playerID, questID, queststate.AVAILABLE, parentID).Scan(&id, &oid)
	return id, oid, err
}

// SetPlayerQuestState sets the state of the player quest.
func (s *Service) SetPlayerQuestState(ctx context.Context, playerOID, playerQuestOID types.OID, state queststate.QuestState) error {
	const op = "game.Service.SetPlayerQuestState"
	err := s.store.DB().QueryRowContext(ctx, `
			update player_quest set quest_state_id = ?
			where gsfoid = ? and player_id = (select id from player where gsfoid = ?)
			returning id;
		`, state, playerQuestOID, playerOID).Scan(new(int))
	if errors.Is(err, sql.ErrNoRows) {
		return wrap.IfErr(op, fmt.Errorf("player quest %v: %w", playerQuestOID, ErrQuestNotFound))
	}
	return wrap.IfErr(op, err)
}

// AddPlayerQuestItem counts a collected item of a STARTED player quest,
// the items of other quests and of unknown player quests are ignored.
// It returns true when every rule amount is collected and the quest becomes COMPLETED.
func (s *Service) AddPlayerQuestItem(ctx context.Context, playerOID, playerQuestOID, itemOID types.OID, count int32) (bool, error) {
	const op = "game.Service.AddPlayerQuestItem"

	// The count comes from the client, a non-positive one would undo the progress.
	if count < 1 {
		return false, nil
	}

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return false, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	var playerQuestID, questID int
	var state queststate.QuestState
	err = tx.QueryRowContext(ctx, `
			select pq.id, pq.quest_id, pq.quest_state_id
			from player_quest as pq
			join player as pl on pl.id = pq.player_id
			where pq.gsfoid = ? and pl.gsfoid = ?;
		`, playerQuestOID, playerOID).Scan(&playerQuestID, &questID, &state)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, wrap.IfErr(op, err)
	}
	if state != queststate.STARTED {
		return false, nil
	}

	_, err = tx.ExecContext(ctx, `
			insert into player_quest_rule (player_quest_id, quest_rule_id, count)
			select ?, qr.id, ?
			from quest_rule as qr
			join asset_container as ac on ac.id = qr.item_container_id
			where qr.quest_id = ? and ac.gsfoid = ?
			on conflict (player_quest_id, quest_rule_id) do update set count = count + excluded.count;
		`, playerQuestID, count, questID, itemOID)
	if err != nil {
		return false, wrap.IfErr(op, err)
	}

	var completed bool
	err = tx.QueryRowContext(ctx, `
			select coalesce(min(coalesce(pqr.count, 0) >= qr.amount), 0)
			from quest_rule as qr
			left join player_quest_rule as pqr on pqr.quest_rule_id = qr.id and pqr.player_quest_id = ?
			where qr.quest_id = ?;
		`, playerQuestID, questID).Scan(&completed)
	if err != nil {
		return false, wrap.IfErr(op, err)
	}

	if completed {
		_, err = tx.ExecContext(ctx, `update player_quest set quest_state_id = ? where id = ?;`, queststate.COMPLETED, playerQuestID)
		if err != nil {
			return false, wrap.IfErr(op, err)
		}
	}

	return completed, wrap.IfErr(op, tx.Commit())
}

// ClosePlayerQuest closes the COMPLETED player quest when the client takes the reward.
// After the last step it also completes the parent quest and returns its OID.
// Returns false when the quest is not completed.
func (s *Service) ClosePlayerQuest(ctx context.Context, playerOID, playerQuestOID types.OID) (bool, *types.OID, error) {
	const op = "game.Service.ClosePlayerQuest"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return false, nil, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	var parentID sql.NullInt64
	err = tx.QueryRowContext(ctx, `
			update player_quest set quest_state_id = ?, completed_at = unixepoch()
			where gsfoid = ? and quest_state_id = ? and player_id = (select id from player where gsfoid = ?)
			returning parent_id;
		`, queststate.CLOSED, playerQuestOID, queststate.COMPLETED, playerOID).Scan(&parentID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, wrap.IfErr(op, err)
	}

	var parentOID *types.OID
	if parentID.Valid {
		var oid types.OID
		err = tx.QueryRowContext(ctx, `
				update player_quest set quest_state_id = ?
				where id = ? and quest_state_id = ?
					and not exists (select 1 from player_quest where parent_id = ? and quest_state_id != ?)
				returning gsfoid;
			`, queststate.COMPLETED, parentID.Int64, queststate.STARTED, parentID.Int64, queststate.CLOSED).Scan(&oid)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return false, nil, wrap.IfErr(op, err)
		}
		if err == nil {
			parentOID = &oid
		}
	}

	return true, parentOID, wrap.IfErr(op, tx.Commit())
}

// AbandonPlayerQuest abandons or fails the player quest and its open steps.
func (s *Service) AbandonPlayerQuest(ctx context.Context, playerOID, playerQuestOID types.OID, state queststate.QuestState) error {
	const op = "game.Service.AbandonPlayerQuest"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	var playerQuestID int
	err = tx.QueryRowContext(ctx, `
			update player_quest set quest_state_id = ?
			where gsfoid = ? and player_id = (select id from player where gsfoid = ?)
			returning id;
		`, state, playerQuestOID, playerOID).Scan(&playerQuestID)
	if errors.Is(err, sql.ErrNoRows) {
		return wrap.IfErr(op, fmt.Errorf("player quest %v: %w", playerQuestOID, ErrQuestNotFound))
	}
	if err != nil {
		return wrap.IfErr(op, err)
	}

	_, err = tx.ExecContext(ctx, `
			update player_quest set quest_state_id = ?
			where parent_id = ? and quest_state_id != ?;
		`, state, playerQuestID, queststate.CLOSED)
	if err != nil {
		return wrap.IfErr(op, err)
	}

	return wrap.IfErr(op, tx.Commit())
}
