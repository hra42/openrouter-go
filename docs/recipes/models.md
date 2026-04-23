# Discovering models

```go
// All models
models, _ := client.ListModels(ctx)

// Endpoints (providers) for a specific model — includes pricing, context length, uptime
endpoints, _ := client.ListModelEndpoints(ctx, "openai/gpt-4o")

// All providers
providers, _ := client.ListProviders(ctx)

// ZDR-compliant endpoints
zdr, _ := client.ListZDREndpoints(ctx)
```

Filter `ListModels` by category (chat, completion, embedding, etc.) — see `models_endpoint.go`.

Examples:

- [`examples/list-models/`](../../examples/list-models/main.go)
- [`examples/list-providers/`](../../examples/list-providers/main.go)
- [`examples/model-endpoints/`](../../examples/model-endpoints/main.go)
