package tracer

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type spanPrefixName string
type spanCustomName string

type customTracerSpan struct {
	trace.Span
}

func PrefixName(spanName string) *spanPrefixName {
	name := spanPrefixName(spanName)
	return &name
}

func (s *spanPrefixName) Start(ctx context.Context, params ...any) (context.Context, *customTracerSpan) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	funcParts := strings.Split(funcName, ".")
	methodName := funcParts[len(funcParts)-1]

	return startSpan(ctx, fmt.Sprintf("%s.%s", string(*s), methodName), funcName, params...)
}

func CustomName(spanName string) *spanCustomName {
	name := spanCustomName(spanName)
	return &name
}

func (s *spanCustomName) Start(ctx context.Context, params ...any) (context.Context, *customTracerSpan) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	return startSpan(ctx, string(*s), funcName, params...)
}

func Start(ctx context.Context, params ...any) (context.Context, *customTracerSpan) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	spanName := parseSpanName(funcName)

	return startSpan(ctx, spanName, funcName, params...)
}

func startSpan(ctx context.Context, spanName string, funcName string, params ...any) (context.Context, *customTracerSpan) {
	if !config.TracerConfig.Enabled {
		return ctx, nil
	}

	ctx, span := otel.Tracer("").Start(ctx, spanName)

	span.SetAttributes(
		attribute.String("code.namespace", funcName),
	)

	for i, param := range params {
		paramJSON, _ := json.Marshal(param)
		span.SetAttributes(
			attribute.String(fmt.Sprintf("func.param.%d", i), string(paramJSON)),
		)
	}

	return ctx, &customTracerSpan{span}
}

func (s *customTracerSpan) Stop(err error, returns ...any) {
	if !config.TracerConfig.Enabled {
		return
	}

	if err != nil {
		s.RecordError(err)
		s.SetStatus(codes.Error, err.Error())
	}

	for i, returnValue := range returns {
		returnJSON, _ := json.Marshal(returnValue)
		s.SetAttributes(
			attribute.String(fmt.Sprintf("func.return.%d", i), string(returnJSON)),
		)
	}

	s.End()
}

func parseSpanName(funcName string) string {
	// Catch usecase, repository and other method name
	re := regexp.MustCompile(`\(\*?([^)]+)\)\.([^.]+)$`)
	matches := re.FindStringSubmatch(funcName)

	if len(matches) == 3 {
		typeName := matches[1]
		methodName := matches[2]

		re = regexp.MustCompile(`([^.]+)$`)
		typeName = re.FindString(typeName)
		if len(typeName) > 0 {
			typeName = strings.ToUpper(string(typeName[0])) + typeName[1:]
		}

		return typeName + "." + methodName
	}

	// Catch middleware method name
	re = regexp.MustCompile(`\.([^.]+)\.func\d+$`)
	matches = re.FindStringSubmatch(funcName)
	if len(matches) >= 2 {
		return matches[1]
	}

	return funcName
}
