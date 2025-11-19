package config

import (
	"context"
	"fmt"

	gcpexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracerExporterGoogle(_ context.Context, traceConfig *TraceConfig) (sdktrace.SpanExporter, error) {
	if traceConfig.Google == nil {
		return nil, fmt.Errorf("google trace configuration is required")
	}

	return gcpexporter.New(gcpexporter.WithProjectID(traceConfig.Google.ProjectID)) //nolint:wrapcheck
}
