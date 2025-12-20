package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/utils/masker"
	"github.com/rs/zerolog"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"
)

var zLogger zerolog.Logger
var oLogger otellog.Logger

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	zLogger = zerolog.New(output).
		With().
		Timestamp().
		Logger()
}

func Disabled() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	zLogger = zerolog.Nop()
}

type LogBuilder struct {
	metadata map[string]any
}

func With() *LogBuilder {
	return &LogBuilder{
		metadata: make(map[string]any),
	}
}

func (b *LogBuilder) Metadata(key string, val any) *LogBuilder {
	b.metadata[key] = val
	return b
}

func (b *LogBuilder) Trace(ctx context.Context, msg string) {
	b.writeLog(ctx, zerolog.TraceLevel, msg, nil)
}

func (b *LogBuilder) Tracef(ctx context.Context, format string, args ...any) {
	b.writeLog(ctx, zerolog.TraceLevel, fmt.Sprintf(format, args...), nil)
}

func (b *LogBuilder) Debug(ctx context.Context, msg string) {
	b.writeLog(ctx, zerolog.DebugLevel, msg, nil)
}

func (b *LogBuilder) Debugf(ctx context.Context, format string, args ...any) {
	b.writeLog(ctx, zerolog.DebugLevel, fmt.Sprintf(format, args...), nil)
}

func (b *LogBuilder) Info(ctx context.Context, msg string) {
	b.writeLog(ctx, zerolog.InfoLevel, msg, nil)
}

func (b *LogBuilder) Infof(ctx context.Context, format string, args ...any) {
	b.writeLog(ctx, zerolog.InfoLevel, fmt.Sprintf(format, args...), nil)
}

func (b *LogBuilder) Warn(ctx context.Context, msg string) {
	b.writeLog(ctx, zerolog.WarnLevel, msg, nil)
}

func (b *LogBuilder) Warnf(ctx context.Context, format string, args ...any) {
	b.writeLog(ctx, zerolog.WarnLevel, fmt.Sprintf(format, args...), nil)
}

func (b *LogBuilder) Error(ctx context.Context, err error, msg string) {
	b.writeLog(ctx, zerolog.ErrorLevel, msg, err)
}

func (b *LogBuilder) Errorf(ctx context.Context, err error, format string, args ...any) {
	b.writeLog(ctx, zerolog.ErrorLevel, fmt.Sprintf(format, args...), err)
}

func (b *LogBuilder) Fatal(ctx context.Context, err error, msg string) {
	b.writeLog(ctx, zerolog.FatalLevel, msg, err)
}

func (b *LogBuilder) Fatalf(ctx context.Context, err error, format string, args ...any) {
	b.writeLog(ctx, zerolog.FatalLevel, fmt.Sprintf(format, args...), err)
}

func (b *LogBuilder) writeLog(ctx context.Context, level zerolog.Level, msg string, err error) {
	file, line, fn := getCaller(3)

	var metadata []byte
	if len(b.metadata) > 0 {
		masked := masker.Mask(b.metadata)
		metadata, _ = json.Marshal(masked)
	}

	var zlog *zerolog.Event
	switch level {
	case zerolog.TraceLevel:
		zlog = zLogger.Trace()
	case zerolog.DebugLevel:
		zlog = zLogger.Debug()
	case zerolog.InfoLevel:
		zlog = zLogger.Info()
	case zerolog.WarnLevel:
		zlog = zLogger.Warn()
	case zerolog.ErrorLevel:
		zlog = zLogger.Error()
	case zerolog.FatalLevel:
		zlog = zLogger.Fatal()
	default:
		zlog = zLogger.Info()
	}

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		zlog.Str("request_id", span.SpanContext().TraceID().String())
	}

	if metadata != nil {
		zlog.Interface("metadata", string(metadata))
	}

	if err != nil {
		zlog.Err(err)
		zlog.Str("filepath", fmt.Sprintf("%s:%d", file, line))
	}

	zlog.Msg(msg)

	if oLogger == nil {
		return
	}

	record := otellog.Record{}
	record.SetSeverity(otellog.Severity(level))
	record.SetSeverityText(level.String())
	record.SetBody(otellog.StringValue(msg))
	record.SetTimestamp(time.Now())
	record.AddAttributes(
		otellog.String("code_filepath", file),
		otellog.Int("code_lineno", line),
		otellog.String("code_function", fn),
	)

	if metadata != nil {
		record.AddAttributes(
			otellog.String("metadata", string(metadata)),
		)
	}

	if err != nil {
		record.AddAttributes(
			otellog.String("error_message", err.Error()),
			otellog.String("error_type", fmt.Sprintf("%T", err)),
		)
	}

	oLogger.Emit(ctx, record)
}

func getCaller(skip int) (file string, line int, function string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0, "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn != nil {
		function = fn.Name()
	}

	return file, line, function
}
