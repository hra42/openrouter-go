# Web search

Augment responses with real-time web data. Two engines are available: the model's native search (if supported) or Exa. They price differently — know which you're asking for.

```go
resp, err := client.ChatComplete(ctx,
    []openrouter.Message{openrouter.CreateUserMessage("What shipped in Go 1.26?")},
    openrouter.WithModel("openai/gpt-4o:online"), // ":online" enables native web search
)
```

Or configure explicitly via plugin options — see [`examples/web_search/main.go`](../../examples/web_search/main.go) and `web_search.go` for `WebSearchEngine` / `WebSearchContextSize`.

Tip: the `:online` suffix is a model-name modifier. Concatenate it; don't pass it as a separate option.
