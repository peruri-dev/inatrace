# inatrace/ddtrace

Integrations Inatrace with DataDog Backend.

Possible configurations:
* Directly to Datadog Agent
* via OpenTelemetry Collector with Datadog Exporter

## Directly to Datadog Agent

```
import (
    "github.com/peruri-dev/inatrace/integrations/ddtrace"
)

tp := inatrace.InitTracerDD()

defer func() {
    if err := tp.Shutdown(context.Background()); err != nil {
        log.Printf("Error shutting down tracer provider: %v", err)
    }
}()
```

`InitTracerDD` will register TraceProvider as default.

App Env:
```
DD_SERVICE=example
OTEL_EXPORTER_OTLP_ENDPOINT=http://<AGENT_IP_ADDRESS>:4318
```

Agent Env:
```
DD_API_KEY=<DATADOG-API-KEY>
DD_SITE=datadoghq.eu
DD_APM_ENABLED=true
DD_APM_NON_LOCAL_TRAFFIC=true
DD_OTLP_CONFIG_RECEIVER_PROTOCOLS_HTTP_ENDPOINT=0.0.0.0:4318
```
