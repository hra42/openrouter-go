# Transforms (context-window management)

Automatically compress the prompt when it exceeds the model's context window:

```go
resp, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4o-mini"),
    openrouter.WithTransforms("middle-out"),
)
```

`"middle-out"` drops content from the middle of the message history. Useful for long chats, but **silently lossy** — don't enable by default in code that needs deterministic inputs.

See [`examples/transforms/main.go`](../../examples/transforms/main.go).
