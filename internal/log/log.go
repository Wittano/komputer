package log

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func Info(ctx context.Context, msg string) {
	logRequest(ctx, log.Info(), msg)
}

func Warn(ctx context.Context, msg string) {
	logRequest(ctx, log.Warn(), msg)
}

func Error(ctx context.Context, msg string, err error) {
	log.Err(err).Ctx(ctx).Str("traceID", ctx.Value("traceID").(string)).Msg(msg)
}

func Fatal(ctx context.Context, msg string, err error) {
	log.Fatal().Err(err).Ctx(ctx).Str("traceID", ctx.Value("traceID").(string)).Msg(msg)
}

func logRequest(ctx context.Context, e *zerolog.Event, msg string) {
	e.Ctx(ctx).Str("traceID", ctx.Value("traceID").(string)).Msg(msg)
}
