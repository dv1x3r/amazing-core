package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"os/exec"
	"sync/atomic"

	"github.com/pressly/goose/v3"
)

var sqliteDriverID atomic.Uint64

func nextSQLiteDriverName(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, sqliteDriverID.Add(1))
}

type SQLiteStore struct {
	logger   *slog.Logger
	sqlDB    *sql.DB
	filePath string
	attached map[string]string
}

func NewSQLiteStore(logger *slog.Logger, filePath string, opts ...SQLiteStoreOption) Store {
	store := &SQLiteStore{
		logger:   logger,
		filePath: filePath,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(store)
	}
	return store
}

type SQLiteStoreOption func(*SQLiteStore)

func WithAttach(name string, filePath string) SQLiteStoreOption {
	return func(s *SQLiteStore) {
		if s.attached == nil {
			s.attached = make(map[string]string)
		}
		s.attached[name] = filePath
	}
}

func (s *SQLiteStore) DB() *sql.DB {
	return s.sqlDB
}

func (s *SQLiteStore) Open() error {
	if s.sqlDB != nil {
		return nil
	}

	s.logger.Info(fmt.Sprintf("connecting to the %s using %s driver", s.filePath, s.DriverName()))

	var sqlDB *sql.DB
	var err error

	if len(s.attached) > 0 {
		for k, v := range s.attached {
			s.logger.Info(fmt.Sprintf("attaching %s AS %s database", v, k))
		}
		if sqlDB, err = s.openWithAttachedDBs(); err != nil {
			return err
		}
	} else {
		if sqlDB, err = s.openSingleDB(); err != nil {
			return err
		}
	}

	if err = sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return err
	}

	s.sqlDB = sqlDB
	return nil
}

func (s *SQLiteStore) Close() error {
	if s.sqlDB == nil {
		return nil
	}
	err := s.sqlDB.Close()
	s.sqlDB = nil
	return err
}

func (s *SQLiteStore) DumpTable(ctx context.Context, tableName string) ([]byte, error) {
	return exec.CommandContext(ctx, "sqlite3", s.filePath, ".dump "+tableName).Output()
}

func (s *SQLiteStore) MigrateUp(fsys fs.FS, dir string) error {
	goose.SetLogger(Logger{s.logger})
	goose.SetBaseFS(fsys)
	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}
	return goose.Up(s.sqlDB, dir)
}

func (s *SQLiteStore) MigrateFile(fsys fs.FS, fileName string) error {
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
		s.logger.Info(fmt.Sprintf("skipping the %s migration: already initialized", fileName))
		return nil
	}

	data, err := fs.ReadFile(fsys, fileName)
	if err != nil {
		return err
	}

	conn.ExecContext(ctx, "PRAGMA foreign_keys = OFF;")
	defer conn.ExecContext(ctx, "PRAGMA foreign_keys = ON;")

	s.logger.Info(fmt.Sprintf("applying the %s migration", fileName))
	_, err = conn.ExecContext(ctx, string(data))
	return err
}

func (s *SQLiteStore) RecreateTableFromFile(fsys fs.FS, fileName string, tableName string, overwrite bool) error {
	data, err := fs.ReadFile(fsys, fileName)
	if err != nil {
		return err
	}
	return s.RecreateTableFromQuery(string(data), tableName, overwrite)
}

func (s *SQLiteStore) RecreateTableFromQuery(query string, tableName string, overwrite bool) error {
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
				s.logger.Info(fmt.Sprintf("skipping the %s table migration: already initialized", tableName))
				return nil
			}
		}
	}

	conn.ExecContext(ctx, "PRAGMA foreign_keys = OFF;")
	defer conn.ExecContext(ctx, "PRAGMA foreign_keys = ON;")

	s.logger.Warn(fmt.Sprintf("preparing the %s table migration: dropping the table", tableName))
	if _, err := conn.ExecContext(ctx, fmt.Sprintf("drop table if exists %q;", tableName)); err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("applying the %s table migration", tableName))
	_, err = conn.ExecContext(ctx, query)
	return err
}
