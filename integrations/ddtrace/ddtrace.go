package ddtrace

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"go.opentelemetry.io/otel"

	//"go.opentelemetry.io/otel/exporters/jaeger"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	ddotel "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentelemetry"

	"github.com/peruri-dev/inatrace/models"
)

// DDTraceProvider
type DDTraceProvider struct {
	tp *ddotel.TracerProvider
}

func (tp *DDTraceProvider) Tracer(name string, options ...trace.TracerOption) trace.Tracer {
	return tp.tp.Tracer(name, options...)
}

func (tp *DDTraceProvider) Shutdown(ctx context.Context) error {
	return tp.tp.Shutdown()
}

func InitTracerDD() models.CustomTraceProvider {
	tp := ddotel.NewTracerProvider()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return &DDTraceProvider{tp}
}

// ExtractTraceSpanID, refers to:
// https://docs.datadoghq.com/tracing/other_telemetry/connect_logs_and_traces/opentelemetry/?tab=go
func ExtractTraceSpanID(ctx context.Context) []slog.Attr {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}

	attrs := []slog.Attr{}
	spanCtx := span.SpanContext()
	ddTraceID := ""
	ddSpanID := ""

	if spanCtx.HasTraceID() {
		ddTraceID = convertTraceID(spanCtx.TraceID().String())
	}

	if spanCtx.HasSpanID() {
		ddSpanID = convertTraceID(spanCtx.SpanID().String())
	}

	attrs = append(attrs, slog.Group("dd",
		"trace_id", ddTraceID,
		"span_id", ddSpanID,
		"service", os.Getenv("DD_SERVICE"),
		"env", os.Getenv("DD_ENV"),
		"version", os.Getenv("DD_VERSION"),
	))

	return attrs
}

func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
