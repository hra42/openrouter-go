# Tool Calling Examples

This directory contains examples demonstrating how to use tool/function calling with the OpenRouter Go client.

## Examples

### main.go - Complete Tool Calling Examples
Demonstrates various tool calling scenarios:
- **Basic tool calling** - Search for books using a simulated Gutenberg library API
- **Tool parameters** - Weather tool with enum parameters and optional fields
- **Multiple tools** - Model selects appropriate tool from multiple available options
- **Forced tool choice** - Force the model to use a specific tool
- **Parallel tool calls** - Control whether tools can be called in parallel

### streaming.go - Streaming with Tool Calls
Shows how to properly handle tool calls in streaming mode:
- **Accumulating tool call arguments** - Handles arguments that arrive in chunks
- **Processing tool results** - Execute tools and send results back to the model
- **Final response streaming** - Stream the model's final response after tool execution

## Running the Examples

1. Set your OpenRouter API key:
```bash
export OPENROUTER_API_KEY="your-api-key"
```

2. Run the examples:
```bash
# Non-streaming tool calling examples
go run main.go

# Streaming with tool calls
go run streaming.go
```

## Key Concepts

### Tool Definition
Tools are defined with a name, description, and JSON Schema parameters:

```go
tools := []openrouter.Tool{
    {
        Type: "function",
        Function: openrouter.Function{
            Name:        "get_weather",
            Description: "Get current weather for a location",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "location": map[string]interface{}{
                        "type":        "string",
                        "description": "City name or zip code",
                    },
                },
                "required": []string{"location"},
            },
        },
    },
}
```

### Tool Execution Flow
1. **Request with tools** - Include tools in the chat completion request
2. **Check for tool calls** - Model may request one or more tool calls
3. **Execute tools** - Run the requested functions with provided arguments
4. **Return results** - Send tool results back as "tool" role messages
5. **Get final response** - Model provides final answer using tool results

### Streaming Tool Calls - Important!
When streaming, tool call information arrives in fragments:
- **First fragment** contains the tool call ID and function name
- **Subsequent fragments** contain argument chunks without IDs
- **Accumulation required** - Track the current tool call ID to properly append argument fragments

```go
// Correct pattern for accumulating streaming tool calls
currentToolCallID := ""
toolCallsMap := make(map[string]*openrouter.ToolCall)

for event := range stream.Events() {
    for _, choice := range event.Choices {
        if choice.Delta != nil && len(choice.Delta.ToolCalls) > 0 {
            for _, tc := range choice.Delta.ToolCalls {
                if tc.ID != "" {
                    // New tool call or update
                    currentToolCallID = tc.ID
                    // Create or update tool call in map
                } else if currentToolCallID != "" {
                    // Append arguments to current tool call
                    toolCallsMap[currentToolCallID].Function.Arguments += tc.Function.Arguments
                }
            }
        }
    }
}
```

## Supported Models

Tool calling is supported by many models including:
- OpenAI GPT-4o, GPT-4o-mini, GPT-3.5-turbo
- Anthropic Claude 3.5 Sonnet, Claude 3 Opus
- Google Gemini models
- Many open-source models

Check model compatibility at [openrouter.ai/models](https://openrouter.ai/models?supported_parameters=tools)

## Best Practices

1. **Validate arguments** - Always validate tool arguments before execution
2. **Handle errors gracefully** - Return error messages in tool results
3. **Preserve context** - Include all messages (including tool calls and results) in subsequent requests
4. **Test streaming** - Streaming tool calls require special handling for argument accumulation
5. **Use structured responses** - Return JSON for consistent tool result parsing