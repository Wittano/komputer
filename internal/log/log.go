package log

import (
	"context"
	"github.com/rs/zerolog"
	"os"
)

// TODO Add saving logs into file
var logger = zerolog.New(os.Stdout).With().Ctx(context.Background()).Timestamp().Logger()

func Info(ctx context.Context, msg string) {
	logger.Info().Ctx(ctx).Str("traceID", getTraceID(ctx)).Msg(msg)
}

func Warn(ctx context.Context, msg string) {
	logger.Warn().Ctx(ctx).Str("traceID", getTraceID(ctx)).Msg(msg)
}

func Error(ctx context.Context, msg string, err error) {
	logger.Err(err).Ctx(ctx).Str("traceID", getTraceID(ctx)).Msg(msg)
}

func Fatal(ctx context.Context, msg string, err error) {
	logger.Fatal().Err(err).Ctx(ctx).Str("traceID", getTraceID(ctx)).Msg(msg)
}

func getTraceID(ctx context.Context) string {
	v := ctx.Value("traceID")

	if v != nil {
		return v.(string)
	}

	return ""
}
