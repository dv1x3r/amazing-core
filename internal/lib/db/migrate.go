package db

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"github.com/pressly/goose/v3"
)

type gLogger struct {
	slog.Logger
}

func (l gLogger) Printf(format string, v ...any) {
	l.Info(fmt.Sprintf(format, v...))
}

func (l gLogger) Fatalf(format string, v ...any) {
	l.Error(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func MigrateBase(logger *slog.Logger, db *sql.DB, fsys fs.FS, sqlFilePath string) error {
	var isInitialized bool
	row := db.QueryRow("select exists (select 1 from sqlite_master where type='table' and name not like 'sqlite_%');")
	if err := row.Scan(&isInitialized); err != nil {
		return err
	}

	if isInitialized {
		return nil
	}

	sqlBytes, err := fs.ReadFile(fsys, sqlFilePath)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("applying the base %s migration", sqlFilePath))
	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return err
	}

	return nil
}

func MigrateUp(logger *slog.Logger, db *sql.DB, fsys fs.FS, dir string) error {
	goose.SetLogger(gLogger{Logger: *logger})
	goose.SetBaseFS(fsys)

	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}

	if err := goose.Up(db, dir); err != nil {
		return err
	}

	return nil
}
