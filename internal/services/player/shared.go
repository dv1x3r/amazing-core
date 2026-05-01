package player

import (
	"context"
	"database/sql"
)

func setActivePlayerAvatarByID(ctx context.Context, tx *sql.Tx, playerAvatarID int, playerID int) error {
	_, err := tx.ExecContext(ctx, `
		update player_avatar
		set is_active = case when id = ? then 1 else 0 end
		where player_id = ?;
	`, playerAvatarID, playerID)
	return err
}
