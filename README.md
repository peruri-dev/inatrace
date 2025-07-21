# inatrace

Simplify APM trace using OpenTelemetry

## Fiber

```go
import (
    "github.com/gofiber/contrib/otelfiber"
    "github.com/gofiber/fiber/v2"
)

f := fiber.New(fiber.Config{})
f.Use(otelfiber.Middleware())
```

## Backends

* [DataDog](./integrations/ddtrace/)
* [Elastic](./integrations/estrace/)
* [Jaeger via OtelCollector](./integrations/octrace/)
* [Uptrace](./integrations/uptrace/)

## Enabling Trace

### Function Span

```go
import (
    "github.com/peruri-dev/inatrace"
    "go.opentelemetry.io/otel/trace"
)

func SendNotif(ctx context.Context, id string) {
    _, span := inatrace.Start(ctx, "SendNotif", trace.WithAttributes(attribute.String("id", id)))
    defer span.End()
}
```

### Resty

Resty is HTTP Client.

```go
import (
    "github.com/dubonzi/otelresty"
    "github.com/go-resty/resty/v2"
)

r = resty.New()
opts := []otelresty.Option{
    otelresty.WithTracerName("example-resty")
}
otelresty.TraceClient(r, opts...)

```

### DB with uptrace/bun

Bun is Database ORM.

```go
import (
    "github.com/uptrace/bun"
    "github.com/uptrace/bun/extra/bunotel"
)

db := bun.NewDB(...)
db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName("example-db")))
```

### Redis

Redis is In-Memory cache.

```go
import (
    "github.com/redis/go-redis/v9"
    "github.com/redis/go-redis/extra/redisotel/v9"
)

rdb := redis.NewClient(...)
redisotel.InstrumentTracing(rdb)
```
