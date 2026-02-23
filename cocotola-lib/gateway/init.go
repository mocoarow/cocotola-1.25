package gateway

import (
	"context"

	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type InitLogExporterFunc func(ctx context.Context, logConfig *LogConfig) (sdklog.Exporter, error)

// var initTracerExporters map[string]InitTracerExporterFunc
var initLogExporters map[string]InitLogExporterFunc

func init() {
	initLogExporters = map[string]InitLogExporterFunc{
		"otlphttp":    initLogExporterOTLPHTTP,
		"uptracehttp": initLogExporterUptraceHTTP,
	}
}
