package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dv1x3r/amazing-core/internal/lib/webapi/sse"
)

func (h *Handler) GetLogsWatch(w http.ResponseWriter, r *http.Request) error {
	sse.SetHeaders(w)

	controller := http.NewResponseController(w)
	if err := controller.SetWriteDeadline(time.Time{}); err != nil {
		return err
	}

	event := sse.Event{Comment: []byte("watching")}
	if err := event.Flush(w); err != nil {
		return err
	}

	records := h.logBus.Watch(r.Context())
	keepalive := time.NewTicker(30 * time.Second)
	defer keepalive.Stop()

	for {
		select {
		case <-r.Context().Done():
			return nil
		case <-keepalive.C:
			event = sse.Event{Comment: []byte("keepalive")}
			if err := event.Flush(w); err != nil {
				return err
			}
		case record := <-records:
			data, err := json.Marshal(record)
			if err != nil {
				return err
			}
			id := strconv.FormatUint(record.ID, 10)
			event := sse.Event{ID: []byte(id), Event: []byte("log"), Data: data}
			if err := event.Flush(w); err != nil {
				return err
			}
		}
	}
}
