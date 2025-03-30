package logger

import "log/slog"

var logger *slog.Logger

func Set(l *slog.Logger) {
	logger = l
}

func Get() *slog.Logger {
	return logger
}
