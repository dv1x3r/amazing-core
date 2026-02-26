package api

import (
	"context"
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/dv1x3r/amazing-core/internal/api/middleware"
	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/web"
	"github.com/dv1x3r/w2go/w2"

	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
	"github.com/dv1x3r/amazing-core/internal/services/randomnames"
)

//go:embed *.gotmpl
var templatesFS embed.FS

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFS(templatesFS, "*.gotmpl"))
}

type Server struct {
	logger *slog.Logger
	server *http.Server
}

func NewServer(
	logger *slog.Logger,
	authService *auth.Service,
	blobService *blob.Service,
	randomNamesService *randomnames.Service,
) *Server {
	router := http.NewServeMux()

	router.Handle("GET /{$}", adminHandler(authService))
	router.Handle("GET /admin/", fsFolderHandler(web.FS, "admin", "/admin"))
	router.Handle("GET /favicon.ico", fsFileHandler(web.FS, "favicon_io/favicon.ico"))
	router.Handle("GET /site.webmanifest", fsFileHandler(web.FS, "favicon_io/site.webmanifest"))
	router.Handle("GET /favicon-16x16.png", fsFileHandler(web.FS, "favicon_io/favicon-16x16.png"))
	router.Handle("GET /favicon-32x32.png", fsFileHandler(web.FS, "favicon_io/favicon-32x32.png"))
	router.Handle("GET /apple-touch-icon.png", fsFileHandler(web.FS, "favicon_io/apple-touch-icon.png"))
	router.Handle("GET /android-chrome-192x192.png", fsFileHandler(web.FS, "favicon_io/android-chrome-192x192.png"))
	router.Handle("GET /android-chrome-512x512.png", fsFileHandler(web.FS, "favicon_io/android-chrome-512x512.png"))

	authHandler := auth.NewAPIHandler(authService)
	router.HandleFunc("POST /login", errorHandler(authHandler.PostLogin))
	router.HandleFunc("POST /logout", errorHandler(authHandler.PostLogout))

	v1 := http.NewServeMux()

	blobHandler := blob.NewAPIHandler(blobService)
	v1.HandleFunc("GET /blob/records", errorHandler(blobHandler.GetRecords))
	v1.HandleFunc("POST /blob/upload", errorHandler(blobHandler.PostUpload))
	v1.HandleFunc("POST /blob/remove", errorHandler(blobHandler.PostRemove))
	v1.HandleFunc("POST /blob/import", errorHandler(blobHandler.PostImport))
	v1.HandleFunc("POST /blob/export", errorHandler(blobHandler.PostExport))
	v1.HandleFunc("POST /blob/s3sync", errorHandler(blobHandler.PostS3Sync))
	if config.Get().Settings.AssetDeliveryAPI {
		router.HandleFunc("GET /cdn/{cdnid}", errorHandler(blobHandler.GetBlob))
	}

	randomNamesHandler := randomnames.NewAPIHandler(randomNamesService)
	v1.HandleFunc("GET /randomnames/form", errorHandler(randomNamesHandler.GetForm))
	v1.HandleFunc("POST /randomnames/form", errorHandler(randomNamesHandler.PostForm))
	v1.HandleFunc("GET /randomnames/records", errorHandler(randomNamesHandler.GetRecords))
	v1.HandleFunc("POST /randomnames/remove", errorHandler(randomNamesHandler.PostRemove))

	protected := middleware.Protected(authService)
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", protected(v1)))

	cop := http.NewCrossOriginProtection()
	for _, origin := range config.Get().Secure.CSRF.TrustedOrigins {
		cop.AddTrustedOrigin(origin)
	}

	stack := middleware.CreateStack(
		middleware.Secure(),
		middleware.IPExtractor(),
		middleware.Logger(logger),
		middleware.Recover(),
		middleware.RateLimiter(50, 100, 3*time.Minute),
		cop.Handler,
	)

	server := &http.Server{
		Handler:      stack(router),
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
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

func adminHandler(authService *auth.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, ok := authService.GetSessionUsername(w, r)
		if ok {
			if err := authService.RefreshSession(w, r); err != nil {
				w2.NewErrorResponse(err.Error()).Write(w, http.StatusInternalServerError)
				return
			}
		}

		data := map[string]any{"username": username}
		if err := tmpl.ExecuteTemplate(w, "admin.gotmpl", data); err != nil {
			w2.NewErrorResponse(err.Error()).Write(w, http.StatusInternalServerError)
		}
	})
}

func fsFileHandler(fsys fs.FS, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, fsys, name)
	}
}

func fsFolderHandler(fsys fs.FS, dir string, strip string) http.Handler {
	subfs, _ := fs.Sub(fsys, dir)
	handler := http.FileServerFS(subfs)
	handler = http.StripPrefix(strip, handler)
	return handler
}

func errorHandler(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		if errors.Is(err, context.Canceled) {
			res := w2.NewErrorResponse("Client Closed Request")
			res.Write(w, 499)
			return
		}

		if e, ok := r.Context().Value("err").(*error); ok {
			*e = err
		}

		status := wrap.HTTPStatus(err)
		message := err.Error()

		// do not expose details of 500 errors
		// but I do not care about that right now
		// if status >= 500 {
		// 	message = http.StatusText(status)
		// }

		res := w2.NewErrorResponse(message)
		res.Write(w, status)
	}
}
