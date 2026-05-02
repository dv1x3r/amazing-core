package player

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetGSFPlayer(ctx context.Context, platform gsf.Platform, playerID int) (types.Player, error) {
	const op = "player.Service.GetGSFPlayer"

	row := s.store.DB().QueryRowContext(ctx, `
			select
				pl.gsfoid,
				pl.created_at,
				pl.is_tutorial_completed,
				pl.is_qa
			from player as pl
			where pl.id = ?;
		`, playerID)

	var player types.Player

	if err := row.Scan(
		&player.OID,
		&player.CreateDate,
		&player.IsTutorialCompleted,
		&player.IsQA,
	); err != nil {
		return player, wrap.IfErr(op, err)
	}

	activeAvatar, err := s.GetGSFActiveAvatar(ctx, platform, playerID)
	if err != nil {
		return player, wrap.IfErr(op, err)
	}
	player.ActivePlayerAvatar = activeAvatar

	return player, nil
}

func (s *Service) GetGSFAvatars(ctx context.Context, platform gsf.Platform, playerID int) ([]types.PlayerAvatar, error) {
	const op = "player.Service.GetGSFAvatars"
	var avatars []types.PlayerAvatar

	rows, err := s.store.DB().QueryContext(ctx, `
			select
				pa.gsfoid,
				pl.gsfoid as player_oid,
				pa.name as player_name,
				av.name as avatar_name,
				av.container_id
			from player_avatar as pa
			join player as pl on pl.id = pa.player_id
			join avatar as av on av.id = pa.avatar_id
			where pa.player_id = ?;
		`, playerID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var avatar types.PlayerAvatar
		var containerID int
		if err := rows.Scan(
			&avatar.OID,
			&avatar.PlayerID,
			&avatar.Name,
			&avatar.Avatar.Name,
			&containerID,
		); err != nil {
			return avatars, wrap.IfErr(op, err)
		}
		container, err := s.assets.GetGSFAssetContainer(ctx, platform, containerID)
		if err != nil {
			return avatars, wrap.IfErr(op, err)
		}
		avatar.Avatar.AssetContainer = container
		avatars = append(avatars, avatar)
	}

	return avatars, wrap.IfErr(op, rows.Err())
}

func (s *Service) GetGSFActiveAvatar(ctx context.Context, platform gsf.Platform, playerID int) (types.PlayerAvatar, error) {
	const op = "player.Service.GetGSFActiveAvatar"

	row := s.store.DB().QueryRowContext(ctx, `
			select
				pa.gsfoid,
				pl.gsfoid as player_oid,
				pa.name as player_name,
				av.name as avatar_name,
				av.container_id
			from player_avatar as pa
			join player as pl on pl.id = pa.player_id
			join avatar as av on av.id = pa.avatar_id
			where pa.player_id = ? and pa.is_active = 1;
		`, playerID)

	var avatar types.PlayerAvatar
	var containerID int

	if err := row.Scan(
		&avatar.OID,
		&avatar.PlayerID,
		&avatar.Name,
		&avatar.Avatar.Name,
		&containerID,
	); err != nil {
		return avatar, wrap.IfErr(op, err)
	}

	container, err := s.assets.GetGSFAssetContainer(ctx, platform, containerID)
	if err != nil {
		return avatar, wrap.IfErr(op, err)
	}
	avatar.Avatar.AssetContainer = container

	return avatar, nil
}

func (s *Service) SetGSFPlayerActiveAvatar(ctx context.Context, platform gsf.Platform, oid types.OID) (types.PlayerAvatar, error) {
	const op = "player.Service.SetGSFPlayerActiveAvatar"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return types.PlayerAvatar{}, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	var playerAvatarID, playerID int
	row := tx.QueryRowContext(ctx, "select id, player_id from player_avatar where gsfoid = ?;", oid)
	if err := row.Scan(&playerAvatarID, &playerID); err != nil {
		return types.PlayerAvatar{}, wrap.IfErr(op, err)
	}

	if err := s.setActivePlayerAvatarByID(ctx, tx, playerAvatarID, playerID); err != nil {
		return types.PlayerAvatar{}, wrap.IfErr(op, err)
	}

	if err := tx.Commit(); err != nil {
		return types.PlayerAvatar{}, wrap.IfErr(op, err)
	}

	avatar, err := s.GetGSFActiveAvatar(ctx, platform, playerID)
	return avatar, wrap.IfErr(op, err)
}
