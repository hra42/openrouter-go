# Broadcast webhook (OTLP)

OpenRouter's Broadcast feature posts OTLP-JSON trace payloads to your webhook. The SDK ships standalone parsing utilities in `broadcast.go` — no `Client` needed.

Typical HTTP handler:

```go
http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    traces, err := openrouter.ParseBroadcastTraces(body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    for _, t := range traces {
        log.Printf("trace: %+v", t)
    }
    w.WriteHeader(http.StatusOK)
})
```

For the raw OTLP wire format use `ParseBroadcastPayload`, which returns `*OTLPExportTraceRequest`.

See [`examples/broadcast-webhook/main.go`](../../examples/broadcast-webhook/main.go) and `broadcast_models.go` for the full types.
