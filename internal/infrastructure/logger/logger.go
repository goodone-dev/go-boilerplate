package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/utils/masker"
	"github.com/rs/zerolog"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"
)

var zLogger zerolog.Logger
var oLogger otellog.Logger

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	zLogger = zerolog.New(output).With().Timestamp().Logger()
	zLogger = zLogger.Hook(TracingHook{})
}

func Disabled() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	zLogger = zerolog.Nop()
}

type TracingHook struct{}

func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		e.Str("request_id", span.SpanContext().TraceID().String()) // TODO: Get from context
	}
}

type Metadata map[string]any

type LogBuilder struct {
	ctx      context.Context
	level    zerolog.Level
	msg      string
	err      error
	metadata Metadata
	fields   map[string]any
}

func Trace(ctx context.Context, msg string) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.TraceLevel,
		msg:    msg,
		fields: make(map[string]any),
	}
}

func Tracef(ctx context.Context, format string, args ...any) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.TraceLevel,
		msg:    fmt.Sprintf(format, args...),
		fields: make(map[string]any),
	}
}

func Debug(ctx context.Context, msg string) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.DebugLevel,
		msg:    msg,
		fields: make(map[string]any),
	}
}

func Debugf(ctx context.Context, format string, args ...any) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.DebugLevel,
		msg:    fmt.Sprintf(format, args...),
		fields: make(map[string]any),
	}
}

func Info(ctx context.Context, msg string) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.InfoLevel,
		msg:    msg,
		fields: make(map[string]any),
	}
}

func Infof(ctx context.Context, format string, args ...any) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.InfoLevel,
		msg:    fmt.Sprintf(format, args...),
		fields: make(map[string]any),
	}
}

func Warn(ctx context.Context, msg string) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.WarnLevel,
		msg:    msg,
		fields: make(map[string]any),
	}
}

func Warnf(ctx context.Context, format string, args ...any) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.WarnLevel,
		msg:    fmt.Sprintf(format, args...),
		fields: make(map[string]any),
	}
}

func Error(ctx context.Context, err error, msg string) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.ErrorLevel,
		msg:    msg,
		err:    err,
		fields: make(map[string]any),
	}
}

func Errorf(ctx context.Context, err error, format string, args ...any) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.ErrorLevel,
		msg:    fmt.Sprintf(format, args...),
		err:    err,
		fields: make(map[string]any),
	}
}

func Fatal(ctx context.Context, err error, msg string) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.FatalLevel,
		msg:    msg,
		err:    err,
		fields: make(map[string]any),
	}
}

func Fatalf(ctx context.Context, err error, format string, args ...any) *LogBuilder {
	return &LogBuilder{
		ctx:    ctx,
		level:  zerolog.FatalLevel,
		msg:    fmt.Sprintf(format, args...),
		err:    err,
		fields: make(map[string]any),
	}
}

func (b *LogBuilder) SetMetadata(metadata Metadata) *LogBuilder {
	b.metadata = metadata
	return b
}

func (b *LogBuilder) AddField(key string, val any) *LogBuilder {
	b.fields[key] = val
	return b
}

func (b *LogBuilder) Write() {
	if int(b.level) < config.Logger.Level {
		return
	}

	var zlog *zerolog.Event
	switch b.level {
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

	var metadata []byte
	if len(b.metadata) > 0 {
		masked := masker.Mask(b.metadata)
		metadata, _ = json.Marshal(masked)
	}

	if metadata != nil {
		zlog.RawJSON("metadata", metadata)
	}

	var fields map[string]any
	if len(b.fields) > 0 {
		fields = masker.Mask(b.fields).(map[string]any)
	}

	for k, v := range fields {
		zlog.Any(k, v)
	}

	file, line, funcName := getCaller(2)

	if b.err != nil {
		zlog.Err(b.err)
		zlog.Str("error_caller", fmt.Sprintf("%s:%d", file, line))
	}

	zlog.Ctx(b.ctx).Msg(b.msg)

	if oLogger == nil {
		return
	}

	record := otellog.Record{}
	record.SetSeverity(otellog.Severity(b.level))
	record.SetSeverityText(b.level.String())
	record.SetBody(otellog.StringValue(b.msg))
	record.SetTimestamp(time.Now())
	record.AddAttributes(
		otellog.String("code.filepath", file),
		otellog.Int("code.lineno", line),
		otellog.String("code.function", funcName),
	)

	if metadata != nil {
		record.AddAttributes(
			otellog.String("metadata", string(metadata)),
		)
	}

	for k, v := range fields {
		record.AddAttributes(
			otellog.String(k, fmt.Sprintf("%v", v)),
		)
	}

	if b.err != nil {
		record.AddAttributes(
			otellog.String("error.message", b.err.Error()),
			otellog.String("error.type", fmt.Sprintf("%T", b.err)),
		)
	}

	oLogger.Emit(b.ctx, record)
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
