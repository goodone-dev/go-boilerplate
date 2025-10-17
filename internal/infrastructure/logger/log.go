package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	otellog "go.opentelemetry.io/otel/log"
)

var zOutput = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
var zLogger = zerolog.New(zOutput).
	With().
	Timestamp().
	Logger()

func Trace(ctx context.Context, msg string) {
	zLogger.Trace().Msg(msg)

	recordLog(ctx, otellog.SeverityTrace, msg)
}

func Tracef(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	zLogger.Trace().Msg(msg)

	recordLog(ctx, otellog.SeverityTrace, msg)
}

func Debug(ctx context.Context, msg string) {
	zLogger.Debug().Msg(msg)

	recordLog(ctx, otellog.SeverityDebug, msg)
}

func Debugf(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	zLogger.Debug().Msg(msg)

	recordLog(ctx, otellog.SeverityDebug, msg)
}

func Info(ctx context.Context, msg string) {
	zLogger.Info().Msg(msg)

	recordLog(ctx, otellog.SeverityInfo, msg)
}

func Infof(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	zLogger.Info().Msg(msg)

	recordLog(ctx, otellog.SeverityInfo, msg)
}

func Warn(ctx context.Context, msg string) {
	zLogger.Warn().Msg(msg)

	recordLog(ctx, otellog.SeverityWarn, msg)
}

func Warnf(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	zLogger.Warn().Msg(msg)

	recordLog(ctx, otellog.SeverityWarn, msg)
}

func Error(ctx context.Context, err error, msg string) {
	zLogger.Error().Err(err).Msg(msg)

	recordLog(ctx, otellog.SeverityError, msg)
}

func Errorf(ctx context.Context, err error, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	zLogger.Error().Err(err).Msg(msg)

	recordLog(ctx, otellog.SeverityError, msg)
}

func Fatal(ctx context.Context, err error, msg string) {
	zLogger.Fatal().Err(err).Msg(msg)

	recordLog(ctx, otellog.SeverityFatal, msg)
}

func Fatalf(ctx context.Context, err error, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	zLogger.Fatal().Err(err).Msg(msg)

	recordLog(ctx, otellog.SeverityFatal, msg)
}

func recordLog(ctx context.Context, severity otellog.Severity, msg string) {
	record := otellog.Record{}
	record.SetSeverity(severity)
	record.SetSeverityText(severity.String())
	record.SetBody(otellog.StringValue(msg))
	record.SetTimestamp(time.Now())

	oLogger.Emit(ctx, record)
}
