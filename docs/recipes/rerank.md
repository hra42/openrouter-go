# Rerank

Score a set of documents by relevance to a query. Useful as a second stage after vector retrieval.

```go
resp, err := client.Rerank(ctx,
    "how do I stream chat completions in Go?",
    []string{
        "Streaming uses Server-Sent Events and requires closing the stream.",
        "To list available models, call ListModels.",
        "The SDK has no external dependencies.",
    },
    "cohere/rerank-v3.5",
)
if err != nil { return err }

for _, r := range resp.Results {
    fmt.Printf("idx=%d score=%.3f\n", r.Index, r.RelevanceScore)
}
```

Results are returned sorted by relevance. `r.Index` is the position in your original `documents` slice.

See [`examples/rerank/main.go`](../../examples/rerank/main.go) and `rerank_options.go` for additional options (top-N, return documents in the response, etc.).
