package uptrace

import (
	"context"
	"log/slog"

	"github.com/peruri-dev/inatrace/models"
	"go.opentelemetry.io/otel/trace"

	"github.com/uptrace/uptrace-go/uptrace"
	uptracego "github.com/uptrace/uptrace-go/uptrace"
)

type UPTraceProvider struct {
	tp trace.TracerProvider
}

func (tp *UPTraceProvider) Tracer(name string, options ...trace.TracerOption) trace.Tracer {
	return tp.tp.Tracer(name, options...)
}

func (tp *UPTraceProvider) Shutdown(ctx context.Context) error {
	uptracego.Shutdown(ctx)
	return nil
}

func InitTracerUP(serviceName string, serviceVersion string) models.CustomTraceProvider {
	// Configure OpenTelemetry with sensible defaults. and call otel.SetTracerProvider()
	uptracego.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("<FIXME>"),

		uptracego.WithServiceName(serviceName),
		uptracego.WithServiceVersion(serviceVersion),
	)

	return &UPTraceProvider{uptrace.TracerProvider()}
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
