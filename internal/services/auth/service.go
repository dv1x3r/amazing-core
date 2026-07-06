package auth

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/pass"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

var (
	ErrBlankCredentials          = errors.New("blank credentials")
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrCreatePlayerAccountFailed = errors.New("create player account failed")
)

type Service struct {
	logger *slog.Logger
	store  db.Store
}

func NewService(logger *slog.Logger, store db.Store) *Service {
	return &Service{
		logger: logger,
		store:  store,
	}
}

func (s *Service) LoginOrCreatePlayer(ctx context.Context, username, password string) (types.OID, error) {
	const op = "auth.Service.LoginOrCreatePlayer"

	username = strings.TrimSpace(username)
	if username == "" || strings.TrimSpace(password) == "" {
		return types.OID{}, wrap.IfErr(op, ErrBlankCredentials)
	}

	var playerOID types.OID
	var passwordHash string
	err := s.store.DB().QueryRowContext(ctx, `
			select gsfoid, password
			from player
			where username = ?;
		`, username).Scan(&playerOID, &passwordHash)
	if errors.Is(err, sql.ErrNoRows) {
		oid, err := s.createPlayer(ctx, username, password)
		return oid, wrap.IfErr(op, err)
	}
	if err != nil {
		return types.OID{}, wrap.IfErr(op, err)
	}

	ok, err := pass.CheckPbkdf2(password, passwordHash)
	if err != nil {
		return types.OID{}, wrap.IfErr(op, err)
	}
	if !ok {
		return types.OID{}, wrap.IfErr(op, ErrInvalidCredentials)
	}

	return playerOID, nil
}

func (s *Service) createPlayer(ctx context.Context, username, password string) (types.OID, error) {
	const op = "auth.Service.createPlayer"

	passwordHash, err := pass.MakePbkdf2(password)
	if err != nil {
		return types.OID{}, wrap.IfErr(op, err)
	}

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return types.OID{}, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	playerID, playerOID, err := s.insertPlayer(ctx, tx, username, passwordHash)
	if err != nil {
		return types.OID{}, wrap.IfErr(op, err)
	}
	if err = s.insertPlayerAvatars(ctx, tx, playerID, playerOID); err != nil {
		return types.OID{}, wrap.IfErr(op, err)
	}
	if err = s.setRandomActiveAvatar(ctx, tx, playerID); err != nil {
		return types.OID{}, wrap.IfErr(op, err)
	}
	if err = s.insertPlayerAvatarOutfits(ctx, tx, playerID); err != nil {
		return types.OID{}, wrap.IfErr(op, err)
	}
	if err = s.insertPlayerItems(ctx, tx, playerID); err != nil {
		return types.OID{}, wrap.IfErr(op, err)
	}

	if err = tx.Commit(); err != nil {
		return types.OID{}, wrap.IfErr(op, err)
	}

	return playerOID, nil
}

func (s *Service) insertPlayer(ctx context.Context, tx *sql.Tx, username, passwordHash string) (int, types.OID, error) {
	const op = "auth.Service.insertPlayer"

	var playerID int
	var playerOID types.OID
	err := tx.QueryRowContext(ctx, `
			insert into player (gsfoid, username, password, is_tutorial_completed, is_qa, max_outfits)
			values ((select coalesce(max(gsfoid), 0) + 1 from player), ?, ?, 1, 1, 1)
			returning id, gsfoid;
		`, username, passwordHash).Scan(&playerID, &playerOID)
	if s.store.IsErrConstraintUnique(err) {
		return 0, types.OID{}, ErrCreatePlayerAccountFailed
	}
	if err != nil {
		return 0, types.OID{}, wrap.IfErr(op, err)
	}
	return playerID, playerOID, nil
}

func (s *Service) insertPlayerAvatars(ctx context.Context, tx *sql.Tx, playerID int, playerOID types.OID) error {
	const op = "auth.Service.insertPlayerAvatars"

	result, err := tx.ExecContext(ctx, `
			insert into player_avatar (gsfoid, player_id, avatar_id, name, outfit_no, is_active)
			select
				(select coalesce(max(gsfoid), 0) from player_avatar) + row_number() over (order by av.id),
				?,
				av.id,
				av.name || ' ' || ?,
				1,
				0
			from avatar as av;
		`, playerID, playerOID.Int64())
	if s.store.IsErrConstraintUnique(err) {
		return ErrCreatePlayerAccountFailed
	}
	if err != nil {
		return wrap.IfErr(op, err)
	}
	return requireRowsAffected(op, result)
}

func (s *Service) setRandomActiveAvatar(ctx context.Context, tx *sql.Tx, playerID int) error {
	const op = "auth.Service.setRandomActiveAvatar"

	result, err := tx.ExecContext(ctx, `
			update player_avatar
			set is_active = 1
			where id = (
				select id
				from player_avatar
				where player_id = ?
				order by random()
				limit 1
			);
		`, playerID)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	return requireRowsAffected(op, result)
}

func (s *Service) insertPlayerAvatarOutfits(ctx context.Context, tx *sql.Tx, playerID int) error {
	const op = "auth.Service.insertPlayerAvatarOutfits"

	result, err := tx.ExecContext(ctx, `
			insert into player_avatar_outfit (gsfoid, player_avatar_id, outfit_no)
			select
				(select coalesce(max(gsfoid), 0) from player_avatar_outfit) + row_number() over (order by pa.id),
				pa.id,
				1
			from player_avatar as pa
			where pa.player_id = ?;
		`, playerID)
	if s.store.IsErrConstraintUnique(err) {
		return ErrCreatePlayerAccountFailed
	}
	if err != nil {
		return wrap.IfErr(op, err)
	}
	return requireRowsAffected(op, result)
}

func (s *Service) insertPlayerItems(ctx context.Context, tx *sql.Tx, playerID int) error {
	const op = "auth.Service.insertPlayerItems"

	result, err := tx.ExecContext(ctx, `
			insert into player_item (gsfoid, player_id, item_id, quantity)
			select
				(select coalesce(max(gsfoid), 0) from player_item) + row_number() over (order by it.id),
				?,
				it.id,
				1
			from item as it;
		`, playerID)
	if s.store.IsErrConstraintUnique(err) {
		return ErrCreatePlayerAccountFailed
	}
	if err != nil {
		return wrap.IfErr(op, err)
	}
	return requireRowsAffected(op, result)
}

func requireRowsAffected(op string, result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return wrap.IfErr(op, err)
	}
	if rowsAffected == 0 {
		return wrap.IfErr(op, ErrCreatePlayerAccountFailed)
	}
	return nil
}
