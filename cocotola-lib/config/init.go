package config

import (
	"context"
	"database/sql"
	"log/slog"

	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"

	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
)

type InitTracerExporterFunc func(ctx context.Context, traceConfig *TraceConfig) (sdktrace.SpanExporter, error)
type InitLogExporterFunc func(ctx context.Context, logConfig *LogConfig) (sdklog.Exporter, error)

var initTracerExporters map[string]InitTracerExporterFunc
var initLogExporters map[string]InitLogExporterFunc

type InitDBFunc func(context.Context, *DBConfig, slog.Level, string) (libgateway.DialectRDBMS, *gorm.DB, *sql.DB, error)

var initDBs map[string]InitDBFunc

func init() {
	initTracerExporters = map[string]InitTracerExporterFunc{
		"google":      initTracerExporterGoogle,
		"otlphttp":    initTracerExporterOTLPHTTP,
		"otlpgrpc":    initTracerExporterOTLPgRPC,
		"none":        initTracerExporterNone,
		"stdout":      initTracerExporterStdout,
		"uptracehttp": initTracerExporterUptraceHTTP,
	}
	initLogExporters = map[string]InitLogExporterFunc{
		"otlphttp":    initLogExporterOTLPHTTP,
		"uptracehttp": initLogExporterUptraceHTTP,
	}
	initDBs = map[string]InitDBFunc{
		"mysql":   initDBMySQL,
		"sqlite3": initDBSQLite3,
	}
}
