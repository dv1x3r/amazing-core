package db

import (
	"context"
	"database/sql"
	"io/fs"
)

type Store interface {
	Open() error
	Close() error
	DB() *sql.DB
	DriverName() string
	DumpTable(ctx context.Context, tableName string) ([]byte, error)
	MigrateUp(fsys fs.FS, dir string) error
	MigrateFile(fsys fs.FS, fileName string) error
	RecreateTableFromFile(fsys fs.FS, fileName string, tableName string, overwrite bool) error
	RecreateTableFromQuery(query string, tableName string, overwrite bool) error
	IsErrConstraintUnique(err error) bool
	IsErrConstraintTrigger(err error) bool
	IsErrConstraintForeignKey(err error) bool
}
