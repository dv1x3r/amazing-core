package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
)

type gsfLogger struct {
	logger *slog.Logger
}

func newGSFLogger(logger *slog.Logger) *gsfLogger {
	return &gsfLogger{logger: logger}
}

func (l *gsfLogger) OnRequest(event gsf.RequestEvent) {
	debug := l.debugEnabled()

	logFn := l.logFunc(event.Err, debug)
	logMessage := "gsf service"
	kind := "service"
	direction := ""
	if event.RequestHeader.IsNotify() {
		logMessage = "gsf notify"
		kind = "sync_notify"
		direction = "inbound"
	}

	attrs := []any{
		slog.String("kind", kind),
		slog.String("remote_ip", event.Session.RemoteIP()),
		slog.String("platform", event.Session.Platform().String()),
		slog.Int("result_code", int(event.ResponseHeader.ResultCode)),
		slog.String("result_code_text", event.ResponseHeader.ResultCodeText()),
		slog.Int("app_code", int(event.ResponseHeader.AppCode)),
		slog.String("app_code_text", event.ResponseHeader.AppCodeText()),
		slog.Int("request_id", int(event.RequestHeader.RequestID)),
		slog.Int("flags", int(event.RequestHeader.Flags)),
		slog.Int("response_flags", int(event.ResponseHeader.Flags)),
		slog.Int("svc_class", int(event.RequestHeader.SvcClass)),
		slog.String("svc_class_text", event.RequestHeader.ServiceClassText()),
		slog.Int("msg_type", int(event.RequestHeader.MsgType)),
		slog.String("msg_type_text", event.RequestHeader.MessageTypeText()),
		slog.String("latency", formatLatency(event.Latency)),
	}
	if direction != "" {
		attrs = append(attrs, slog.String("direction", direction))
	}

	attrs = appendPlayerOID(attrs, event.Session)

	if debug {
		attrs = append(attrs, slog.Any("request", event.RequestBody))
		attrs = append(attrs, slog.Any("response", event.ResponseBody))
	}

	attrs = appendErrorAttrs(attrs, event.Err)

	logFn(logMessage, attrs...)
}

func (l *gsfLogger) OnNotify(event gsf.NotifyEvent) {
	debug := l.debugEnabled()

	logFn := l.logFunc(event.Err, debug)

	attrs := []any{
		slog.String("kind", "client_notify"),
		slog.String("direction", "outbound"),
		slog.String("remote_ip", event.Session.RemoteIP()),
		slog.String("platform", event.Session.Platform().String()),
		slog.Int("flags", int(event.Header.Flags)),
		slog.Int("svc_class", int(event.Header.SvcClass)),
		slog.String("svc_class_text", event.Header.ServiceClassText()),
		slog.Int("msg_type", int(event.Header.MsgType)),
		slog.String("msg_type_text", event.Header.MessageTypeText()),
		slog.String("latency", formatLatency(event.Latency)),
	}

	attrs = appendPlayerOID(attrs, event.Session)

	if debug {
		attrs = append(attrs, slog.Any("notify", event.Body))
	}

	attrs = appendErrorAttrs(attrs, event.Err)

	logFn("gsf notify", attrs...)
}

func (l *gsfLogger) OnUnhandled(session *gsf.Session, header *gsf.Header, data []byte) {
	kind := "unhandled_service"
	if header.IsNotify() {
		kind = "unhandled_notify"
	}
	attrs := []any{
		slog.String("kind", kind),
		slog.String("direction", "inbound"),
		slog.String("remote_ip", session.RemoteIP()),
		slog.String("platform", session.Platform().String()),
		slog.Int("svc_class", int(header.SvcClass)),
		slog.String("svc_class_text", header.ServiceClassText()),
		slog.Int("msg_type", int(header.MsgType)),
		slog.String("msg_type_text", header.MessageTypeText()),
		slog.Int("request_id", int(header.RequestID)),
		slog.Int("flags", int(header.Flags)),
		slog.String("hex", fmt.Sprintf("%x", data)),
	}

	attrs = appendPlayerOID(attrs, session)

	l.logger.Info("gsf unhandled", attrs...)
}

func (l *gsfLogger) debugEnabled() bool {
	return l.logger.Handler().Enabled(context.TODO(), slog.LevelDebug)
}

func (l *gsfLogger) logFunc(err error, debug bool) func(string, ...any) {
	if err != nil {
		return l.logger.Error
	}
	if debug {
		return l.logger.Debug
	}
	return l.logger.Info
}

func appendPlayerOID(attrs []any, session *gsf.Session) []any {
	if playerOID, ok := session.PlayerOID(); ok {
		attrs = append(attrs, slog.Int64("player_oid", playerOID))
	}
	return attrs
}

func formatLatency(latency time.Duration) string {
	return fmt.Sprintf("%.2fms", float64(latency.Microseconds())/1000)
}

func appendErrorAttrs(attrs []any, err error) []any {
	var panicError wrap.PanicError
	if errors.As(err, &panicError) {
		attrs = append(attrs, slog.String("error", panicError.Error()))
		attrs = append(attrs, slog.String("stack", panicError.Stack()))
	} else if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	return attrs
}
