package specctx

import (
	"context"

	"golang.org/x/exp/slog"
)

type ctxKey int

var logCtxKey ctxKey

func WithLogger(parent context.Context, log *slog.Logger) context.Context {
	return context.WithValue(parent, logCtxKey, log)
}

func GetLogger(ctx context.Context) *slog.Logger {
	out := ctx.Value(logCtxKey)
	if out == nil {
		return nil
	}

	return out.(*slog.Logger)
}
