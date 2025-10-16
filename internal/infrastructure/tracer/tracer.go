package tracer

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-logr/stdr"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func NewProvider(ctx context.Context) *trace.TracerProvider {
	logger := stdr.New(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile))
	otel.SetLogger(logger)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

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
		log.Fatalf("❌ Could not to create tracer exporter: %v", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			traceExporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(config.ApplicationConfig.Name),
				semconv.ServiceInstanceIDKey.String(config.ApplicationConfig.URL),
			),
		),
	)

	return traceProvider
}
