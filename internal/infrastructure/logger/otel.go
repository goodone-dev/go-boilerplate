package logger

import (
	"context"
	"fmt"
	"log"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otelsdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func NewProvider(ctx context.Context) *otelsdklog.LoggerProvider {
	if !config.LoggerConfig.Enabled {
		return nil
	}

	exporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(fmt.Sprintf("%s:%d", config.LoggerConfig.Host, config.LoggerConfig.Port)),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		log.Fatal("❌ Could not create logger exporter", err)
	}

	provider := otelsdklog.NewLoggerProvider(
		otelsdklog.WithProcessor(
			otelsdklog.NewBatchProcessor(exporter),
		),
		otelsdklog.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(config.ApplicationConfig.Name),
				semconv.ServiceInstanceIDKey.String(config.ApplicationConfig.URL),
			),
		),
	)

	oLogger = provider.Logger(config.ApplicationConfig.Name)
	return provider
}
