package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
)

func Logger(logger *slog.Logger) gsf.Middleware {
	return func(next gsf.HandlerFunc) gsf.HandlerFunc {
		return func(w gsf.ResponseWriter, r *gsf.Request) error {
			startTime := time.Now()
			err := next(w, r)
			took := fmt.Sprintf("%.2fms", float64(time.Since(startTime).Microseconds())/1000)

			logFn := logger.Info
			attrs := []any{
				slog.Any("request", r.Body()),
				slog.Any("response", w.Body()),
				slog.String("took", took),
			}

			if err != nil {
				logFn = logger.Error
				attrs = append(attrs, slog.String("error", err.Error()))
			}

			message := fmt.Sprintf("[gsf] %s %+v", r.RemoteAddr, w.Header())

			logFn(message, attrs...)
			return err
		}
	}
}
