package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetPlayerAvatars(ctx context.Context, platform gsf.Platform, playerOID types.OID) ([]types.PlayerAvatar, error) {
	const op = "game.Service.GetPlayerAvatars"

	rows, err := s.store.DB().QueryContext(ctx, `
			select
				pa.gsfoid,
				pl.gsfoid as player_gsfoid,
				pa.name as player_avatar_name,
				pa.outfit_no,
				pao.gsfoid as player_avatar_outfit_gsfoid,
				pa.avatar_id,
				pl.max_outfits
			from player_avatar as pa
			join player as pl on pl.id = pa.player_id
			left join player_avatar_outfit as pao on pao.player_avatar_id = pa.id and pao.outfit_no = pa.outfit_no
			where pl.gsfoid = ?;
		`, playerOID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	playerAvatars := []types.PlayerAvatar{}
	for rows.Next() {
		var playerAvatar types.PlayerAvatar
		var avatarID int
		var maxOutfits int16
		if err := rows.Scan(
			&playerAvatar.OID,
			&playerAvatar.PlayerOID,
			&playerAvatar.Name,
			&playerAvatar.OutfitNo,
			&playerAvatar.PlayerAvatarOutfitOID,
			&avatarID,
			&maxOutfits,
		); err != nil {
			return playerAvatars, wrap.IfErr(op, err)
		}
		avatar, err := s.GetAvatar(ctx, platform, avatarID)
		if err != nil {
			return playerAvatars, wrap.IfErr(op, err)
		}
		avatar.MaxOutfits = maxOutfits
		playerAvatar.Avatar = avatar
		playerAvatars = append(playerAvatars, playerAvatar)
	}

	return playerAvatars, wrap.IfErr(op, rows.Err())
}

func (s *Service) GetActivePlayerAvatar(ctx context.Context, platform gsf.Platform, playerOID types.OID) (types.PlayerAvatar, error) {
	const op = "game.Service.GetActivePlayerAvatar"

	row := s.store.DB().QueryRowContext(ctx, `
			select
				pa.gsfoid,
				pl.gsfoid as player_gsfoid,
				pa.name as player_avatar_name,
				pa.outfit_no,
				pao.gsfoid as player_avatar_outfit_gsfoid,
				pa.avatar_id
			from player_avatar as pa
			join player as pl on pl.id = pa.player_id
			left join player_avatar_outfit as pao on pao.player_avatar_id = pa.id and pao.outfit_no = pa.outfit_no
			where pl.gsfoid = ? and pa.is_active = 1;
		`, playerOID)

	var playerAvatar types.PlayerAvatar
	var avatarID int

	if err := row.Scan(
		&playerAvatar.OID,
		&playerAvatar.PlayerOID,
		&playerAvatar.Name,
		&playerAvatar.OutfitNo,
		&playerAvatar.PlayerAvatarOutfitOID,
		&avatarID,
	); err != nil {
		return playerAvatar, wrap.IfErr(op, err)
	}

	avatar, err := s.GetAvatar(ctx, platform, avatarID)
	if err != nil {
		return playerAvatar, wrap.IfErr(op, err)
	}
	playerAvatar.Avatar = avatar

	return playerAvatar, nil
}

func (s *Service) SetPlayerActiveAvatar(ctx context.Context, platform gsf.Platform, playerAvatarOID types.OID) (types.PlayerAvatar, error) {
	const op = "game.Service.SetPlayerActiveAvatar"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return types.PlayerAvatar{}, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	var playerAvatarID, playerID int
	var playerOID types.OID
	row := tx.QueryRowContext(ctx, `
			select pa.id, pa.player_id, pl.gsfoid
			from player_avatar as pa
			join player as pl on pl.id = pa.player_id
			where pa.gsfoid = ?;
		`, playerAvatarOID)
	if err := row.Scan(&playerAvatarID, &playerID, &playerOID); err != nil {
		return types.PlayerAvatar{}, wrap.IfErr(op, err)
	}

	if _, err := tx.ExecContext(ctx, `
		update player_avatar
		set is_active = case when id = ? then 1 else 0 end
		where player_id = ?;
	`, playerAvatarID, playerID); err != nil {
		return types.PlayerAvatar{}, wrap.IfErr(op, err)
	}

	if err := tx.Commit(); err != nil {
		return types.PlayerAvatar{}, wrap.IfErr(op, err)
	}

	avatar, err := s.GetActivePlayerAvatar(ctx, platform, playerOID)
	return avatar, wrap.IfErr(op, err)
}
