package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetOtherPlayerDetails(ctx context.Context, platform gsf.Platform, playerOID types.OID) (types.OtherPlayerDetails, error) {
	const op = "game.Service.GetOtherPlayerDetails"

	player, err := s.GetPlayer(ctx, platform, playerOID)
	if err != nil {
		return types.OtherPlayerDetails{}, wrap.IfErr(op, err)
	}

	var avatarCount int
	if err := s.store.DB().QueryRowContext(ctx, `
			select count(*)
			from player_avatar as pa
			join player as pl on pl.id = pa.player_id
			where pl.gsfoid = ?;
		`, playerOID).Scan(&avatarCount); err != nil {
		return types.OtherPlayerDetails{}, wrap.IfErr(op, err)
	}

	clothing := []types.PlayerItem{}
	if player.ActivePlayerAvatar.PlayerAvatarOutfitOID.Int64() != 0 {
		clothing, err = s.GetOutfitItems(ctx, platform, player.ActivePlayerAvatar.PlayerAvatarOutfitOID, player.OID)
		if err != nil {
			return types.OtherPlayerDetails{}, wrap.IfErr(op, err)
		}
	}

	return types.OtherPlayerDetails{
		PlayerAvatar:      player.ActivePlayerAvatar,
		Clothing:          clothing,
		PlayerName:        player.ActivePlayerAvatar.Name,
		WorldName:         player.ActivePlayerAvatar.Name,
		TierOID:           types.OID{},
		PlayerAvatarCount: int32(avatarCount),
		Level:             1,
		XP:                0,
		Token:             0,
		Energy:            0,
		PlayerFriendCount: 0,
		Findable:          true,
		FindableDuration:  0,
		ExternalSites:     []types.PlayerExternalSite{},
	}, nil
}
