package player

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
	"github.com/dv1x3r/w2go/w2sql"

	"github.com/huandu/go-sqlbuilder"
)

type PlayerList struct {
	ID         int    `json:"id"`
	OID        int64  `json:"oid"`
	OIDStr     string `json:"oid_str"`
	ActiveName string `json:"active_name"`
}

type PlayerDetails struct {
	ID                  int         `json:"id"`
	OID                 int64       `json:"oid"`
	CreatedAt           w2.UnixTime `json:"created_at"`
	IsTutorialCompleted bool        `json:"is_tutorial_completed"`
	IsQA                bool        `json:"is_qa"`
	ActiveAvatar        w2.Dropdown `json:"active_avatar"`
}

type PlayerAvatar struct {
	ID       int              `json:"id"`
	OID      w2.Field[int64]  `json:"oid"`
	OIDStr   string           `json:"oid_str"`
	Name     w2.Field[string] `json:"name"`
	Avatar   w2.Dropdown      `json:"avatar"`
	OutfitNo w2.Field[int]    `json:"outfit_no"`
	IsActive bool             `json:"is_active"`
}

type PlayerOutfit struct {
	ID           int             `json:"id"`
	OID          w2.Field[int64] `json:"oid"`
	OIDStr       string          `json:"oid_str"`
	PlayerAvatar w2.Dropdown     `json:"player_avatar"`
	OutfitNo     w2.Field[int]   `json:"outfit_no"`
}

func (s *Service) GetPlayerListGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[PlayerList], error) {
	const op = "player.Service.GetPlayerListGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerList]{
		From: "player as p",
		Select: []string{
			"p.id",
			"p.gsfoid",
			"coalesce(pa.name, '[None]') as active_name",
		},
		WhereMapping: map[string]string{
			"id":          "p.id",
			"oid":         "p.gsfoid",
			"active_name": "pa.name",
		},
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin, `(
				select id, player_id, name,
					row_number() over (partition by player_id order by id) as rn
				from player_avatar
				where is_active = 1
				) as pa`,
				"pa.player_id = p.id and pa.rn = 1",
			)
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.OrderByAsc("p.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerList) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.ActiveName,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromInt64(record.OID).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetPlayerDetailsForm(ctx context.Context, req w2.GetFormRequest) (w2.GetFormResponse[PlayerDetails], error) {
	const op = "player.Service.GetPlayerDetailsForm"
	res, err := w2db.GetFormContext(ctx, s.store.DB(), req, w2db.GetFormOptions[PlayerDetails]{
		From:    "player as p",
		IDField: "p.id",
		Select: []string{
			"p.id",
			"p.gsfoid",
			"p.created_at",
			"p.is_tutorial_completed",
			"p.is_qa",
			"pa.id as active_avatar_id",
			"pa.name as active_avatar_name",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin, `(
				select id, player_id, name,
					row_number() over (partition by player_id order by id) as rn
				from player_avatar
				where is_active = 1
				) as pa`,
				"pa.player_id = p.id and pa.rn = 1",
			)
		},
		Scan: func(row *sql.Row, record *PlayerDetails) error {
			return row.Scan(
				&record.ID,
				&record.OID,
				&record.CreatedAt,
				&record.IsTutorialCompleted,
				&record.IsQA,
				&record.ActiveAvatar.ID,
				&record.ActiveAvatar.Text,
			)
		},
	})
	if errors.Is(err, sql.ErrNoRows) {
		return res, wrap.IfErr(op, ErrPlayerNotFound)
	}
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetPlayerAvatarGrid(ctx context.Context, req w2.GetGridRequest, id int) (w2.GetGridResponse[PlayerAvatar], error) {
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
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.Where(sb.EQ("pa.player_id", id))
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("avatar as av", "av.id = pa.avatar_id")
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
			record.OIDStr = types.OIDFromInt64(record.OID.V).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetPlayerOutfitGrid(ctx context.Context, req w2.GetGridRequest, id int) (w2.GetGridResponse[PlayerOutfit], error) {
	const op = "player.Service.GetPlayerOutfitGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerOutfit]{
		From: "player_avatar_outfit as po",
		Select: []string{
			"po.id",
			"po.gsfoid",
			"po.player_avatar_id",
			"pa.name as player_avatar_name",
			"po.outfit_no",
		},
		WhereMapping: map[string]string{
			"id":            "po.id",
			"oid":           "po.gsfoid",
			"player_avatar": "pa.name",
			"outfit_no":     "po.outfit_no",
		},
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("player_avatar as pa", "pa.id = po.player_avatar_id")
			sb.Where(sb.EQ("pa.player_id", id))
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.OrderByAsc("po.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerOutfit) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.PlayerAvatar.ID,
				&record.PlayerAvatar.Text,
				&record.OutfitNo,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromInt64(record.OID.V).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetPlayerAvatarsDropdown(ctx context.Context, req w2.GetDropdownRequest, id int) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "player.Service.GetPlayerAvatarsDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "player_avatar",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "name",
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Where(sb.EQ("player_id", id))
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerDetails(ctx context.Context, req w2.SaveFormRequest[PlayerDetails]) error {
	const op = "player.Service.UpdatePlayerDetails"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		affected, err := w2db.UpdateFormContext(ctx, tx, req, w2db.UpdateFormOptions{
			Update:  "player",
			IDField: "id",
			Cols:    []string{"gsfoid", "is_tutorial_completed", "is_qa"},
			Values:  []any{req.Record.OID, req.Record.IsTutorialCompleted, req.Record.IsQA},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrPlayerExists
		} else if affected == 0 && err == nil {
			return ErrPlayerNotFound
		} else if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx,
			"update player_avatar set is_active = case when id = ? then 1 else 0 end where player_id = ?;",
			req.Record.ActiveAvatar.ID,
			req.RecID,
		)
		return err
	})
	return wrap.IfErr(op, err)
}

func (s *Service) CreatePlayerAvatar(ctx context.Context, req w2.SaveFormRequest[PlayerAvatar], playerID int) error {
	const op = "player.Service.CreatePlayerAvatar"
	_, err := w2db.InsertFormContext(ctx, s.store.DB(), req, w2db.InsertFormOptions{
		Into: "player_avatar",
		Cols: []string{"player_id", "name", "avatar_id", "gsfoid"},
		Values: []any{
			playerID,
			req.Record.Name,
			req.Record.Avatar.ID,
			sqlbuilder.Buildf("(select coalesce(max(gsfoid) + 1, 1) from player_avatar where player_id = %v)", playerID),
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrPlayerAvatarExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerAvatars(ctx context.Context, req w2.SaveGridRequest[PlayerAvatar]) error {
	const op = "player.Service.UpdatePlayerAvatars"
	_, err := w2db.SaveGridContext(ctx, s.store.DB(), req, w2db.SaveGridOptions[PlayerAvatar]{
		BuildUpdate: func(change PlayerAvatar) *sqlbuilder.UpdateBuilder {
			ub := sqlbuilder.Update("player_avatar")
			w2sql.Set(ub, change.OID, "gsfoid")
			w2sql.Set(ub, change.Name, "name")
			w2sql.Set(ub, change.OutfitNo, "outfit_no")
			w2sql.Set(ub, change.Avatar.ID, "avatar_id")
			ub.Where(ub.EQ("id", change.ID))
			return ub
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrPlayerAvatarExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) CreatePlayerOutfit(ctx context.Context, req w2.SaveFormRequest[PlayerOutfit]) error {
	const op = "player.Service.CreatePlayerOutfit"
	_, err := w2db.InsertFormContext(ctx, s.store.DB(), req, w2db.InsertFormOptions{
		Into: "player_avatar_outfit",
		Cols: []string{"player_avatar_id", "outfit_no", "gsfoid"},
		Values: []any{
			req.Record.PlayerAvatar.ID,
			req.Record.OutfitNo,
			sqlbuilder.Buildf("(select coalesce(max(gsfoid) + 1, 1) from player_avatar_outfit)"),
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrPlayerOutfitExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) UpdatePlayerOutfits(ctx context.Context, req w2.SaveGridRequest[PlayerOutfit]) error {
	const op = "player.Service.UpdatePlayerOutfits"
	_, err := w2db.SaveGridContext(ctx, s.store.DB(), req, w2db.SaveGridOptions[PlayerOutfit]{
		BuildUpdate: func(change PlayerOutfit) *sqlbuilder.UpdateBuilder {
			ub := sqlbuilder.Update("player_avatar_outfit")
			w2sql.Set(ub, change.OID, "gsfoid")
			w2sql.Set(ub, change.PlayerAvatar.ID, "player_avatar_id")
			w2sql.Set(ub, change.OutfitNo, "outfit_no")
			ub.Where(ub.EQ("id", change.ID))
			return ub
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrPlayerOutfitExists
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

func (s *Service) DeletePlayerOutfits(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "player.Service.DeletePlayerOutfits"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "player_avatar_outfit",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
