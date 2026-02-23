package gateway

import (
	"context"
	"io"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

func initTracerExporterNone(_ context.Context, _ *TraceConfig) (*stdouttrace.Exporter, error) {
	return stdouttrace.New( //nolint:wrapcheck
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(io.Discard),
	)
}
