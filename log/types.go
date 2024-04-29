package log

import (
	"context"
	"log/slog"
	"time"
)

const RequestIDKey = "requestID"

type LogFunc func(l slog.Logger)

type Context struct {
	Logger *slog.Logger
	Ctx    context.Context
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return c.Ctx.Deadline()
}

func (c Context) Done() <-chan struct{} {
	return c.Ctx.Done()
}

func (c Context) Err() error {
	return c.Ctx.Err()
}

func (c Context) Value(key any) any {
	if key == "Logger" {
		return c.Logger
	}

	return c.Ctx.Value(key)
}

// NewCtxWithRequestID returns new empty context with logger and requestID.
// It doesn't affect cancel or deadline from parent context
func NewCtxWithRequestID(ctx context.Context) Context {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		requestID = ""
	}

	return NewContext(requestID)
}

func NewContext(uuid string) Context {
	return Context{
		slog.With(RequestIDKey, uuid),
		context.WithValue(context.Background(), RequestIDKey, uuid),
	}
}

func Log(ctx context.Context, logFunc LogFunc) {
	loggerCtx, ok := ctx.(Context)

	if ok {
		logFunc(*loggerCtx.Logger)
	} else {
		logFunc(*slog.Default())
	}
}
