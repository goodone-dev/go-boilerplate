package logger

import (
	"context"
	"fmt"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var oLogger otellog.Logger

func NewProvider(ctx context.Context) *log.LoggerProvider {
	if !config.LoggerConfig.Enabled {
		return nil
	}

	logExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(fmt.Sprintf("%s:%d", config.LoggerConfig.Host, config.LoggerConfig.Port)),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(
			log.NewBatchProcessor(logExporter),
		),
		log.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(config.ApplicationConfig.Name),
				semconv.ServiceInstanceIDKey.String(config.ApplicationConfig.URL),
			),
		),
	)

	oLogger = loggerProvider.Logger("")

	return loggerProvider
}
