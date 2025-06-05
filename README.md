# inatrace

Simplify APM trace using OpenTelemetry

## Fiber

```
import (
    "github.com/gofiber/contrib/otelfiber"
    "github.com/gofiber/fiber/v2"
)

f := fiber.New(fiber.Config{})
f.Use(otelfiber.Middleware())
```

### Backends

* [DataDog](./integrations/ddtrace/)
* [Elastic](./integrations/estrace/)

## Resty

```
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