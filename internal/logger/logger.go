package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

func Init(level slog.Level, production bool) {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: !production && level == slog.LevelDebug,
	}

	var handler slog.Handler
	if production {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	slog.SetDefault(slog.New(handler))
}

func InitWithWriter(w io.Writer, level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{Level: level}
	return slog.New(slog.NewTextHandler(w, opts))
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

type loggerKey struct{}

func With(args ...any) *slog.Logger {
	return slog.Default().With(args...)
}
