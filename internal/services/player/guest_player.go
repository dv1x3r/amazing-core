package player

import (
	"context"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) CreateGSFGuestPlayer(ctx context.Context, platform gsf.Platform) (types.Player, error) {
	const op = "player.Service.CreateGSFGuestPlayer"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return types.Player{}, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	var playerID int
	var playerOID types.OID
	err = tx.QueryRowContext(ctx, `
			insert into player (gsfoid, is_tutorial_completed, is_qa, max_outfits)
			values ((select coalesce(max(gsfoid), 0) + 1 from player), 1, 1, 1)
			returning id, gsfoid;
		`).Scan(&playerID, &playerOID)
	if s.store.IsErrConstraintUnique(err) {
		return types.Player{}, wrap.IfErr(op, ErrPlayerExists)
	}
	if err != nil {
		return types.Player{}, wrap.IfErr(op, err)
	}

	var playerAvatarID int
	err = tx.QueryRowContext(ctx, `
			insert into player_avatar (gsfoid, player_id, avatar_id, name, outfit_no, is_active)
			values (
				(select coalesce(max(gsfoid), 0) + 1 from player_avatar),
				?,
				(select id from avatar order by random() limit 1),
				?,
				1,
				1
			)
			returning id;
		`, playerID, fmt.Sprintf("Guest %d", playerOID.Int64())).Scan(&playerAvatarID)
	if s.store.IsErrConstraintUnique(err) {
		return types.Player{}, wrap.IfErr(op, ErrPlayerAvatarExists)
	}
	if err != nil {
		return types.Player{}, wrap.IfErr(op, err)
	}

	if _, err = tx.ExecContext(ctx, `
			insert into player_avatar_outfit (gsfoid, player_avatar_id, outfit_no)
			values (
				(select coalesce(max(gsfoid), 0) + 1 from player_avatar_outfit),
				?,
				1
			);
		`, playerAvatarID); err != nil {
		return types.Player{}, wrap.IfErr(op, err)
	}

	if _, err = tx.ExecContext(ctx, `
			insert into player_item (gsfoid, player_id, item_id, quantity)
			select
				(select coalesce(max(gsfoid), 0) from player_item) + row_number() over (order by it.id),
				?,
				it.id,
				1
			from item as it;
		`, playerID); err != nil {
		return types.Player{}, wrap.IfErr(op, err)
	}

	if err = tx.Commit(); err != nil {
		return types.Player{}, wrap.IfErr(op, err)
	}

	player, err := s.GetGSFPlayer(ctx, platform, playerOID)
	return player, wrap.IfErr(op, err)
}
