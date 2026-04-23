# Responses API (beta)

> **Warning:** The Responses API is in beta and may change without notice. Prefer `ChatComplete` for production workloads.

OpenAI-compatible stateless Responses API: reasoning, tool calling, web search, streaming.

```go
resp, err := client.CreateResponse(ctx,
    "Explain quantum entanglement.",              // input (string or structured)
    openrouter.WithResponsesModel("openai/o1"),
    openrouter.WithResponsesReasoningEffort("medium"),
)
```

Streaming:

```go
stream, _ := client.CreateResponseStream(ctx, "...",
    openrouter.WithResponsesModel("openai/o1"),
)
defer stream.Close()
for ev := range stream.Events() {
    // ev is a typed Responses stream event
}
if err := stream.Err(); err != nil { return err }
```

See [`examples/responses/main.go`](../../examples/responses/main.go) and `responses_options.go` for the full option set (tools, web search, metadata, max output tokens, etc.).
