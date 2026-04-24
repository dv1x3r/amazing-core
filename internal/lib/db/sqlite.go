package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"os/exec"

	"github.com/pressly/goose/v3"
)

type SQLiteStore struct {
	filePath string
	sqlDB    *sql.DB
}

func (s *SQLiteStore) DB() *sql.DB {
	return s.sqlDB
}

func (s *SQLiteStore) DumpTable(ctx context.Context, tableName string) ([]byte, error) {
	return exec.CommandContext(ctx, "sqlite3", s.filePath, ".dump "+tableName).Output()
}

func (s *SQLiteStore) MigrateUp(logger *slog.Logger, fsys fs.FS, dir string) error {
	goose.SetLogger(Logger{logger})
	goose.SetBaseFS(fsys)
	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}
	return goose.Up(s.sqlDB, dir)
}

func (s *SQLiteStore) MigrateFile(logger *slog.Logger, fsys fs.FS, fileName string) error {
	ctx := context.Background()

	conn, err := s.sqlDB.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	var isInitialized bool
	row := conn.QueryRowContext(ctx, "select exists (select 1 from sqlite_master where type='table' and name not like 'sqlite_%');")
	if err := row.Scan(&isInitialized); err != nil {
		return err
	}

	if isInitialized {
		logger.Info(fmt.Sprintf("skipping the %s migration: already initialized", fileName))
		return nil
	}

	data, err := fs.ReadFile(fsys, fileName)
	if err != nil {
		return err
	}

	conn.ExecContext(ctx, "PRAGMA foreign_keys = OFF;")
	defer conn.ExecContext(ctx, "PRAGMA foreign_keys = ON;")

	logger.Info(fmt.Sprintf("applying the %s migration", fileName))
	_, err = conn.ExecContext(ctx, string(data))
	return err
}

func (s *SQLiteStore) RecreateTableFromFile(logger *slog.Logger, fsys fs.FS, fileName string, tableName string, overwrite bool) error {
	data, err := fs.ReadFile(fsys, fileName)
	if err != nil {
		return err
	}
	return s.RecreateTableFromQuery(logger, string(data), tableName, overwrite)
}

func (s *SQLiteStore) RecreateTableFromQuery(logger *slog.Logger, query string, tableName string, overwrite bool) error {
	ctx := context.Background()

	conn, err := s.sqlDB.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if !overwrite {
		var isInitialized bool
		row := conn.QueryRowContext(ctx, "select exists (select 1 from sqlite_master where type='table' and name = ?);", tableName)
		if err := row.Scan(&isInitialized); err != nil {
			return err
		}
		if isInitialized {
			var count int
			row := conn.QueryRowContext(ctx, fmt.Sprintf("select count(*) from %q;", tableName))
			if err := row.Scan(&count); err != nil {
				return err
			}
			if count > 0 {
				logger.Info(fmt.Sprintf("skipping the %s table migration: already initialized", tableName))
				return nil
			}
		}
	}

	conn.ExecContext(ctx, "PRAGMA foreign_keys = OFF;")
	defer conn.ExecContext(ctx, "PRAGMA foreign_keys = ON;")

	logger.Warn(fmt.Sprintf("preparing the %s table migration: dropping the table", tableName))
	if _, err := conn.ExecContext(ctx, fmt.Sprintf("drop table if exists %q;", tableName)); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("applying the %s table migration", tableName))
	_, err = conn.ExecContext(ctx, query)
	return err
}
