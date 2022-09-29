// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracing

import (
	"context"

	"github.com/finleap-connect/monoskope/internal/version"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

/*
Find detailed documentation at https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace#readme-examples.
*/

// InitOpenTelemetry configures and sets the global MeterProvider and TracerProvider for OpenTelemetry
func InitOpenTelemetry(ctx context.Context) (func() error, error) {
	meterProviderShutdown, err := InitMeterProvider(ctx)
	if err != nil {
		return nil, err
	}

	tracerProviderShutdown, err := InitTracerProvider(ctx)
	if err != nil {
		return nil, err
	}

	return func() error {
		if err := meterProviderShutdown(); err != nil {
			return err
		}
		if err := tracerProviderShutdown(); err != nil {
			return err
		}
		return nil
	}, nil
}

// InitMeterProvider configures and sets the global MeterProvider
func InitMeterProvider(ctx context.Context) (func() error, error) {
	meterExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(meterExporter)))
	global.SetMeterProvider(meterProvider)

	return func() error { return meterProvider.Shutdown(ctx) }, nil
}

// InitTracerProvider configures and sets the global TracerProvider
func InitTracerProvider(ctx context.Context) (func() error, error) {
	client := otlptracegrpc.NewClient()
	traceExporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(version.Name),
			semconv.ServiceVersionKey.String(version.Version),
			semconv.ServiceInstanceIDKey.String(uuid.New().String()),
		),
	)
	if err != nil {
		return nil, err
	}
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() error { return tracerProvider.Shutdown(ctx) }, nil
}
