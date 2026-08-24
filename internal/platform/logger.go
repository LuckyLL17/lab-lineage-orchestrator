package platform

import (
	"log/slog"
	"os"
)

type Logger struct {
	base *slog.Logger
}

func NewLogger() *Logger {
	return &Logger{base: slog.New(slog.NewTextHandler(os.Stdout, nil))}
}

func (l *Logger) Info(message string, args ...any) {
	l.base.Info(message, args...)
}

func (l *Logger) Error(message string, args ...any) {
	l.base.Error(message, args...)
}

func (l *Logger) Warn(message string, args ...any) {
	l.base.Warn(message, args...)
}
