package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
)

func main() {
	mode := flag.String("mode", "", "Operation mode: import or export.")
	dbPath := flag.String("db", "", "Path to the blob.db database.")
	dir := flag.String("dir", "", "Cache source directory for import mode and destination directory for export mode.")
	overwrite := flag.Bool("overwrite", false, "Overwrite files on disk in export mode. Default: false.")

	flag.Parse()

	if *mode == "" || *dir == "" || *dbPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if err := run(logger, *mode, *dbPath, *dir, *overwrite); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func run(logger *slog.Logger, mode string, dbPath string, dir string, overwrite bool) error {
	store, err := db.NewSQLiteStore(dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer store.DB().Close()

	service := blob.NewService(logger, store, "")
	ctx := context.Background()

	switch mode {
	case "import":
		result, err := service.ImportFromFolder(ctx, dir)
		if err != nil {
			return err
		}
		fmt.Printf("imported: %d, skipped: %d\n", result.ImportedFiles, result.SkippedFiles)

	case "export":
		result, err := service.ExportToFolder(ctx, dir, overwrite)
		if err != nil {
			return err
		}
		fmt.Printf("exported: %d, skipped: %d\n", result.ExportedFiles, result.SkippedFiles)

	default:
		return fmt.Errorf("unknown mode: %s (expected import or export)", mode)
	}

	return nil
}
