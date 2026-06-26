//go:build !cgo

package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"modernc.org/sqlite"
	sqliteLib "modernc.org/sqlite/lib"
)

const dataSourceArgs = "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)&_pragma=busy_timeout(10000)"

func (s SQLiteStore) openSingleDB() (*sql.DB, error) {
	return sql.Open("sqlite", s.filePath+dataSourceArgs)
}

func (s SQLiteStore) openWithAttachedDBs() (*sql.DB, error) {
	driverName := nextSQLiteDriverName("sqlite")
	sqliteDriver := &sqlite.Driver{}
	sqliteDriver.RegisterConnectionHook(
		func(conn sqlite.ExecQuerierContext, dsn string) error {
			for name, filePath := range s.attached {
				query := fmt.Sprintf("ATTACH DATABASE ? AS %s;", name)
				args := []driver.NamedValue{{Ordinal: 1, Value: filePath}}
				if _, err := conn.ExecContext(context.Background(), query, args); err != nil {
					return fmt.Errorf("attach database %q: %w", name, err)
				}
			}
			return nil
		},
	)
	sql.Register(driverName, sqliteDriver)
	return sql.Open(driverName, s.filePath+dataSourceArgs)
}

func (SQLiteStore) DriverName() string {
	return "nocgo-sqlite"
}

func (SQLiteStore) IsErrConstraintUnique(err error) bool {
	var sqliteErr *sqlite.Error
	ok := errors.As(err, &sqliteErr)
	return ok && sqliteErr.Code() == sqliteLib.SQLITE_CONSTRAINT_UNIQUE
}

func (SQLiteStore) IsErrConstraintTrigger(err error) bool {
	var sqliteErr *sqlite.Error
	ok := errors.As(err, &sqliteErr)
	return ok && sqliteErr.Code() == sqliteLib.SQLITE_CONSTRAINT_TRIGGER
}

func (SQLiteStore) IsErrConstraintForeignKey(err error) bool {
	var sqliteErr *sqlite.Error
	ok := errors.As(err, &sqliteErr)
	return ok && sqliteErr.Code() == sqliteLib.SQLITE_CONSTRAINT_FOREIGNKEY
}
