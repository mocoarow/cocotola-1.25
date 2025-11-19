package config

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracerExporterOTLPHTTP(ctx context.Context, traceConfig *TraceConfig) (sdktrace.SpanExporter, error) {
	if traceConfig.OTLP == nil {
		return nil, fmt.Errorf("otlp trace configuration is required")
	}

	options := make([]otlptracehttp.Option, 0)
	options = append(options, otlptracehttp.WithEndpoint(traceConfig.OTLP.Endpoint))
	if traceConfig.OTLP.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	}

	return otlptracehttp.New(ctx, options...) //nolint:wrapcheck
}

func initTracerExporterOTLPgRPC(ctx context.Context, traceConfig *TraceConfig) (sdktrace.SpanExporter, error) {
	if traceConfig.OTLP == nil {
		return nil, fmt.Errorf("otlp trace configuration is required")
	}

	options := make([]otlptracegrpc.Option, 0)
	options = append(options, otlptracegrpc.WithEndpoint(traceConfig.OTLP.Endpoint))
	if traceConfig.OTLP.Insecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	return otlptracegrpc.New(ctx, options...) //nolint:wrapcheck
}
