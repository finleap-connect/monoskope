package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitOpenTelemetry creates a new MeterProvider and TracerProvider for OpenTelemetry from env vars
func InitOpenTelemetry(ctx context.Context) (func() error, error) {
	res, err := resource.New(ctx,
		resource.WithFromEnv(), // pull attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables
		resource.WithProcess(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithContainerID(),
	)
	if err != nil {
		return nil, err
	}

	meterExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(meterExporter)))
	global.SetMeterProvider(meterProvider)

	client := otlptracegrpc.NewClient()
	traceExporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() error {
		if err := meterProvider.Shutdown(ctx); err != nil {
			return err
		}
		if err := tracerProvider.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}, nil
}
