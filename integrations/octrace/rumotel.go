package octrace

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// for datatace using @datadog/browser-rum
// for truewatch using @cloudcate/browser-rum
// datadog trace-id 64-bit unsigned decimal integer
// datadog span-id 64-bit unsigned decimal integer
// otel trace-id 128-bit lowercase hexadecimal (32 chars)
// otel span-id 64-bit lowercase hexadecimal (16 chars)

// DatadogRUMPropagator extracts Datadog RUM trace headers and converts them to OpenTelemetry trace context
type DatadogRUMPropagator struct{}

var _ propagation.TextMapPropagator = &DatadogRUMPropagator{}

func ddTraceIDtoOtel(ctx context.Context, traceID string) (trace.TraceID, error) {
	// Convert Datadog trace ID (64-bit unsigned int as string) to OpenTelemetry TraceID
	traceIDUint, err := strconv.ParseUint(traceID, 10, 64)
	if err != nil {
		slog.DebugContext(ctx, "Failed to parse Datadog trace ID", slog.String("error", err.Error()))
		return trace.TraceID{}, err
	}

	// Convert to 16-character hex string (128-bit), pad with zeros on the left
	traceIDHex := fmt.Sprintf("%016x", traceIDUint)

	// Parse as OpenTelemetry TraceID (128-bit)
	otelTraceID, err := trace.TraceIDFromHex(traceIDHex + strings.Repeat("0", 16))
	if err != nil {
		slog.DebugContext(ctx, "Failed to convert to OTel TraceID", slog.String("error", err.Error()))
		return trace.TraceID{}, err
	}

	return otelTraceID, nil
}

func ddSpanIDtoOtel(parentSpanID string) trace.SpanID {
	if parentSpanID != "" {

		parentSpanIDUint, err := strconv.ParseUint(parentSpanID, 10, 64)
		if err == nil {
			spanIDHex := fmt.Sprintf("%016x", parentSpanIDUint)
			otelSpanID, _ := trace.SpanIDFromHex(spanIDHex)
			return otelSpanID
		}

	}

	return trace.SpanID{}

}

// Extract extracts Datadog trace context from headers
func (p *DatadogRUMPropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	traceID := carrier.Get("x-datadog-trace-id")
	parentSpanID := carrier.Get("x-datadog-parent-id")

	if traceID != "" {
		// Convert Datadog trace ID (64-bit unsigned int as string) to OpenTelemetry TraceID
		otelTraceID, err := ddTraceIDtoOtel(ctx, traceID)
		if err != nil {
			return ctx
		}

		// Parse parent span ID if available
		otelSpanID := ddSpanIDtoOtel(parentSpanID)

		// Create span context with the extracted IDs
		sc := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    otelTraceID,
			SpanID:     otelSpanID,
			TraceFlags: trace.FlagsSampled,
			Remote:     true,
		})

		// slog.DebugContext(ctx, "Allow propagate trace id from datadog rum", slog.String("x-datadog-trace-id", traceID))

		return trace.ContextWithSpanContext(ctx, sc)
	}

	traceParent := carrier.Get("traceparent")
	traceState := carrier.Get("tracestate")

	if traceParent != "" {
		// TODO:
		_ = traceState
	}

	return ctx
}

// Inject injects Datadog trace headers into the carrier
func (p *DatadogRUMPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()

	if !sc.IsValid() {
		return
	}

	// Convert OpenTelemetry TraceID to Datadog format (lower 64 bits as unsigned int)
	traceID := sc.TraceID()
	traceIDLower := uint64(0)
	for i := 0; i < 8; i++ {
		traceIDLower = traceIDLower<<8 + uint64(traceID[i])
	}

	// Convert SpanID to Datadog format
	spanID := sc.SpanID()
	spanIDUint := uint64(0)
	for i := 0; i < 8; i++ {
		spanIDUint = spanIDUint<<8 + uint64(spanID[i])
	}

	carrier.Set("x-datadog-trace-id", strconv.FormatUint(traceIDLower, 10))
	carrier.Set("x-datadog-parent-id", strconv.FormatUint(spanIDUint, 10))
}

// Fields returns the list of headers used by this propagator
func (p *DatadogRUMPropagator) Fields() []string {
	return []string{"x-datadog-trace-id", "x-datadog-parent-id"}
}
