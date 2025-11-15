package config

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracerExporterUptraceHTTP(ctx context.Context, traceConfig *TraceConfig) (sdktrace.SpanExporter, error) {
	return otlptracehttp.New(ctx, //nolint:wrapcheck
		otlptracehttp.WithEndpoint(traceConfig.Uptrace.Endpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"uptrace-dsn": traceConfig.Uptrace.DSN,
		}),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
}
