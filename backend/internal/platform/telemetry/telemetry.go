package telemetry

import (
	"context"
	"os"
	"log/slog"

	"github.com/gin-gonic/gin"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	otel "go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	otlptracegrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitTracer initializes the OpenTelemetry tracer provider.
func InitTracer() func(context.Context) error {
	// Create the OTLP gRPC exporter
	conn, err := grpc.NewClient(
		"otel-collector:4317", // OTLP gRPC receiver address
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		slog.Error("failed to create gRPC connection to collector", "error", err)
		os.Exit(1)
	}

	exporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		slog.Error("failed to create trace exporter", "error", err)
		os.Exit(1)
	}

	// Create a new tracer provider with the batch span processor.
	// The SimpleSpanProcessor is for demonstration purposes. Use BatchSpanProcessor in production.
	bsp := oteltrace.NewBatchSpanProcessor(exporter)
	tracerProvider := oteltrace.NewTracerProvider(
		oteltrace.WithSampler(oteltrace.AlwaysSample()),
		oteltrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set the global error handler
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		slog.Error("OpenTelemetry error", "error", err)
	}))

	return tracerProvider.Shutdown
}

// GinMiddleware returns a Gin middleware for OpenTelemetry tracing.
func GinMiddleware() gin.HandlerFunc {
	return otelgin.Middleware(os.Getenv("SERVICE_NAME"))
}
