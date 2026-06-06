package slogbus

import (
	"log/slog"
	"time"
)

type Record struct {
	ID      uint64         `json:"id"`
	Time    time.Time      `json:"time"`
	Level   string         `json:"level"`
	Message string         `json:"message"`
	Attrs   map[string]any `json:"attrs"`
}

func attrsToMap(attrs []slog.Attr) map[string]any {
	data := make(map[string]any, len(attrs))
	for _, attr := range attrs {
		if attr.Key == "" {
			continue
		}
		data[attr.Key] = valueToAny(attr.Value)
	}
	return data
}

func valueToAny(value slog.Value) any {
	value = value.Resolve()
	if value.Kind() == slog.KindGroup {
		return attrsToMap(value.Group())
	}
	return value.Any()
}
