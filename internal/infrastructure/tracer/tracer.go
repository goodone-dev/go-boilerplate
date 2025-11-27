package tracer

import (
	"context"
	"fmt"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func NewProvider(ctx context.Context) *trace.TracerProvider {
	if !config.TracerConfig.Enabled {
		return nil
	}

	traceExporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(fmt.Sprintf("%s:%d", config.TracerConfig.Host, config.TracerConfig.Port)),
			otlptracehttp.WithHeaders(map[string]string{
				"content-type": "application/json",
			}),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		logger.Fatal(ctx, err, "❌ Could not create tracer exporter")
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			traceExporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(config.ApplicationConfig.Name),
				semconv.ServiceInstanceIDKey.String(config.ApplicationConfig.URL),
			),
		),
	)

	otel.SetTracerProvider(traceProvider)

	return traceProvider
}
