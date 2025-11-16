package config

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/contrib/processors/baggagecopy"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

type OTLPTraceConfig struct {
	Endpoint string `yaml:"endpoint" validate:"required"`
	Insecure bool   `yaml:"insecure"`
}

type GoogleTraceConfig struct {
	ProjectID string `yaml:"projectId" validate:"required"`
}

type UptraceTraceConfig struct {
	Endpoint string `yaml:"endpoint" validate:"required"`
	DSN      string `yaml:"dsn" validate:"required"`
}

type TraceConfig struct {
	Exporter           string              `yaml:"exporter" validate:"required,oneof=otlphttp uptracehttp google none"`
	SamplingPercentage int                 `yaml:"samplingPercentage" validate:"gte=0,lte=100"`
	OTLP               *OTLPTraceConfig    `yaml:"otlp"`
	Google             *GoogleTraceConfig  `yaml:"google"`
	Uptrace            *UptraceTraceConfig `yaml:"uptrace"`
}

const traceShutdownTimeout = 5 * time.Second

func initTracerExporter(ctx context.Context, traceConfig *TraceConfig) (sdktrace.SpanExporter, error) {
	initTracerExporter, ok := initTracerExporters[traceConfig.Exporter]
	if !ok {
		return nil, fmt.Errorf("invalid trace exporter: %s", traceConfig.Exporter)
	}

	return initTracerExporter(ctx, traceConfig)
}

func initTraceSampler(samplingPercentage int) sdktrace.Sampler {
	if samplingPercentage >= 100 {
		return sdktrace.AlwaysSample()
	}
	if samplingPercentage <= 0 {
		return sdktrace.NeverSample()
	}
	return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(float64(samplingPercentage) / 100.0))
}

func InitTracerProvider(ctx context.Context, traceConfig *TraceConfig, appName string) (func(), error) {
	exp, err := initTracerExporter(ctx, traceConfig)
	if err != nil {
		return nil, fmt.Errorf("initTracerExporter: %w", err)
	}

	sampler := initTraceSampler(traceConfig.SamplingPercentage)

	bp := sdktrace.NewBatchSpanProcessor(exp,
		sdktrace.WithMaxQueueSize(10_000),
		sdktrace.WithMaxExportBatchSize(10_000),
		sdktrace.WithExportTimeout(10*time.Second),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bp),
		sdktrace.WithSpanProcessor(baggagecopy.NewSpanProcessor(baggagecopy.AllowAllMembers)),
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), traceShutdownTimeout)
		defer cancel()

		if err := tp.ForceFlush(shutdownCtx); err != nil {
			slog.Default().Error("failed to force flush tracer provider", slog.Any("error", err))
		}
		if err := bp.Shutdown(shutdownCtx); err != nil {
			slog.Default().Error("failed to shutdown span processor", slog.Any("error", err))
		}
	}, nil
}
