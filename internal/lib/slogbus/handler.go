package slogbus

import (
	"context"
	"log/slog"
)

type handler struct {
	slog.Handler
	bus *Bus
}

func (b *Bus) Handler(next slog.Handler) slog.Handler {
	return &handler{
		Handler: next,
		bus:     b,
	}
}

func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	attrs := make([]slog.Attr, 0, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})

	h.bus.publish(Record{
		Time:    record.Time,
		Level:   record.Level.String(),
		Message: record.Message,
		Attrs:   attrsToMap(attrs),
	})

	return h.Handler.Handle(ctx, record)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		Handler: h.Handler.WithAttrs(attrs),
		bus:     h.bus,
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		Handler: h.Handler.WithGroup(name),
		bus:     h.bus,
	}
}
