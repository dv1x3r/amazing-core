package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/dv1x3r/amazing-core/internal/lib/prettyslog"
)

type loggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggerResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			lw := &loggerResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			var err error
			r = r.WithContext(context.WithValue(r.Context(), "err", &err))
			next.ServeHTTP(lw, r)

			logFn := logger.Info
			color := prettyslog.Green

			if lw.statusCode >= 500 {
				logFn = logger.Error
				color = prettyslog.Red
			} else if lw.statusCode >= 300 {
				logFn = logger.Warn
				color = prettyslog.Yellow
			}

			remoteIP := r.Context().Value(IPExtractorKey).(string)
			took := fmt.Sprintf("%.2fms", float64(time.Since(startTime).Microseconds())/1000)

			statusCode := prettyslog.Colorize(color, strconv.Itoa(lw.statusCode))
			statusText := http.StatusText(lw.statusCode)
			if lw.statusCode == 499 {
				statusText = "Client Closed Request"
			}

			message := fmt.Sprintf("[api] %s %s %s %s %s %s", remoteIP, r.Method, statusCode, statusText, r.URL.Path, took)

			var attrs []any
			if err != nil {
				attrs = append(attrs, slog.String("error", err.Error()))
			}

			if len(r.URL.Query()) > 0 {
				attrs = append(attrs, slog.Any("query", r.URL.Query()))
			}

			logFn(message, attrs...)
		})
	}
}
