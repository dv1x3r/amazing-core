package api

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"time"

	"github.com/dv1x3r/amazing-core/internal/api/admin"
	"github.com/dv1x3r/amazing-core/internal/api/middleware"
	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/lib/logger"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/web"

	"github.com/dv1x3r/amazing-core/internal/services/auth"
	"github.com/dv1x3r/amazing-core/internal/services/blob"
	"github.com/dv1x3r/amazing-core/internal/services/randomnames"

	"github.com/dv1x3r/w2go/w2"
	"github.com/gorilla/csrf"
)

type Server struct {
	server *http.Server
}

func NewServer(
	authService *auth.Service,
	blobService *blob.Service,
	randomNamesService *randomnames.Service,
) *Server {
	router := http.NewServeMux()

	router.Handle("GET /{$}", admin.Handler(authService))
	router.Handle("GET /web/", http.StripPrefix("/web", http.FileServerFS(web.FS)))
	router.Handle("GET /favicon.ico", fsFileHandler(web.FS, "favicon_io/favicon.ico"))
	router.Handle("GET /site.webmanifest", fsFileHandler(web.FS, "favicon_io/site.webmanifest"))
	router.Handle("GET /favicon-16x16.png", fsFileHandler(web.FS, "favicon_io/favicon-16x16.png"))
	router.Handle("GET /favicon-32x32.png", fsFileHandler(web.FS, "favicon_io/favicon-32x32.png"))
	router.Handle("GET /apple-touch-icon.png", fsFileHandler(web.FS, "favicon_io/apple-touch-icon.png"))
	router.Handle("GET /android-chrome-192x192.png", fsFileHandler(web.FS, "favicon_io/android-chrome-192x192.png"))
	router.Handle("GET /android-chrome-512x512.png", fsFileHandler(web.FS, "favicon_io/android-chrome-512x512.png"))

	blobHandler := blob.NewAPIHandler(blobService)
	router.HandleFunc("GET /cdn/{cdnid}", errorHandler(blobHandler.GetBlob))

	authHandler := auth.NewAPIHandler(authService)
	router.HandleFunc("POST /login", errorHandler(authHandler.PostLogin))
	router.HandleFunc("POST /logout", errorHandler(authHandler.PostLogout))

	v1 := http.NewServeMux()

	v1.HandleFunc("GET /blob/records", errorHandler(blobHandler.GetRecords))
	v1.HandleFunc("POST /blob/upload", errorHandler(blobHandler.PostUpload))
	v1.HandleFunc("POST /blob/remove", errorHandler(blobHandler.PostRemove))

	randomNamesHandler := randomnames.NewAPIHandler(randomNamesService)
	v1.HandleFunc("GET /randomnames/form", errorHandler(randomNamesHandler.GetForm))
	v1.HandleFunc("POST /randomnames/form", errorHandler(randomNamesHandler.PostForm))
	v1.HandleFunc("GET /randomnames/records", errorHandler(randomNamesHandler.GetRecords))
	v1.HandleFunc("POST /randomnames/remove", errorHandler(randomNamesHandler.PostRemove))

	protected := middleware.Protected(authService)
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", protected(v1)))

	stack := middleware.CreateStack(
		middleware.Gzip(),
		middleware.Secure(),
		middleware.IPExtractor(),
		middleware.Logger(logger.Get()),
		middleware.Recover(),
		middleware.RateLimiter(50, 50, 3*time.Minute),
		csrf.Protect(
			[]byte(config.Get().Secure.CSRF.Key),
			csrf.Secure(config.Get().Secure.CSRF.Secure),
			csrf.TrustedOrigins(config.Get().Secure.CSRF.TrustedOrigins),
			csrf.CookieName("csrfcookie"),
			csrf.Path("/"),
		),
	)

	server := &http.Server{
		Handler:      stack(router),
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &Server{server: server}
}

func (s *Server) Start(address string) error {
	s.server.Addr = address
	logger.Get().Info("starting the api server on " + address)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	logger.Get().Info("shutting down the api server")
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Get().Error("[api]" + err.Error())
	}
}

func fsFileHandler(fsys fs.FS, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, fsys, name)
	}
}

func errorHandler(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		if errors.Is(err, context.Canceled) {
			w2.NewErrorResponse("Client Closed Request").Write(w, 499)
			return
		}

		if ctxe, ok := r.Context().Value("err").(*error); ok {
			*ctxe = err
		}

		status := wrap.HTTPStatus(err)
		message := err.Error()
		if status >= 500 {
			message = http.StatusText(status)
		}

		w2.NewErrorResponse(message).Write(w, status)
	}
}
