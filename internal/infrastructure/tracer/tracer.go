package tracer

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/utils/masker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type tracer struct {
	span       trace.Span
	attributes map[string]any
	funcInput  map[string]any
	funcOutput map[string]any
}

type Metadata map[string]any

func Start(ctx context.Context) (context.Context, *tracer) {
	if !config.Tracer.Enabled {
		return ctx, nil
	}

	file, line, funcName := getCaller(2)
	spanName := parseSpanName(funcName)

	ctx, span := otel.Tracer(config.Application.Name).Start(ctx, spanName)
	span.SetAttributes(
		attribute.String("code.filepath", file),
		attribute.Int("code.lineno", line),
		attribute.String("code.function", funcName),
	)

	return ctx, &tracer{
		span:       span,
		attributes: make(map[string]any),
	}
}

func (t *tracer) AddAttribute(key string, value any) *tracer {
	t.attributes[key] = value
	return t
}

func (t *tracer) SetFunctionInput(metadata Metadata) *tracer {
	t.funcInput = metadata
	return t
}

func (t *tracer) SetFunctionOutput(metadata Metadata) *tracer {
	t.funcOutput = metadata
	return t
}

func (t *tracer) End(err error) {
	if !config.Tracer.Enabled {
		return
	}

	var funcInput []byte
	if len(t.funcInput) > 0 {
		masked := masker.Mask(t.funcInput)
		funcInput, _ = json.Marshal(masked)
	}

	if funcInput != nil {
		t.span.SetAttributes(attribute.String("function.input", string(funcInput)))
	}

	var funcOutput []byte
	if len(t.funcOutput) > 0 {
		masked := masker.Mask(t.funcOutput)
		funcOutput, _ = json.Marshal(masked)
	}

	if funcOutput != nil {
		t.span.SetAttributes(attribute.String("function.output", string(funcOutput)))
	}

	var attributes map[string]any
	if len(t.attributes) > 0 {
		attributes = masker.Mask(t.attributes).(map[string]any)
	}

	for k, v := range attributes {
		t.span.SetAttributes(attribute.String(k, fmt.Sprintf("%v", v)))
	}

	if err != nil {
		t.span.RecordError(err)
		t.span.SetStatus(codes.Error, err.Error())
		t.span.SetAttributes(
			attribute.String("error.message", err.Error()),
			attribute.String("error.type", fmt.Sprintf("%T", err)),
		)
	}

	t.span.End()
}

func getCaller(skip int) (file string, line int, funcName string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0, "unknown"
	}

	return file, line, runtime.FuncForPC(pc).Name()
}

func parseSpanName(funcName string) string {
	if funcName == "" || funcName == "unknown" {
		return "unknown"
	}

	// Find the last "/" to get the package and function part
	// e.g., "github.com/org/repo/pkg/subpkg.(*Type).Method" → "subpkg.(*Type).Method"
	lastSlashIdx := strings.LastIndex(funcName, "/")
	shortName := funcName
	if lastSlashIdx != -1 {
		shortName = funcName[lastSlashIdx+1:]
	}

	// Split by "." to separate package from type/function
	// e.g., "subpkg.(*Type).Method" → ["subpkg", "(*Type).Method"]
	parts := strings.SplitN(shortName, ".", 2)
	if len(parts) < 2 {
		return shortName
	}

	result := parts[1]

	// Clean up pointer receiver: "(*redisClient).Get" → "redisClient.Get"
	result = strings.ReplaceAll(result, "(*", "")
	result = strings.ReplaceAll(result, ")", "")

	// Clean up generics: "baseRepo[...].FindById" → "baseRepo.FindById"
	if bracketStart := strings.Index(result, "["); bracketStart != -1 {
		if bracketEnd := strings.Index(result, "]"); bracketEnd != -1 {
			result = result[:bracketStart] + result[bracketEnd+1:]
		}
	}

	// Handle chained methods: "NewRouter.RateLimiterHandler.func10" → "RateLimiterHandler.func10"
	// Keep only the last two parts (Type.Method)
	dotParts := strings.Split(result, ".")
	if len(dotParts) > 2 {
		result = dotParts[len(dotParts)-2] + "." + dotParts[len(dotParts)-1]
	}

	return result
}
