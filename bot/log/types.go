package log

import (
	"context"
	"log/slog"
	"time"
)

const RequestIDKey = "requestID"

type Func func(l slog.Logger)

type Context struct {
	Logger *slog.Logger
	ctx    context.Context
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c Context) Err() error {
	return c.ctx.Err()
}

func (c Context) Value(key any) any {
	if key == "Logger" {
		return c.Logger
	}

	return c.ctx.Value(key)
}

// NewCtxWithRequestID returns new empty context with logger and requestID.
// It doesn't affect cancel or deadline from parent context
func NewCtxWithRequestID(ctx context.Context) Context {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		requestID = ""
	}

	return NewContext(nil, requestID)
}

func NewContext(ctx context.Context, uuid string) Context {
	return Context{
		slog.With(RequestIDKey, uuid),
		context.WithValue(ctx, RequestIDKey, uuid),
	}
}
