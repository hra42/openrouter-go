# Tool / function calling

The model must support tools (e.g. `openai/gpt-4o`, `anthropic/claude-3.5-sonnet`). Non-tool-capable models return an error.

```go
tools := []openrouter.Tool{{
    Type: "function",
    Function: openrouter.Function{
        Name:        "get_weather",
        Description: "Get current weather for a city.",
        Parameters: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "location": map[string]any{"type": "string", "description": "City and state"},
            },
            "required": []string{"location"},
        },
    },
}}

resp, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithTools(tools...),
    openrouter.WithToolChoice("auto"),
)

for _, tc := range resp.Choices[0].Message.ToolCalls {
    fmt.Printf("call %s: %s(%s)\n", tc.ID, tc.Function.Name, tc.Function.Arguments)
}
```

Feed the tool result back with `CreateToolMessage`:

```go
messages = append(messages, resp.Choices[0].Message)
messages = append(messages, openrouter.CreateToolMessage(`{"temp":72}`, toolCallID))
resp2, _ := client.ChatComplete(ctx, messages, openrouter.WithModel("openai/gpt-4o"), openrouter.WithTools(tools...))
```

## Streaming tool calls

Tool calls arrive as deltas — accumulate by `Index`:

```go
type pending struct{ id, name, args string }
calls := map[int]*pending{}

for event := range stream.Events() {
    for _, choice := range event.Choices {
        if choice.Delta == nil { continue }
        for _, d := range choice.Delta.ToolCalls {
            p, ok := calls[d.Index]
            if !ok { p = &pending{}; calls[d.Index] = p }
            if d.ID != "" { p.id = d.ID }
            if d.Function.Name != "" { p.name = d.Function.Name }
            p.args += d.Function.Arguments
        }
    }
}
```

See [`examples/tool-calling/main.go`](../../examples/tool-calling/main.go).
