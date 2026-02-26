package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type loggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
	err        error
}

func (w *loggerResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggerResponseWriter) SetError(err error) {
	w.err = err
}

func Logger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := &loggerResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			startTime := time.Now()
			next.ServeHTTP(lw, r)
			latency := fmt.Sprintf("%.2fms", float64(time.Since(startTime).Microseconds())/1000)

			logFn := logger.Info
			if lw.statusCode >= 500 {
				logFn = logger.Error
			} else if lw.statusCode >= 300 {
				logFn = logger.Warn
			}

			remoteIP := r.Context().Value(IPExtractorKey).(string)
			statusText := http.StatusText(lw.statusCode)

			attrs := []any{
				slog.String("remote_ip", remoteIP),
				slog.Int("status", lw.statusCode),
				slog.String("status_text", statusText),
				slog.String("host", r.Host),
				slog.String("method", r.Method),
				slog.String("uri", r.URL.String()),
				slog.String("latency", latency),
			}

			if lw.err != nil {
				attrs = append(attrs, slog.String("error", lw.err.Error()))
			}

			logFn("http", attrs...)
		})
	}
}
