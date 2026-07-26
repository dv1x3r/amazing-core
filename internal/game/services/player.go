package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetPlayer(ctx context.Context, platform gsf.Platform, playerOID types.OID) (types.Player, error) {
	const op = "game.Service.GetPlayer"

	row := s.store.DB().QueryRowContext(ctx, `
			select
				pl.gsfoid,
				pl.created_at,
				pl.is_tutorial_completed,
				pl.is_qa
			from player as pl
			where pl.gsfoid = ?;
		`, playerOID)

	var player types.Player

	if err := row.Scan(
		&player.OID,
		&player.CreateDate,
		&player.IsTutorialCompleted,
		&player.IsQA,
	); err != nil {
		return player, wrap.IfErr(op, err)
	}

	activeAvatar, err := s.GetActivePlayerAvatar(ctx, platform, playerOID)
	if err != nil {
		return player, wrap.IfErr(op, err)
	}
	player.ActivePlayerAvatar = activeAvatar

	return player, nil
}

func (s *Service) GetPlayerInfoTO(ctx context.Context, playerOID types.OID) (types.PlayerInfoTO, error) {
	const op = "game.Service.GetPlayerInfoTO"

	row := s.store.DB().QueryRowContext(ctx, `
			select
				pl.username,
				coalesce(pl.is_chat_allowed, 1),
				pl.is_qa
			from player as pl
			where pl.gsfoid = ?;
		`, playerOID)

	var playerInfo types.PlayerInfoTO
	if err := row.Scan(
		&playerInfo.PlayerName,
		&playerInfo.ChatAllowed,
		&playerInfo.QA,
	); err != nil {
		return playerInfo, wrap.IfErr(op, err)
	}

	return playerInfo, nil
}
