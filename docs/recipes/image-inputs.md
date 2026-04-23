# Image inputs

```go
msg := openrouter.CreateMultiModalMessage(
    "user",
    "What's in this image?",
    "https://example.com/image.jpg",
)
resp, _ := client.ChatComplete(ctx, []openrouter.Message{msg},
    openrouter.WithModel("openai/gpt-4o"),
)
```

For multiple images or base64 data URLs, build `[]ContentPart` directly — see [`examples/image-inputs/main.go`](../../examples/image-inputs/main.go).

The model must support vision. Check `ListModelEndpoints` or the OpenRouter models page.
