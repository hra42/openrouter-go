# Structured output (JSON schema)

Force the model to return JSON matching a schema:

```go
format := openrouter.ResponseFormat{
    Type: "json_schema",
    JSONSchema: map[string]any{
        "name":   "weather_report",
        "strict": true,
        "schema": map[string]any{
            "type": "object",
            "properties": map[string]any{
                "city":        map[string]any{"type": "string"},
                "temperature": map[string]any{"type": "number"},
                "conditions":  map[string]any{"type": "string"},
            },
            "required":             []string{"city", "temperature", "conditions"},
            "additionalProperties": false,
        },
    },
}

resp, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithResponseFormat(format),
)
```

`resp.Choices[0].Message.Content` will be a JSON string; `json.Unmarshal` it into your Go type.

For plain JSON mode without a schema, use `Type: "json_object"`.

See [`examples/structured-output/main.go`](../../examples/structured-output/main.go).
