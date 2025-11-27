package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/rs/zerolog"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"
)

var zOutput = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
var zLogger = zerolog.New(zOutput).
	With().
	Timestamp().
	Logger()

func Disabled() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	zLogger = zerolog.Nop()
}

func Trace(ctx context.Context, msg string) {
	log := zLogger.Trace()

	recordLog(ctx, log, otellog.SeverityTrace, msg)
}

func Tracef(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Trace()

	recordLog(ctx, log, otellog.SeverityTrace, msg)
}

func Debug(ctx context.Context, msg string) {
	log := zLogger.Debug()

	recordLog(ctx, log, otellog.SeverityDebug, msg)
}

func Debugf(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Debug()

	recordLog(ctx, log, otellog.SeverityDebug, msg)
}

func Info(ctx context.Context, msg string) {
	log := zLogger.Info()

	recordLog(ctx, log, otellog.SeverityInfo, msg)
}

func Infof(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Info()

	recordLog(ctx, log, otellog.SeverityInfo, msg)
}

func Warn(ctx context.Context, msg string) {
	log := zLogger.Warn()

	recordLog(ctx, log, otellog.SeverityWarn, msg)
}

func Warnf(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Warn()

	recordLog(ctx, log, otellog.SeverityWarn, msg)
}

func Error(ctx context.Context, err error, msg string) {
	log := zLogger.Error().Err(err)

	recordLog(ctx, log, otellog.SeverityError, msg)
}

func Errorf(ctx context.Context, err error, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Error().Err(err)

	recordLog(ctx, log, otellog.SeverityError, msg)
}

func Fatal(ctx context.Context, err error, msg string) {
	log := zLogger.Fatal().Err(err)

	recordLog(ctx, log, otellog.SeverityFatal, msg)
}

func Fatalf(ctx context.Context, err error, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Fatal().Err(err)

	recordLog(ctx, log, otellog.SeverityFatal, msg)
}

func recordLog(ctx context.Context, log *zerolog.Event, severity otellog.Severity, msg string) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		log.Str("request_id", span.SpanContext().TraceID().String())
	}

	log.Msg(msg)

	if !config.LoggerConfig.Enabled {
		return
	}

	record := otellog.Record{}
	record.SetSeverity(severity)
	record.SetSeverityText(severity.String())
	record.SetBody(otellog.StringValue(msg))
	record.SetTimestamp(time.Now())

	oLogger.Emit(ctx, record)
}
