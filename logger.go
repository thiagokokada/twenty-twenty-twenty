package main

import (
	"context"
	"log/slog"
)

type logHandler struct {
	slog.Handler
	*slog.HandlerOptions
}

func (l logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= l.HandlerOptions.Level.Level()
}
func (l logHandler) Handle(ctx context.Context, r slog.Record) error {
	return l.Handler.Handle(ctx, r)
}

func (l logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return logHandler{Handler: l.Handler.WithAttrs(attrs)}
}

func (l logHandler) WithGroup(name string) slog.Handler {
	return logHandler{Handler: l.Handler.WithGroup(name)}
}
