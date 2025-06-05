# inatrace/estrace

Integrations Inatrace with Elastic APM.

Possible configurations:
* Directly to Elastic APM Server via Elastic APM Agent (embedded)
* via OpenTelemetry Collector

## Directly to Elastic APM

```
import (
    "github.com/peruri-dev/inatrace/integrations/estrace"
)

tp := inatrace.InitTracerES(serviceName, servieVersion)

defer func() {
    if err := tp.Shutdown(context.Background()); err != nil {
        log.Printf("Error shutting down tracer provider: %v", err)
    }
}()
```

`InitTracerES` will register TraceProvider as default. The `Shutdown` actually noop.

App Env:
```
ELASTIC_APM_SERVICE_NAME=my-example
ELASTIC_APM_SECRET_TOKEN=<ELASTIC_APM_TOKEN>
ELASTIC_APM_SERVER_URL=<ELASTIC_APM_SERVER>:443
ELASTIC_APM_ENVIRONMENT=development
```

### Configuration

You need to setup Elastic APM Server then Elastic APM Agent (the collector and client).

**Elastic APM Server** is the server that capture the ingest. It is different product,
and differs from **ElasticSearch** and **Kibana**.

![apm-architecture-diy](https://www.elastic.co/guide/en/observability/8.15/images/apm-architecture-diy.png)

**Elastic APM Agent** is embedded via go package, just provide the env `ELASTIC_APM_`
alongside the app.

Follow these instructions: [add-apm-integration](https://www.elastic.co/guide/en/observability/8.15/traces-get-started.html#add-apm-integration)
