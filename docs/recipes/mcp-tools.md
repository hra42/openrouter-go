# MCP tool conversion

Convert [Model Context Protocol](https://modelcontextprotocol.io) tool definitions into the OpenAI-compatible shape the chat API expects. No `Client` required — these are pure utilities in `mcp.go`.

```go
mcpTool := openrouter.MCPTool{
    Name:        "get_weather",
    Description: "Get current weather",
    InputSchema: openrouter.MCPInputSchema{
        Type: "object",
        Properties: map[string]any{
            "location": map[string]any{"type": "string"},
        },
        Required: []string{"location"},
    },
}

tool := openrouter.ConvertMCPTool(mcpTool)

resp, _ := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithTools(tool),
)
```

See [`examples/mcp-tools/main.go`](../../examples/mcp-tools/main.go).
