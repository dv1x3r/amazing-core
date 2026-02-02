package logger

import (
	"log/slog"
	"os"

	"github.com/dv1x3r/amazing-core/internal/lib/prettyslog"
)

type HandlerType int

const (
	_ HandlerType = iota
	TextHandler
	JsonHandler
	PrettyHandler
)

var (
	level  slog.LevelVar
	logger *slog.Logger
)

func Create(handlerType HandlerType) {
	switch handlerType {
	case TextHandler:
		logger = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: &level},
		))
	case JsonHandler:
		logger = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: &level},
		))
	case PrettyHandler:
		logger = slog.New(prettyslog.NewHandler(
			&slog.HandlerOptions{Level: &level},
		))
	}
}

func SetLevel(l slog.Level) {
	level.Set(l)
}

func Get() *slog.Logger {
	return logger
}
