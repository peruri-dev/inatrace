package estrace

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"log"

	"go.elastic.co/apm/module/apmotel/v2"
	"go.elastic.co/apm/v2"
	"go.opentelemetry.io/otel"

	//"go.opentelemetry.io/otel/exporters/jaeger"

	"github.com/peruri-dev/inatrace/models"
)

// ESTraceProvider
type ESTraceProvider struct {
	tp trace.TracerProvider
}

func (tp *ESTraceProvider) Tracer(name string, options ...trace.TracerOption) trace.Tracer {
	return tp.tp.Tracer(name, options...)
}

func (tp *ESTraceProvider) Shutdown(ctx context.Context) error {
	return nil
}

func InitTracerES(serviceName string, serviceVersion string) models.CustomTraceProvider {
	apmTracer, err := apm.NewTracer(serviceName, serviceVersion)
	if err != nil {
		log.Fatal(err)
	}
	provider, err := apmotel.NewTracerProvider(apmotel.WithAPMTracer(
		apmTracer,
	))
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return &ESTraceProvider{provider}
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
