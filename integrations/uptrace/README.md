# inatrace/uptrace

Integrations to Uptrace backend.

App Env:
```
UPTRACE_DSN=https://secret@api.uptrace.dev?grpc=4317
```

Changes v1.0.3:

* auto set env UPTRACE_DISABLED when UPTRACE_DSN empty
* user span_id key for otel's spanID

