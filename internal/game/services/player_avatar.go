package services

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2db"
)

func (s *Service) GetPlayerAvatars(ctx context.Context, platform gsf.Platform, playerOID types.OID) ([]types.PlayerAvatar, error) {
	const op = "game.Service.GetPlayerAvatars"

	type row struct {
		playerAvatar types.PlayerAvatar
		avatarID     int
		maxOutfits   int16
	}
	avatarRows, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[row]{
		Query: `
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
		`,
		Args: []any{playerOID},
		Scan: func(rows *sql.Rows, r *row) error {
			return rows.Scan(
				&r.playerAvatar.OID,
				&r.playerAvatar.PlayerOID,
				&r.playerAvatar.Name,
				&r.playerAvatar.OutfitNo,
				&r.playerAvatar.PlayerAvatarOutfitOID,
				&r.avatarID,
				&r.maxOutfits,
			)
		},
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	playerAvatars := make([]types.PlayerAvatar, 0, len(avatarRows))
	for _, r := range avatarRows {
		avatar, err := s.GetAvatar(ctx, platform, r.avatarID)
		if err != nil {
			return playerAvatars, wrap.IfErr(op, err)
		}
		avatar.MaxOutfits = r.maxOutfits
		playerAvatar := r.playerAvatar
		playerAvatar.Avatar = avatar
		playerAvatars = append(playerAvatars, playerAvatar)
	}

	return playerAvatars, nil
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

func (s *Service) SetPlayerActiveAvatar(ctx context.Context, platform gsf.Platform, playerOID, playerAvatarOID types.OID) (types.PlayerAvatar, error) {
	const op = "game.Service.SetPlayerActiveAvatar"

	// Only the selected avatar of the player is active, an avatar of another player changes nothing.
	err := s.store.DB().QueryRowContext(ctx, `
			update player_avatar set is_active = (gsfoid = ?)
			where player_id = (
				select pa.player_id
				from player_avatar as pa
				join player as pl on pl.id = pa.player_id
				where pa.gsfoid = ? and pl.gsfoid = ?
			)
			returning id;
		`, playerAvatarOID, playerAvatarOID, playerOID).Scan(new(int))
	if err != nil {
		return types.PlayerAvatar{}, wrap.IfErr(op, err)
	}

	avatar, err := s.GetActivePlayerAvatar(ctx, platform, playerOID)
	return avatar, wrap.IfErr(op, err)
}
