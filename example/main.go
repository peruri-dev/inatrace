package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/peruri-dev/inatrace/integrations/uptrace"
	// OR "github.com/peruri-dev/inatrace/integrations/estrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	os.Setenv("UPTRACE_DSN", "https://secret@api.uptrace.dev?grpc=4317")

	ctx := context.Background()

	tp := uptrace.InitTracer("example", "v0.1.0")
	// OR
	// tp := estrace.InitTracer("example", "v0.1.0")

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Create a tracer. Usually, tracer is a global variable.
	tracer := otel.Tracer("app_or_package_name")

	// Create a root span (a trace) to measure some operation.
	ctx, main := tracer.Start(ctx, "main-operation")
	// End the span when the operation we are measuring is done.
	defer main.End()

	// The passed ctx carries the parent span (main).
	// That is how OpenTelemetry manages span relations.
	_, child1 := tracer.Start(ctx, "GET /posts/:id")
	child1.SetAttributes(
		attribute.String("http.method", "GET"),
		attribute.String("http.route", "/posts/:id"),
		attribute.String("http.url", "http://localhost:8080/posts/123"),
		attribute.Int("http.status_code", 200),
	)
	if err := errors.New("dummy error"); err != nil {
		child1.RecordError(err, trace.WithStackTrace(true))
		child1.SetStatus(codes.Error, err.Error())
		child1.End()
	}

	_, child2 := tracer.Start(ctx, "SELECT")
	child2.SetAttributes(
		attribute.String("db.system", "mysql"),
		attribute.String("db.statement", "SELECT * FROM posts LIMIT 100"),
	)
	child2.End()
}
