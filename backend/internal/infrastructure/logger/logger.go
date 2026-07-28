package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New() *Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Logger{
		logger,
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	args = enrichLogWithCtx(ctx, args...)
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	args = enrichLogWithCtx(ctx, args...)
	l.logger.ErrorContext(ctx, msg, args...)
}

func enrichLogWithCtx(ctx context.Context, args ...any) []any {
	_ = ctx
	return args
}
