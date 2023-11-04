package log

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
)

var (
	FileLog io.WriteCloser = os.Stdout
)

var logger zerolog.Logger

func init() {
	logFilePath := os.Getenv("LOG_FILE")
	if logFilePath != "" {
		var err error

		FileLog, err = os.Open(logFilePath)
		if err != nil {
			Fatal(context.Background(), fmt.Sprintf("Failed to open file '%s'", logFilePath), err)
			return
		}
	}

	logger = zerolog.New(FileLog).With().Ctx(context.Background()).Timestamp().Logger()
}

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
