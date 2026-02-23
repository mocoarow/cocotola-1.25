package gateway

import (
	"context"
	"errors"

	gcpexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
)

func initTracerExporterGoogle(_ context.Context, traceConfig *TraceConfig) (*gcpexporter.Exporter, error) {
	if traceConfig.Google == nil {
		return nil, errors.New("google trace configuration is required")
	}

	return gcpexporter.New(gcpexporter.WithProjectID(traceConfig.Google.ProjectID)) //nolint:wrapcheck
}
