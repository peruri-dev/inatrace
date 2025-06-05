module github.com/peruri-dev/inatrace/integrations/estrace

go 1.23.2

require (
	github.com/peruri-dev/inatrace v1.0.0
	go.elastic.co/apm/module/apmotel/v2 v2.6.2
	go.elastic.co/apm/v2 v2.6.2
	go.opentelemetry.io/otel v1.31.0
	go.opentelemetry.io/otel/trace v1.31.0
)

require (
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/elastic/go-sysinfo v1.15.0 // indirect
	github.com/elastic/go-windows v1.0.2 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	go.elastic.co/apm/module/apmhttp/v2 v2.6.2 // indirect
	go.elastic.co/fastjson v1.4.0 // indirect
	go.opentelemetry.io/otel/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/sdk v1.31.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.31.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	howett.net/plist v1.0.1 // indirect
)

replace github.com/peruri-dev/inatrace => ../..
