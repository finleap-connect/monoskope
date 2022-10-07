package main

import (
	"context"
	"errors"
	"os"

	"github.com/finleap-connect/monoskope/internal/telemetry"
	"go.opentelemetry.io/otel/codes"
)

func test() {
	_, span := telemetry.GetTracer().Start(context.Background(), "test-span")
	defer span.End()
	span.RecordError(errors.New("just a test"))
	span.SetStatus(codes.Error, "ja kapott")
}

func main() {
	os.Setenv("OTEL_ENABLED", "true")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	shutdown, err := telemetry.InitOpenTelemetry(context.Background())
	if err != nil {
		panic(err)
	}
	defer shutdown()
	test()
}
