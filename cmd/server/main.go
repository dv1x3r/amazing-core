package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/dv1x3r/amazing-core/data"
	"github.com/dv1x3r/w2go/w2db"
	"github.com/huandu/go-sqlbuilder"

	"github.com/dv1x3r/amazing-core/internal/api"
	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/game"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/downloader"
	"github.com/dv1x3r/amazing-core/internal/lib/logger"

	"github.com/dv1x3r/amazing-core/internal/services/asset"
	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
	"github.com/dv1x3r/amazing-core/internal/services/dummy"
	"github.com/dv1x3r/amazing-core/internal/services/randname"
	"github.com/dv1x3r/amazing-core/internal/services/siteframe"
)

var (
	version = "nightly build"
)

const AMAZING_CORE = `
           ___   __  ______ ____  _____  _______
          / _ | /  |/  / _ /_  / /  _/ |/ / ___/
         / __ |/ /|_/ / __ |/ /__/ //    / (_ / 
        /_/ |_/_/  /_/_/ |_/___/___/_/|_/\___/  
                        / ___/ __ \/ _ \/ __/   
                       / /__/ /_/ / , _/ _/     
                       \___/\____/_/|_/___/     
  Amazing Core %s - https://amazingcore.org/
`

func main() {
	cfg := config.Get()

	// ── Logger ──────────────────────────────────────────────────────────────────
	switch cfg.Logger.Level {
	case "debug+sql":
		logger.SetLevel(slog.LevelDebug)
	case "debug":
		logger.SetLevel(slog.LevelDebug)
	case "info":
		logger.SetLevel(slog.LevelInfo)
	case "error":
		logger.SetLevel(slog.LevelError)
	case "warn":
		logger.SetLevel(slog.LevelWarn)
	default:
		logger.SetLevel(slog.LevelDebug)
	}

	switch cfg.Logger.Handler {
	case "text":
		logger.Create(logger.TextHandler)
	case "json":
		logger.Create(logger.JsonHandler)
	case "pretty":
		logger.Create(logger.PrettyHandler)
		logger.Get().Info(fmt.Sprintf(AMAZING_CORE, version), "config", cfg)
	default:
		logger.Create(logger.PrettyHandler)
		logger.Get().Info(fmt.Sprintf(AMAZING_CORE, version), "config", cfg)
	}

	w2db.SetFlavor(sqlbuilder.SQLite)
	if cfg.Logger.Level == "debug+sql" {
		w2db.SetLogger(logger.Get())
	}

	// ── Config check ────────────────────────────────────────────────────────────
	if cfg.Servers.API == "" || cfg.Servers.Game == "" {
		logger.Get().Error("missing server configuration", "api", cfg.Servers.API, "game", cfg.Servers.Game)
		fmt.Scanln()
		os.Exit(1)
	}

	if cfg.Storage.Databases.Core == "" || cfg.Storage.Databases.Blob == "" {
		logger.Get().Error("missing database configuration", "core", cfg.Storage.Databases.Core, "blob", cfg.Storage.Databases.Blob)
		fmt.Scanln()
		os.Exit(1)
	}

	// ── Database folders ────────────────────────────────────────────────────────
	if err := os.MkdirAll(path.Dir(cfg.Storage.Databases.Core), os.ModePerm); err != nil {
		logger.Get().Error("unable to access the folder for the core database", "err", err)
		fmt.Scanln()
		os.Exit(1)
	}

	if err := os.MkdirAll(path.Dir(cfg.Storage.Databases.Blob), os.ModePerm); err != nil {
		logger.Get().Error("unable to access the folder for the blob database", "err", err)
		fmt.Scanln()
		os.Exit(1)
	}

	// ── Download blob.db ────────────────────────────────────────────────────────
	if cfg.Blob.Download {
		if err := downloader.DownloadIfNotExists(logger.Get(), cfg.Storage.Databases.Blob, cfg.Blob.DownloadURL); err != nil {
			logger.Get().Error("unable to download blob.db", "err", err)
			fmt.Scanln()
			os.Exit(1)
		}
	}

	// ── Database stores ─────────────────────────────────────────────────────────
	coreStore, err := db.NewSQLiteStore(cfg.Storage.Databases.Core)
	if err != nil {
		logger.Get().Error("unable to connect to core.db", "err", err)
		fmt.Scanln()
		os.Exit(1)
	}
	defer coreStore.DB().Close()

	logger.Get().Info(fmt.Sprintf("connected to the %s using the %s driver", cfg.Storage.Databases.Core, coreStore.DriverName()))

	blobStore, err := db.NewSQLiteStore(cfg.Storage.Databases.Blob)
	if err != nil {
		logger.Get().Error("unable to connect to blob.db", "err", err)
		fmt.Scanln()
		os.Exit(1)
	}
	defer blobStore.DB().Close()

	logger.Get().Info(fmt.Sprintf("connected to the %s using the %s driver", cfg.Storage.Databases.Blob, blobStore.DriverName()))

	// ── Database migrations ─────────────────────────────────────────────────────
	if err := coreStore.MigrateBaseFile(logger.Get(), data.FS, "sql/core_db/base.sql"); err != nil {
		logger.Get().Error("unable to initialize core.db", "err", err)
		fmt.Scanln()
		os.Exit(1)
	}

	if err := blobStore.MigrateBaseFile(logger.Get(), data.FS, "sql/blob_db/base.sql"); err != nil {
		logger.Get().Error("unable to initialize blob.db", "err", err)
		fmt.Scanln()
		os.Exit(1)
	}

	if err := coreStore.MigrateUp(logger.Get(), data.FS, "sql/core_db/updates"); err != nil {
		logger.Get().Error("unable to apply core.db updates", "err", err)
		fmt.Scanln()
		os.Exit(1)
	}

	// ── Services & Servers ──────────────────────────────────────────────────────
	authService := auth.NewService(cfg.Secure.Session.Secure, cfg.Secure.Session.Key, cfg.Secure.Auth.Username, cfg.Secure.Auth.Password)
	dummyService := dummy.NewService(coreStore)
	blobService := blob.NewService(logger.Get(), blobStore, cfg.Settings.AssetDeliveryURL)
	assetService := asset.NewService(logger.Get(), coreStore, cfg.Settings.AssetDeliveryURL)
	randnameService := randname.NewService(coreStore)
	siteFrameService := siteframe.NewService(logger.Get(), coreStore)

	apiHandler := api.NewHandler(
		authService,
		assetService,
		blobService,
		siteFrameService,
		dummyService,
		randnameService,
	)

	apiServer := api.NewServer(
		logger.Get(),
		coreStore,
		apiHandler,
		authService,
	)

	gameHandler := game.NewHandler(
		dummyService,
		assetService,
		randnameService,
	)

	gameServer := game.NewServer(
		logger.Get(),
		gameHandler,
	)

	interruptCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		apiServer.ListenAndServe(cfg.Servers.API)
		stop()
	}()

	go func() {
		gameServer.ListenAndServe(cfg.Servers.Game)
		stop()
	}()

	<-interruptCtx.Done()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	apiServer.Shutdown(shutdownCtx)
	gameServer.Shutdown(shutdownCtx)
}
