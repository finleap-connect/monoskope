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

package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/finleap-connect/monoskope/internal/version"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	otelEnabled        = "OTEL_ENABLED"
	serviceNamePrefix  = "OTEL_SERVICE_NAME_PREFIX"
	serviceName        = "OTEL_SERVICE_NAME"
	defaultNamePrefix  = "m8"
	otelEndpointEnvVar = "OTEL_EXPORTER_OTLP_ENDPOINT"
)

var (
	instanceKey = uuid.New().String()
)

// GetIsOpenTelemetryEnabled returns if environment variable OTEL_ENABLED is set to "true"
func GetIsOpenTelemetryEnabled() bool {
	return os.Getenv(otelEnabled) == "true"
}

// GetServiceName combines the values of the environment variables OTEL_SERVICE_NAME_PREFIX and OTEL_SERVICE_NAME with fallback to "m8" and "version.Name"
func GetServiceName() string {
	name := version.Name
	prefix := defaultNamePrefix

	if sn := os.Getenv(serviceName); sn != "" {
		name = sn
	}
	if sp := os.Getenv(serviceNamePrefix); sp != "" {
		prefix = sp
	}

	return fmt.Sprintf("%s%s", prefix, name)
}

// InitOpenTelemetry configures and sets the global MeterProvider and TracerProvider for OpenTelemetry
func InitOpenTelemetry(ctx context.Context) (func() error, error) {
	meterProviderShutdown, err := initMeterProvider(ctx)
	if err != nil {
		return nil, err
	}

	tracerProviderShutdown, err := initTracerProvider(ctx)
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

// GetTracer returns a new traces with the current service name set
func GetTracer() trace.Tracer {
	return otel.Tracer(GetServiceName())
}

// initMeterProvider configures and sets the global MeterProvider
func initMeterProvider(ctx context.Context) (func() error, error) {
	meterExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(meterExporter)))
	global.SetMeterProvider(meterProvider)

	return func() error { return meterProvider.Shutdown(ctx) }, nil
}

// initTracerProvider configures and sets the global TracerProvider
func initTracerProvider(ctx context.Context) (func() error, error) {
	var err error
	var spanExporter sdktrace.SpanExporter
	serviceName := GetServiceName()

	log := logger.WithName("telemetry").WithValues("serviceName", serviceName, "version", version.Version, "instance", instanceKey)

	endpoint, exists := os.LookupEnv(otelEndpointEnvVar)
	if !exists {
		return nil, fmt.Errorf("%s env var must be set", otelEndpointEnvVar)
	}

	timeoutContext, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	log.Info("Establishing connection to OpenTelemetry collector...", "endpoint", endpoint)
	conn, err := grpc.DialContext(timeoutContext, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err, "unable to connect to OpenTelemetry collector", "endpoint", endpoint)
		return nil, err
	}

	log.V(logger.DebugLevel).Info("Initializing otlptracegrpc...")
	spanExporter, err = otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, err
	}

	log.V(logger.DebugLevel).Info("Configuring resource...")
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(GetServiceName()),
			semconv.ServiceVersionKey.String(version.Version),
			semconv.ServiceInstanceIDKey.String(instanceKey),
		))

	if err != nil {
		return nil, err
	}

	options := []sdktrace.TracerProviderOption{
		sdktrace.WithBatcher(spanExporter),
		sdktrace.WithResource(res),
	}
	tracerProvider := sdktrace.NewTracerProvider(options...)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	log.Info("TraceProvider configured.")

	return func() error { return tracerProvider.Shutdown(ctx) }, nil
}
