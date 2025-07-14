package octrace

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/peruri-dev/inatrace/models"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type OCTraceProvider struct {
	tp trace.TracerProvider
}

func (tp *OCTraceProvider) Tracer(name string, options ...trace.TracerOption) trace.Tracer {
	return tp.tp.Tracer(name, options...)
}

func (tp *OCTraceProvider) Shutdown(ctx context.Context) error {
	return nil
}

func InitTracerOC(serviceName string, serviceVersion string) models.CustomTraceProvider {
	headers := map[string]string{
		"content-type": "application/json",
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("localhost:4318"),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		log.Fatalf("creating new exporter: %w", err)
	}

	tracerprovider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
		),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		),
	)

	otel.SetTracerProvider(tracerprovider)
	//otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	fmt.Println(exporter)

	return &OCTraceProvider{tracerprovider}
}

// ExtractTraceSpanID
func ExtractTraceSpanID(ctx context.Context) []slog.Attr {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}

	attrs := []slog.Attr{}
	spanCtx := span.SpanContext()
	traceID := ""
	spanID := ""

	if spanCtx.HasTraceID() {
		traceID = spanCtx.TraceID().String()
	}

	if spanCtx.HasSpanID() {
		spanID = spanCtx.SpanID().String()
	}

	attrs = append(attrs, slog.String("trace.id", traceID), slog.String("transaction.id", spanID))

	return attrs
}
