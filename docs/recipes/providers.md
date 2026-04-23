# Provider preferences

Pin routing to specific providers, require parameter support, opt out of data collection, and control fallbacks:

```go
provider := openrouter.Provider{
    Order:             []string{"OpenAI", "Anthropic"},
    RequireParameters: true,
    DataCollection:    "deny",   // "allow" | "deny"
    AllowFallbacks:    true,
}

resp, _ := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithProvider(provider),
)
```

For per-request Zero Data Retention enforcement, use `ListZDREndpoints` to discover ZDR-compliant providers before making calls.

See `models.go` for the full `Provider` struct.
