package gateway

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

func initTracerExporterStdout(_ context.Context, _ *TraceConfig) (*stdouttrace.Exporter, error) {
	return stdouttrace.New( //nolint:wrapcheck
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(os.Stderr),
	)
}
