package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/utils/masker"
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

type Log struct {
	metadatas []any
}

func WithMetadata(metadatas ...any) *Log {
	return &Log{metadatas: metadatas}
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

func (l *Log) Trace(ctx context.Context, msg string) {
	log := zLogger.Trace()

	recordLog(ctx, log, otellog.SeverityTrace, msg, l.metadatas...)
}

func (l *Log) Tracef(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Trace()

	recordLog(ctx, log, otellog.SeverityTrace, msg, l.metadatas...)
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

func (l *Log) Debug(ctx context.Context, msg string) {
	log := zLogger.Debug()

	recordLog(ctx, log, otellog.SeverityDebug, msg, l.metadatas...)
}

func (l *Log) Debugf(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Debug()

	recordLog(ctx, log, otellog.SeverityDebug, msg, l.metadatas...)
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

func (l *Log) Info(ctx context.Context, msg string) {
	log := zLogger.Info()

	recordLog(ctx, log, otellog.SeverityInfo, msg, l.metadatas...)
}

func (l *Log) Infof(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Info()

	recordLog(ctx, log, otellog.SeverityInfo, msg, l.metadatas...)
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

func (l *Log) Warn(ctx context.Context, msg string) {
	log := zLogger.Warn()

	recordLog(ctx, log, otellog.SeverityWarn, msg, l.metadatas...)
}

func (l *Log) Warnf(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Warn()

	recordLog(ctx, log, otellog.SeverityWarn, msg, l.metadatas...)
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

func (l *Log) Error(ctx context.Context, err error, msg string) {
	log := zLogger.Error().Err(err)

	recordLog(ctx, log, otellog.SeverityError, msg, l.metadatas...)
}

func (l *Log) Errorf(ctx context.Context, err error, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Error().Err(err)

	recordLog(ctx, log, otellog.SeverityError, msg, l.metadatas...)
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

func (l *Log) Fatal(ctx context.Context, err error, msg string) {
	log := zLogger.Fatal().Err(err)

	recordLog(ctx, log, otellog.SeverityFatal, msg, l.metadatas...)
}

func (l *Log) Fatalf(ctx context.Context, err error, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log := zLogger.Fatal().Err(err)

	recordLog(ctx, log, otellog.SeverityFatal, msg, l.metadatas...)
}

func recordLog(ctx context.Context, log *zerolog.Event, severity otellog.Severity, msg string, metadatas ...any) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		log.Str("request_id", span.SpanContext().TraceID().String())
	}

	metadata := masker.Mask(metadatas)
	if metadata != nil && len(metadatas) > 0 {
		log.Interface("metadata", metadata)
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

	if metadata != nil && len(metadatas) > 0 {
		if jsonMetadata, err := json.Marshal(metadata); err == nil {
			record.AddAttributes(otellog.String("metadata", string(jsonMetadata)))
		}
	}

	oLogger.Emit(ctx, record)
}
