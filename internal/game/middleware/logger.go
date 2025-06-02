package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dv1x3r/amazing-core/internal/game/gsf"
	"github.com/dv1x3r/amazing-core/internal/game/types/clientmessagetype"
	"github.com/dv1x3r/amazing-core/internal/game/types/serviceclass"
	"github.com/dv1x3r/amazing-core/internal/game/types/syncmessagetype"
	"github.com/dv1x3r/amazing-core/internal/game/types/usermessagetype"
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

			svcClass := serviceclass.ServiceClass(w.Header().SvcClass).String()
			msgType := fmt.Sprint(w.Header().MsgType)

			switch w.Header().SvcClass {
			case int32(serviceclass.USER_SERVER):
				msgType = usermessagetype.UserMessageType(w.Header().MsgType).String()
			case int32(serviceclass.SYNC_SERVER):
				msgType = syncmessagetype.SyncMessageType(w.Header().MsgType).String()
			case int32(serviceclass.CLIENT):
				msgType = clientmessagetype.ClientMessageType(w.Header().MsgType).String()
			}

			message := fmt.Sprintf("[gsf] %s %s.%s %+v", r.RemoteAddr, svcClass, msgType, w.Header())

			logFn(message, attrs...)
			return err
		}
	}
}
