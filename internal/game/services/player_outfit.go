package services

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2db"
)

func (s *Service) GetOutfits(ctx context.Context, playerAvatarOID, playerOID types.OID) ([]types.PlayerAvatarOutfit, error) {
	const op = "game.Service.GetOutfits"

	outfits, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[types.PlayerAvatarOutfit]{
		Query: `
			select
				po.gsfoid,
				pl.gsfoid as player_gsfoid,
				pa.gsfoid as player_avatar_gsfoid,
				po.outfit_no
			from player_avatar_outfit as po
			join player_avatar as pa on pa.id = po.player_avatar_id
			join player as pl on pl.id = pa.player_id
			where pa.gsfoid = ? and pl.gsfoid = ?
			order by po.outfit_no;
		`,
		Args: []any{playerAvatarOID, playerOID},
		Scan: func(rows *sql.Rows, outfit *types.PlayerAvatarOutfit) error {
			return rows.Scan(
				&outfit.OID,
				&outfit.PlayerOID,
				&outfit.PlayerAvatarOID,
				&outfit.OutfitNo,
			)
		},
	})
	if outfits == nil {
		outfits = []types.PlayerAvatarOutfit{}
	}
	return outfits, wrap.IfErr(op, err)
}

// CreateOutfit creates an outfit of the player avatar, an avatar of another player fails the not null constraint.
func (s *Service) CreateOutfit(ctx context.Context, playerOID, playerAvatarOID types.OID, outfitNo int16) (types.OID, error) {
	const op = "game.Service.CreateOutfit"
	var outfitOID types.OID
	row := s.store.DB().QueryRowContext(ctx, `
			insert into player_avatar_outfit (gsfoid, player_avatar_id, outfit_no)
			values (
				(select max(gsfoid) + 1 from player_avatar_outfit),
				(
					select pa.id
					from player_avatar as pa
					join player as pl on pl.id = pa.player_id
					where pa.gsfoid = ? and pl.gsfoid = ?
				),
				?
			) returning gsfoid;
		`, playerAvatarOID, playerOID, outfitNo)
	return outfitOID, wrap.IfErr(op, row.Scan(&outfitOID))
}

func (s *Service) SetPlayerActiveOutfit(ctx context.Context, playerOID, playerAvatarOID types.OID, outfitNo int16) error {
	const op = "game.Service.SetPlayerActiveOutfit"
	_, err := s.store.DB().ExecContext(ctx, `
			update player_avatar set outfit_no = ?
			where gsfoid = ? and player_id = (select id from player where gsfoid = ?);
		`, outfitNo, playerAvatarOID, playerOID)
	return wrap.IfErr(op, err)
}
