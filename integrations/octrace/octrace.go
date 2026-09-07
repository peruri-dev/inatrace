package octrace

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	inatraceModels "github.com/peruri-dev/inatrace/models"

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

type ExporterFormat string

const (
	FORMAT_HTTP ExporterFormat = "http"
	FORMAT_GRPC ExporterFormat = "grpc"
)

func InitTracerOC(serviceName string, serviceVersion string, exporterFormat ExporterFormat) inatraceModels.CustomTraceProvider {
	// OTEL_EXPORTER_OTLP_HEADERS=content-type=application/json
	// OTEL_EXPORTER_OTLP_ENDPOINT=https://localhost:4318
	// OTEL_EXPORTER_OTLP_INSECURE=true (if using http)
	// OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=https://localhost:4318/v1/traces
	// OTEL_EXPORTER_OTLP_TRACES_INSECURE=true (if using http)

	var client otlptrace.Client

	switch exporterFormat {
	case FORMAT_HTTP, "":
		slog.Info("[OCtrace] using exporter http format")
		client = otlptracehttp.NewClient()
	case FORMAT_GRPC:
		slog.Info("[OCtrace] using exporter grpc format")
		client = otlptracegrpc.NewClient()
	default:
		slog.Error(fmt.Sprintf("[OCtrace] undefined format %s", exporterFormat))
		log.Fatalf("[OCtrace] undefined format %s", exporterFormat)
	}

	exporter, err := otlptrace.New(
		context.Background(),
		client,
	)
	if err != nil {
		slog.Error("[OCtrace] fail to creating new exporter", slog.String("error", err.Error()))
		log.Fatalf("[OCtrace] fail to creating new exporter, error: %s", err)
	}

	kv := []attribute.KeyValue{}

	if serviceName != "" {
		kv = append(kv, semconv.ServiceNameKey.String(serviceName))
		slog.Info(fmt.Sprintf("[OCtrace] using provided serivcename %s", serviceName))
	} else {
		slog.Info(fmt.Sprintf("[OCtrace] using default, maybe: %s", os.Getenv("OTEL_SERVICE_NAME")))
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
				kv...,
			),
		),
	)

	otel.SetTracerProvider(tracerprovider)
	//otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return &OCTraceProvider{tracerprovider}
}
