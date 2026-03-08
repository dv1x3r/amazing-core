package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/dv1x3r/amazing-core/data"

	"github.com/dv1x3r/amazing-core/internal/api"
	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/game"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/downloader"
	"github.com/dv1x3r/amazing-core/internal/lib/logger"

	"github.com/dv1x3r/amazing-core/internal/services/asset"
	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
	"github.com/dv1x3r/amazing-core/internal/services/randname"

	"github.com/gorilla/sessions"
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
  %s - github.com/dv1x3r/amazing-core
`

func main() {
	cfg := config.Get()

	switch cfg.Logger.Level {
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

	if cfg.Servers.API == "" || cfg.Servers.Game == "" {
		logger.Get().Error("missing server configuration", "api", cfg.Servers.API, "game", cfg.Servers.Game)
		os.Exit(1)
	}

	if cfg.Storage.Databases.Core == "" || cfg.Storage.Databases.Blob == "" {
		logger.Get().Error("missing database configuration", "core", cfg.Storage.Databases.Core, "blob", cfg.Storage.Databases.Blob)
		os.Exit(1)
	}

	// prepare blob.db if missing & download is enabled
	if cfg.Blob.Download {
		if err := downloader.DownloadIfNotExists(logger.Get(), cfg.Storage.Databases.Blob, cfg.Blob.DownloadURL); err != nil {
			logger.Get().Error("unable to download blob.db", "err", err)
			os.Exit(1)
		}
	}

	// prepare database folders

	if err := os.MkdirAll(path.Dir(cfg.Storage.Databases.Core), os.ModePerm); err != nil {
		logger.Get().Error("unable to access the folder for the core database", "err", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(path.Dir(cfg.Storage.Databases.Blob), os.ModePerm); err != nil {
		logger.Get().Error("unable to access the folder for the blob database", "err", err)
		os.Exit(1)
	}

	// connect to the databases

	coreStore, err := db.NewSQLiteStore(cfg.Storage.Databases.Core)
	if err != nil {
		logger.Get().Error("unable to connect to core.db", "err", err)
		os.Exit(1)
	}
	defer coreStore.DB().Close()

	logger.Get().Info(fmt.Sprintf("connected to the %s using the %s driver", cfg.Storage.Databases.Core, coreStore.DriverName()))

	blobStore, err := db.NewSQLiteStore(cfg.Storage.Databases.Blob)
	if err != nil {
		logger.Get().Error("unable to connect to blob.db", "err", err)
		os.Exit(1)
	}
	defer blobStore.DB().Close()

	logger.Get().Info(fmt.Sprintf("connected to the %s using the %s driver", cfg.Storage.Databases.Blob, blobStore.DriverName()))

	// apply base migrations

	if err := coreStore.MigrateBaseFile(logger.Get(), data.FS, "sql/core_db/base.sql"); err != nil {
		logger.Get().Error("unable to initialize core.db", "err", err)
		os.Exit(1)
	}

	if err := blobStore.MigrateBaseFile(logger.Get(), data.FS, "sql/blob_db/base.sql"); err != nil {
		logger.Get().Error("unable to initialize blob.db", "err", err)
		os.Exit(1)
	}

	// apply update migrations

	if err := coreStore.MigrateUp(logger.Get(), data.FS, "sql/core_db/updates"); err != nil {
		logger.Get().Error("unable to apply core.db updates", "err", err)
		os.Exit(1)
	}

	if err := coreStore.RecreateTableFromFile(logger.Get(), data.FS, "sql/core_db/assets.sql", "asset", false); err != nil {
		logger.Get().Error("unable to apply assets.sql", "err", err)
		os.Exit(1)
	}

	if _, err := asset.ImportCacheJSON(logger.Get(), coreStore.DB(), data.FS, "cache.json", false); err != nil {
		logger.Get().Error("unable to import cache.json", "err", err)
		os.Exit(1)
	}

	session := sessions.NewCookieStore([]byte(cfg.Secure.Session.Key))
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 14,
		Secure:   cfg.Secure.Session.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	authService := auth.NewService(session)
	blobService := blob.NewService(logger.Get(), blobStore)
	assetService := asset.NewService(logger.Get(), coreStore)
	randnameService := randname.NewService(coreStore)

	apiServer := api.NewServer(
		coreStore.DB(),
		logger.Get(),
		authService,
		blobService,
		assetService,
		randnameService,
	)

	gameServer := game.NewServer(
		logger.Get(),
		randnameService,
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
