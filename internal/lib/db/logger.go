package db

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func (l Logger) Printf(format string, v ...any) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

func (l Logger) Fatalf(format string, v ...any) {
	l.logger.Error(fmt.Sprintf(format, v...))
	fmt.Scanln()
	os.Exit(1)
}
