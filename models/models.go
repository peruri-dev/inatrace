package models

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type CustomTraceProvider interface {
	Tracer(name string, options ...trace.TracerOption) trace.Tracer
	Shutdown(ctx context.Context) error
}
