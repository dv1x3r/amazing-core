package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	sqldata "github.com/dv1x3r/amazing-core/data/sql"

	"github.com/dv1x3r/amazing-core/internal/api"
	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/game"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/logger"
	"github.com/dv1x3r/amazing-core/internal/lib/prettyslog"

	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
	"github.com/dv1x3r/amazing-core/internal/services/randomnames"

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

	logger.Set(slog.New(prettyslog.NewHandler(&slog.HandlerOptions{Level: slog.LevelDebug})))
	logger.Get().Info(fmt.Sprintf(AMAZING_CORE, version), "config", cfg)

	if cfg.Servers.API == "" || cfg.Servers.Game == "" {
		logger.Get().Error("missing server configuration", "api", cfg.Servers.API, "game", cfg.Servers.Game)
		os.Exit(1)
	}

	if cfg.Storage.Databases.Core == "" || cfg.Storage.Databases.Blob == "" {
		logger.Get().Error("missing database configuration", "core", cfg.Storage.Databases.Core, "blob", cfg.Storage.Databases.Blob)
		os.Exit(1)
	}

	if err := os.MkdirAll(path.Dir(cfg.Storage.Databases.Core), os.ModePerm); err != nil {
		logger.Get().Error("unable to access the folder for the core database", "err", err)
		os.Exit(1)
	}

	store, err := db.NewSQLiteStore(cfg.Storage.Databases.Core)
	if err != nil {
		logger.Get().Error("unable to connect to the core database", "err", err)
		os.Exit(1)
	}
	defer store.DB().Close()

	logger.Get().Info(fmt.Sprintf("connected to the %s using the %s driver", cfg.Storage.Databases.Core, store.DriverName()))
	logger.Get().Info(fmt.Sprintf("attaching the %s database", cfg.Storage.Databases.Blob))
	if _, err := store.DB().Exec("attach database ? as blob;", cfg.Storage.Databases.Blob); err != nil {
		logger.Get().Error("unable to attach the blob database", "err", err)
		os.Exit(1)
	}

	if err := db.MigrateBase(logger.Get(), store.DB(), sqldata.FS, "base/core_db.sql"); err != nil {
		logger.Get().Error("unable to initialize the core database", "err", err)
		os.Exit(1)
	}

	if err := db.MigrateUp(logger.Get(), store.DB(), sqldata.FS, "updates"); err != nil {
		logger.Get().Error("unable to migrate the core database", "err", err)
		os.Exit(1)
	}

	session := sessions.NewCookieStore([]byte(config.Get().Secure.Session.Key))
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 14,
		Secure:   config.Get().Secure.Session.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	authService := auth.NewService(session)
	blobService := blob.NewService(store)
	randomNamesService := randomnames.NewService(store)

	apiServer := api.NewServer(
		authService,
		blobService,
		randomNamesService,
	)

	gameServer := game.NewServer(
		randomNamesService,
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		err := apiServer.Start(cfg.Servers.API)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Get().Error(err.Error())
			stop()
		}
	}()

	go func() {
		err := gameServer.Start(cfg.Servers.Game)
		if err != nil && !errors.Is(err, net.ErrClosed) {
			logger.Get().Error(err.Error())
			stop()
		}
	}()

	<-ctx.Done()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	apiServer.Shutdown(shutdownCtx)
	gameServer.Shutdown(shutdownCtx)
}
