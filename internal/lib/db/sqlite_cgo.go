//go:build cgo

package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

const dataSourceArgs = "?_journal_mode=WAL&_foreign_keys=1&_busy_timeout=10000"

func (s SQLiteStore) openSingleDB() (*sql.DB, error) {
	return sql.Open("sqlite3", s.filePath+dataSourceArgs)
}

func (s SQLiteStore) openWithAttachedDBs() (*sql.DB, error) {
	driverName := nextSQLiteDriverName("sqlite3")
	sql.Register(driverName, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			for name, filePath := range s.attached {
				query := fmt.Sprintf("ATTACH DATABASE ? AS %s;", name)
				args := []driver.NamedValue{{Ordinal: 1, Value: filePath}}
				if _, err := conn.ExecContext(context.Background(), query, args); err != nil {
					return fmt.Errorf("attach database %q: %w", name, err)
				}
			}
			return nil
		},
	})
	return sql.Open(driverName, s.filePath+dataSourceArgs)
}

func (SQLiteStore) DriverName() string {
	return "cgo-sqlite3"
}

func (SQLiteStore) IsErrConstraintUnique(err error) bool {
	var sqliteErr sqlite3.Error
	ok := errors.As(err, &sqliteErr)
	return ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}

func (SQLiteStore) IsErrConstraintTrigger(err error) bool {
	var sqliteErr sqlite3.Error
	ok := errors.As(err, &sqliteErr)
	return ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintTrigger
}

func (SQLiteStore) IsErrConstraintForeignKey(err error) bool {
	var sqliteErr sqlite3.Error
	ok := errors.As(err, &sqliteErr)
	return ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey
}
