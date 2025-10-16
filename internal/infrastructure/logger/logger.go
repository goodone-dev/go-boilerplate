package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

var slogger = otelslog.NewLogger("demo")
var logger = zerolog.New(os.Stderr).
	With().
	Timestamp().
	Logger().
	Output(zerolog.ConsoleWriter{Out: os.Stderr})

func Trace(ctx context.Context, msg string, args ...any) {
	slogger.DebugContext(ctx, msg, args...)
	printLog(logger.Trace(), msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	slogger.DebugContext(ctx, msg, args...)
	printLog(logger.Debug(), msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	slogger.InfoContext(ctx, msg, args...)
	printLog(logger.Info(), msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	slogger.WarnContext(ctx, msg, args...)
	printLog(logger.Warn(), msg, args...)
}

func Error(ctx context.Context, msg string, err error, args ...any) {
	slogger.ErrorContext(ctx, msg, args...)
	printLog(logger.Error().Err(err), msg, args...)
}

func Fatal(ctx context.Context, msg string, err error, args ...any) {
	slogger.ErrorContext(ctx, msg, args...)
	printLog(logger.Fatal().Err(err), msg, args...)
}

func Panic(ctx context.Context, msg string, err error, args ...any) {
	slogger.ErrorContext(ctx, msg, args...)
	printLog(logger.Panic().Err(err), msg, args...)
}

func printLog(log *zerolog.Event, msg string, args ...any) {
	if len(args) > 0 {
		arr := zerolog.Arr()
		for _, arg := range args {
			arr.Interface(arg)
		}

		log.Array("args", arr)
	}

	log.Msg(msg)
}
