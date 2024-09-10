package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracing() (*trace.TracerProvider, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	exporter, err := otlptrace.New(context.Background(), otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithHeaders(headers),
		otlptracehttp.WithInsecure(),
	))
	if err != nil {
		return nil, err
	}
	tracerprovider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("goorm-class"),
			),
		),
	)

	otel.SetTracerProvider(tracerprovider)

	return tracerprovider, nil
}
