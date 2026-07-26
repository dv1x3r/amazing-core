package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetAvatar(ctx context.Context, platform gsf.Platform, avatarID int) (types.Avatar, error) {
	const op = "game.Service.GetAvatar"

	row := s.store.DB().QueryRowContext(ctx, `
			select container_id, name
			from avatar
			where id = ?;
		`, avatarID)

	var containerID int
	var avatar types.Avatar

	if err := row.Scan(
		&containerID,
		&avatar.Name,
	); err != nil {
		return avatar, wrap.IfErr(op, err)
	}

	container, err := s.GetAssetContainer(ctx, platform, containerID)
	if err != nil {
		return avatar, wrap.IfErr(op, err)
	}
	avatar.AssetContainer = container

	return avatar, nil
}
