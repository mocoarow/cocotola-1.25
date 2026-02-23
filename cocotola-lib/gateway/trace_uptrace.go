package gateway

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func initTracerExporterUptraceHTTP(ctx context.Context, traceConfig *TraceConfig) (*otlptrace.Exporter, error) {
	if traceConfig.Uptrace == nil {
		return nil, errors.New("uptrace trace configuration is required")
	}

	return otlptracehttp.New(ctx, //nolint:wrapcheck
		otlptracehttp.WithEndpoint(traceConfig.Uptrace.Endpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"uptrace-dsn": traceConfig.Uptrace.DSN,
		}),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
}
