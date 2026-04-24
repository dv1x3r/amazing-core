package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dv1x3r/amazing-core/internal/network/gsf"
)

func Logger(logger *slog.Logger) gsf.Middleware {
	return func(next gsf.HandlerFunc) gsf.HandlerFunc {
		return func(w gsf.ResponseWriter, r *gsf.Request) error {
			debug := logger.Handler().Enabled(context.TODO(), slog.LevelDebug)

			startTime := time.Now()
			err := next(w, r)
			latency := fmt.Sprintf("%.2fms", float64(time.Since(startTime).Microseconds())/1000)

			logFn := logger.Info
			if err != nil {
				logFn = logger.Error
			} else if debug {
				logFn = logger.Debug
			}

			attrs := []any{
				slog.String("remote_ip", r.RemoteIP()),
				slog.String("platform", r.Platform().String()),
				slog.Int("result_code", int(w.Header().ResultCode)),
				slog.Int("app_code", int(w.Header().AppCode)),
				slog.Int("request_id", int(r.Header().RequestID)),
				slog.Int("req_flags", int(r.Header().Flags)),
				slog.Int("res_flags", int(w.Header().Flags)),
				slog.Int("svc_class", int(r.Header().SvcClass)),
				slog.Int("msg_type", int(r.Header().MsgType)),
				slog.String("latency", latency),
			}

			if debug {
				attrs = append(attrs, slog.String("result_code_text", w.Header().ResultCodeText()))
				attrs = append(attrs, slog.String("app_code_text", w.Header().AppCodeText()))
				attrs = append(attrs, slog.String("svc_class_text", r.Header().ServiceClassText()))
				attrs = append(attrs, slog.String("msg_type_text", r.Header().MessageTypeText()))
				attrs = append(attrs, slog.Any("request", r.Body()))
				attrs = append(attrs, slog.Any("response", w.Body()))
			}

			if err != nil {
				attrs = append(attrs, slog.String("error", err.Error()))
			}

			logFn("gsf", attrs...)
			return err
		}
	}
}
