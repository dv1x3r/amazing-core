//go:build cgo

package db

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

func NewSQLiteStore(filePath string) (Store, error) {
	const args = "?_journal=WAL&_fk=1&_busy_timeout=10000"
	sqlDB, err := sql.Open("sqlite3", filePath+args)
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{filePath: filePath, sqlDB: sqlDB}, nil
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
