# Streaming

The SDK uses Server-Sent Events. Streams expose a channel via `Events()` and must be closed.

```go
stream, err := client.ChatCompleteStream(ctx,
    []openrouter.Message{openrouter.CreateUserMessage("Write a haiku.")},
    openrouter.WithModel("openai/gpt-4o-mini"),
)
if err != nil {
    return err
}
defer func() { _ = stream.Close() }()

for event := range stream.Events() {
    for _, choice := range event.Choices {
        if choice.Delta != nil {
            if s, ok := choice.Delta.Content.(string); ok {
                fmt.Print(s)
            }
        }
    }
}
if err := stream.Err(); err != nil {
    return err
}
```

## Cancellation

Cancelling the `context.Context` interrupts the stream cleanly:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

stream, _ := client.ChatCompleteStream(ctx, messages, openrouter.WithModel(m))
defer stream.Close()
for event := range stream.Events() { /* ... */ }
```

## Legacy (prompt) streaming

```go
stream, _ := client.CompleteStream(ctx, "Once upon a time",
    openrouter.WithCompletionModel("meta-llama/llama-3.1-8b-instruct"),
)
defer stream.Close()
for event := range stream.Events() {
    for _, choice := range event.Choices {
        fmt.Print(choice.Text)
    }
}
```

## Pitfalls

- Always `defer stream.Close()` right after checking the returned error.
- Read `stream.Err()` **after** the `Events()` channel closes.
- Don't `break` out of the loop without cancelling the context — you'll leak the HTTP connection.

See [`examples/streaming/main.go`](../../examples/streaming/main.go).
