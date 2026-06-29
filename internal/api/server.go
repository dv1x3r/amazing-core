package api

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/dv1x3r/amazing-core/data"
	"github.com/dv1x3r/amazing-core/internal/api/middleware"
	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/web"

	"github.com/dv1x3r/w2go/w2lib"
	"github.com/dv1x3r/w2go/w2widget"
)

type Server struct {
	logger *slog.Logger
	server *http.Server
}

func NewServer(logger *slog.Logger, handler *Handler, store db.Store) *Server {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", errorHandler(handler.GetDashboard))
	router.HandleFunc("POST /login", errorHandler(handler.PostLogin))
	router.HandleFunc("POST /logout", errorHandler(handler.PostLogout))
	router.HandleFunc("GET /cdn/{cdnid}", errorHandler(handler.GetBlob))

	router.Handle("GET /lib/", http.StripPrefix("/lib/", w2lib.FileServerFS()))
	router.Handle("GET /admin/", http.FileServerFS(web.FS))
	router.Handle("GET /queries/", http.FileServerFS(mustSubFS(data.FS, "sql")))
	router.Handle("GET /favicon.ico", serveFSFile(web.FS, "favicon_io/favicon.ico"))
	router.Handle("GET /site.webmanifest", serveFSFile(web.FS, "favicon_io/site.webmanifest"))
	router.Handle("GET /favicon-16x16.png", serveFSFile(web.FS, "favicon_io/favicon-16x16.png"))
	router.Handle("GET /favicon-32x32.png", serveFSFile(web.FS, "favicon_io/favicon-32x32.png"))
	router.Handle("GET /apple-touch-icon.png", serveFSFile(web.FS, "favicon_io/apple-touch-icon.png"))
	router.Handle("GET /android-chrome-192x192.png", serveFSFile(web.FS, "favicon_io/android-chrome-192x192.png"))
	router.Handle("GET /android-chrome-512x512.png", serveFSFile(web.FS, "favicon_io/android-chrome-512x512.png"))

	v1 := http.NewServeMux()

	v1.HandleFunc("GET /asset", errorHandler(handler.GetAsset))
	v1.HandleFunc("GET /asset/grid", errorHandler(handler.GetAssetGrid))
	v1.HandleFunc("POST /asset/grid", errorHandler(handler.PostAssetGrid))
	v1.HandleFunc("POST /asset/cache.json", errorHandler(handler.PostAssetCacheJSON))

	v1.HandleFunc("GET /asset/filetype", errorHandler(handler.GetAssetFileType))
	v1.HandleFunc("GET /asset/assettype", errorHandler(handler.GetAssetType))
	v1.HandleFunc("GET /asset/assetgroup", errorHandler(handler.GetAssetGroup))

	v1.HandleFunc("GET /avatar", errorHandler(handler.GetAvatar))
	v1.HandleFunc("GET /avatar/grid", errorHandler(handler.GetAvatarGrid))
	v1.HandleFunc("POST /avatar/form", errorHandler(handler.PostAvatarForm))
	v1.HandleFunc("POST /avatar/remove", errorHandler(handler.PostAvatarRemove))

	v1.HandleFunc("GET /avatar/slot", errorHandler(handler.GetAvatarSlot))

	v1.HandleFunc("GET /blob/grid", errorHandler(handler.GetBlobGrid))
	v1.HandleFunc("POST /blob/remove", errorHandler(handler.PostBlobRemove))
	v1.HandleFunc("POST /blob/import", errorHandler(handler.PostBlobImport))
	v1.HandleFunc("POST /blob/extract", errorHandler(handler.PostBlobExtract))

	v1.HandleFunc("GET /container", errorHandler(handler.GetContainer))
	v1.HandleFunc("GET /container/grid", errorHandler(handler.GetContainerGrid))
	v1.HandleFunc("POST /container/form", errorHandler(handler.PostContainerForm))
	v1.HandleFunc("POST /container/remove", errorHandler(handler.PostContainerRemove))

	v1.HandleFunc("GET /container/{id}/asset/grid", errorHandler(handler.GetContainerAssetGrid))
	v1.HandleFunc("POST /container/{id}/asset/form", errorHandler(handler.PostContainerAssetForm))
	v1.HandleFunc("POST /container/asset/grid", errorHandler(handler.PostContainerAssetGrid))
	v1.HandleFunc("POST /container/asset/remove", errorHandler(handler.PostContainerAssetRemove))
	v1.HandleFunc("POST /container/asset/reorder", errorHandler(handler.PostContainerAssetReorder))

	v1.HandleFunc("GET /container/{id}/package/grid", errorHandler(handler.GetContainerPackageGrid))
	v1.HandleFunc("POST /container/{id}/package/form", errorHandler(handler.PostContainerPackageForm))
	v1.HandleFunc("POST /container/package/grid", errorHandler(handler.PostContainerPackageGrid))
	v1.HandleFunc("POST /container/package/remove", errorHandler(handler.PostContainerPackageRemove))
	v1.HandleFunc("POST /container/package/reorder", errorHandler(handler.PostContainerPackageReorder))

	v1.HandleFunc("GET /item", errorHandler(handler.GetItem))
	v1.HandleFunc("GET /item/grid", errorHandler(handler.GetItemGrid))
	v1.HandleFunc("POST /item/form", errorHandler(handler.PostItemForm))
	v1.HandleFunc("POST /item/remove", errorHandler(handler.PostItemRemove))

	v1.HandleFunc("GET /item/category", errorHandler(handler.GetItemCategory))
	v1.HandleFunc("GET /item/category/grid", errorHandler(handler.GetItemCategoryGrid))
	v1.HandleFunc("POST /item/category/form", errorHandler(handler.PostItemCategoryForm))
	v1.HandleFunc("POST /item/category/remove", errorHandler(handler.PostItemCategoryRemove))

	v1.HandleFunc("GET /player/grid", errorHandler(handler.GetPlayerGrid))
	v1.HandleFunc("GET /player/form", errorHandler(handler.GetPlayerForm))
	v1.HandleFunc("POST /player/form", errorHandler(handler.PostPlayerForm))

	v1.HandleFunc("GET /player/{id}/item", errorHandler(handler.GetPlayerItem))
	v1.HandleFunc("GET /player/{id}/item/grid", errorHandler(handler.GetPlayerItemGrid))
	v1.HandleFunc("POST /player/{id}/item/form", errorHandler(handler.PostPlayerItemForm))
	v1.HandleFunc("POST /player/item/remove", errorHandler(handler.PostPlayerItemRemove))

	v1.HandleFunc("GET /player/{id}/avatar", errorHandler(handler.GetPlayerAvatar))
	v1.HandleFunc("GET /player/{id}/avatar/grid", errorHandler(handler.GetPlayerAvatarGrid))
	v1.HandleFunc("POST /player/{id}/avatar/form", errorHandler(handler.PostPlayerAvatarForm))
	v1.HandleFunc("POST /player/avatar/remove", errorHandler(handler.PostPlayerAvatarRemove))

	v1.HandleFunc("GET /player/{id}/outfit", errorHandler(handler.GetPlayerOutfit))
	v1.HandleFunc("GET /player/{id}/outfit/grid", errorHandler(handler.GetPlayerOutfitGrid))
	v1.HandleFunc("POST /player/outfit/form", errorHandler(handler.PostPlayerOutfitForm))
	v1.HandleFunc("POST /player/outfit/remove", errorHandler(handler.PostPlayerOutfitRemove))

	v1.HandleFunc("GET /randname/grid", errorHandler(handler.GetRandnameGrid))
	v1.HandleFunc("POST /randname/form", errorHandler(handler.PostRandnameForm))
	v1.HandleFunc("POST /randname/remove", errorHandler(handler.PostRandnameRemove))

	v1.HandleFunc("GET /siteframe/grid", errorHandler(handler.GetSiteFrameGrid))
	v1.HandleFunc("POST /siteframe/form", errorHandler(handler.PostSiteFrameForm))
	v1.HandleFunc("POST /siteframe/remove", errorHandler(handler.PostSiteFrameRemove))

	v1.HandleFunc("GET /zone/grid", errorHandler(handler.GetZoneGrid))
	v1.HandleFunc("POST /zone/form", errorHandler(handler.PostZoneForm))
	v1.HandleFunc("POST /zone/remove", errorHandler(handler.PostZoneRemove))

	v1.HandleFunc("GET /logs/watch", errorHandler(handler.GetLogsWatch))

	if config.Get().Storage.SQLExplorer {
		v1.HandleFunc("GET /sql", errorHandler(w2widget.SQLiteSchemaHandler(store.DB())))
		v1.HandleFunc("POST /sql", errorHandler(w2widget.SQLExecHandler(store.DB())))
	}

	protected := middleware.Protected(handler.svc.Auth)
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", protected(v1)))

	stack := middleware.CreateStack(
		middleware.Secure(),
		middleware.IPExtractor(),
		middleware.Logger(logger),
		middleware.Recover(),
		middleware.RateLimiter(200, 300, 3*time.Minute),
	)

	server := &http.Server{
		Handler:           stack(router),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return &Server{
		logger: logger,
		server: server,
	}
}

func (s *Server) ListenAndServe(address string) {
	s.server.Addr = address
	s.logger.Info("starting the api server on " + address)
	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(err.Error())
		}
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	s.logger.Info("shutting down the api server")
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(err.Error())
	}
}

func serveFSFile(fsys fs.FS, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, fsys, name)
	}
}

func mustSubFS(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
