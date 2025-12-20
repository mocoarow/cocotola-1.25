package gateway

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracerExporterUptraceHTTP(ctx context.Context, traceConfig *TraceConfig) (sdktrace.SpanExporter, error) {
	if traceConfig.Uptrace == nil {
		return nil, fmt.Errorf("uptrace trace configuration is required")
	}

	return otlptracehttp.New(ctx, //nolint:wrapcheck
		otlptracehttp.WithEndpoint(traceConfig.Uptrace.Endpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"uptrace-dsn": traceConfig.Uptrace.DSN,
		}),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
}
