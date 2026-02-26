package db

import (
	"context"
	"database/sql"
	"io/fs"
	"log/slog"
)

type Store interface {
	DB() *sql.DB
	DriverName() string
	DumpTable(ctx context.Context, tableName string) ([]byte, error)
	MigrateUp(logger *slog.Logger, fsys fs.FS, dir string) error
	MigrateBaseFile(logger *slog.Logger, fsys fs.FS, fileName string) error
	RecreateTableFromFile(logger *slog.Logger, fsys fs.FS, fileName string, tableName string, overwrite bool) error
	RecreateTableFromQuery(logger *slog.Logger, query string, tableName string, overwrite bool) error
	IsErrConstraintUnique(err error) bool
	IsErrConstraintTrigger(err error) bool
	IsErrConstraintForeignKey(err error) bool
}
