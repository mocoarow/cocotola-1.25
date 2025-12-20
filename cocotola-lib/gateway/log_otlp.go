package gateway

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func initLogExporterOTLPHTTP(ctx context.Context, logConfig *LogConfig) (sdklog.Exporter, error) {
	options := make([]otlploghttp.Option, 0)
	options = append(options, otlploghttp.WithEndpoint(logConfig.OTLP.Endpoint))
	if logConfig.OTLP.Insecure {
		options = append(options, otlploghttp.WithInsecure())
	}

	return otlploghttp.New(ctx, options...) //nolint:wrapcheck
}
