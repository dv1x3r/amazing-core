package player

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type PlayerAvatar struct {
	ID       int              `json:"id"`
	OID      w2.Field[string] `json:"oid"`
	OIDStr   string           `json:"oid_str"`
	Name     w2.Field[string] `json:"name"`
	Avatar   w2.Dropdown      `json:"avatar"`
	OutfitNo w2.Field[int]    `json:"outfit_no"`
	IsActive bool             `json:"is_active"`
}

func (s *Service) GetPlayerAvatarGrid(ctx context.Context, req w2.GetGridRequest, playerID int) (w2.GetGridResponse[PlayerAvatar], error) {
	const op = "player.Service.GetPlayerAvatarGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerAvatar]{
		From: "player_avatar as pa",
		Select: []string{
			"pa.id",
			"pa.gsfoid",
			"pa.name",
			"pa.outfit_no",
			"pa.is_active",
			"pa.avatar_id",
			"av.name as avatar_name",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("avatar as av", "av.id = pa.avatar_id")
			sb.Where(sb.EQ("pa.player_id", playerID))
			sb.OrderByAsc("pa.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerAvatar) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.Name,
				&record.OutfitNo,
				&record.IsActive,
				&record.Avatar.ID,
				&record.Avatar.Text,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromString(record.OID.V).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreatePlayerAvatar(ctx context.Context, req w2.SaveFormRequest[PlayerAvatar], playerID int) (int, error) {
	const op = "player.Service.CreatePlayerAvatar"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "player_avatar",
		Cols: []string{"player_id", "gsfoid", "avatar_id", "name", "outfit_no"},
		Values: []any{
			playerID,
			sqlbuilder.Buildf("(select coalesce(max(gsfoid) + 1, 1) from player_avatar)"),
			req.Record.Avatar.ID,
			req.Record.Name,
			req.Record.OutfitNo,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrPlayerAvatarExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerAvatar(ctx context.Context, req w2.SaveFormRequest[PlayerAvatar]) error {
	const op = "player.Service.UpdatePlayerAvatar"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update:  "player_avatar",
		Cols:    []string{"gsfoid", "avatar_id", "name", "outfit_no"},
		Values:  []any{req.Record.OID, req.Record.Avatar.ID, req.Record.Name, req.Record.OutfitNo},
		IDField: "id",
		IDValue: req.Record.ID,
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrPlayerAvatarExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeletePlayerAvatars(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "player.Service.DeletePlayerAvatars"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "player_avatar",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) GetGSFPlayerAvatars(ctx context.Context, platform gsf.Platform, playerOID types.OID) ([]types.PlayerAvatar, error) {
	const op = "player.Service.GetGSFPlayerAvatars"

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
		avatar, err := s.avatars.GetGSFAvatar(ctx, platform, avatarID)
		if err != nil {
			return playerAvatars, wrap.IfErr(op, err)
		}
		avatar.MaxOutfits = maxOutfits
		playerAvatar.Avatar = avatar
		playerAvatars = append(playerAvatars, playerAvatar)
	}

	return playerAvatars, wrap.IfErr(op, rows.Err())
}

func (s *Service) GetGSFActivePlayerAvatar(ctx context.Context, platform gsf.Platform, playerOID types.OID) (types.PlayerAvatar, error) {
	const op = "player.Service.GetGSFActivePlayerAvatar"

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

	avatar, err := s.avatars.GetGSFAvatar(ctx, platform, avatarID)
	if err != nil {
		return playerAvatar, wrap.IfErr(op, err)
	}
	playerAvatar.Avatar = avatar

	return playerAvatar, nil
}

func (s *Service) SetGSFPlayerActiveAvatar(ctx context.Context, platform gsf.Platform, playerAvatarOID types.OID) (types.PlayerAvatar, error) {
	const op = "player.Service.SetGSFPlayerActiveAvatar"

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

	avatar, err := s.GetGSFActivePlayerAvatar(ctx, platform, playerOID)
	return avatar, wrap.IfErr(op, err)
}
