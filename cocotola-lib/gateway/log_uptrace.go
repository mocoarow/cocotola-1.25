package gateway

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func initLogExporterUptraceHTTP(ctx context.Context, logConfig *LogConfig) (sdklog.Exporter, error) {
	return otlploghttp.New(ctx, //nolint:wrapcheck
		otlploghttp.WithEndpoint(logConfig.Uptrace.Endpoint),
		otlploghttp.WithHeaders(map[string]string{
			"uptrace-dsn": logConfig.Uptrace.DSN,
		}),
		otlploghttp.WithCompression(otlploghttp.GzipCompression),
	)
}
