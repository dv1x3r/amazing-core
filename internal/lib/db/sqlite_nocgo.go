//go:build !cgo

package db

import (
	"database/sql"
	"errors"

	"modernc.org/sqlite"
	sqliteLib "modernc.org/sqlite/lib"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(filePath string) (Store, error) {
	const args = "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)&_pragma=busy_timeout(10000)"
	db, err := sql.Open("sqlite", filePath+args)
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) DB() *sql.DB {
	return s.db
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
