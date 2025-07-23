package logging

import (
	"context"
	"log/slog"
	"os"

	oteltrace "go.opentelemetry.io/otel/trace"
)

// InitLogger initializes the structured logger.
func InitLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

// ContextLogger returns a logger with trace and span IDs from the context.
func ContextLogger(ctx context.Context) *slog.Logger {
	spanCtx := oteltrace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return slog.Default().With(
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)
	}
	return slog.Default()
}
