# OpenRouter Go Client Documentation

A zero-dependency Go client library for the OpenRouter API.

OpenRouter is a unified API that provides access to multiple AI model providers through a single interface. This package implements complete bindings for the OpenRouter API, including support for chat completions, legacy completions, and full streaming capabilities.

## Installation

To install the package, use go get:

```bash
go get github.com/hra42/openrouter-go
```

## Basic Usage

Create a client and make a simple chat completion request:

```go
package main

import (
    "context"
    "fmt"
    "github.com/hra42/openrouter-go"
)

func main() {
    client := openrouter.NewClient("your-api-key")

    messages := []openrouter.Message{
        openrouter.CreateUserMessage("Hello, how are you?"),
    }

    response, err := client.ChatComplete(context.Background(), messages,
        openrouter.WithModel("openai/gpt-3.5-turbo"),
    )
    if err != nil {
        panic(err)
    }

    fmt.Println(response.Choices[0].Message.Content)
}
```

## Streaming

The library supports Server-Sent Events (SSE) streaming for real-time responses:

```go
stream, err := client.ChatCompleteStream(context.Background(), messages,
    openrouter.WithModel("openai/gpt-3.5-turbo"),
)
if err != nil {
    return err
}
defer stream.Close()

for event := range stream.Events() {
    // Process each streaming event
    fmt.Print(event.Choices[0].Delta.Content)
}

// Check for any errors that occurred during streaming
if err := stream.Err(); err != nil {
    return err
}
```

## Configuration

The client can be configured with various options:

```go
client := openrouter.NewClient("your-api-key",
    openrouter.WithBaseURL("https://custom.openrouter.ai"),
    openrouter.WithTimeout(60 * time.Second),
    openrouter.WithRetry(3, time.Second),
    openrouter.WithReferer("https://myapp.com"),
    openrouter.WithAppName("MyApplication"),
    openrouter.WithHeader("X-Custom", "value"),
)
```

### Available Client Options

| Option | Description |
|--------|-------------|
| `WithBaseURL(url)` | Set a custom API base URL |
| `WithHTTPClient(client)` | Use a custom HTTP client |
| `WithTimeout(duration)` | Set request timeout |
| `WithDefaultModel(model)` | Set default model for requests |
| `WithReferer(url)` | Set HTTP-Referer header |
| `WithAppName(name)` | Set application name (X-Title header) |
| `WithRetry(max, delay)` | Configure retry behavior |
| `WithHeader(key, value)` | Add custom headers |

## Error Handling

The library provides typed errors for different failure scenarios:

```go
resp, err := client.ChatComplete(ctx, messages, opts...)
if err != nil {
    if openrouter.IsRequestError(err) {
        reqErr := err.(*openrouter.RequestError)

        switch {
        case reqErr.IsRateLimitError():
            // Handle rate limiting - wait and retry
            time.Sleep(time.Second * 5)

        case reqErr.IsAuthenticationError():
            // Handle auth errors - check API key
            return errors.New("invalid API key")

        case reqErr.IsServerError():
            // Handle server errors - maybe retry
            return fmt.Errorf("server error: %s", reqErr.Message)
        }
    }
}
```

### Error Types

- `RequestError` - API request errors with status codes
- `ValidationError` - Input validation errors
- `StreamError` - Streaming-specific errors

### Error Checking Functions

- `IsRequestError(err)` - Check if error is a RequestError
- `IsValidationError(err)` - Check if error is a ValidationError
- `IsStreamError(err)` - Check if error is a StreamError

## Features

- ✅ Zero external dependencies
- ✅ Full API coverage (chat and legacy completions)
- ✅ SSE streaming support with automatic reconnection
- ✅ Comprehensive error handling with typed errors
- ✅ Exponential backoff retry logic
- ✅ Rate limiting support
- ✅ Context-aware cancellation
- ✅ Functional options pattern for flexible configuration
- ✅ Thread-safe operations

## API Methods

### Chat Completions

```go
// Non-streaming chat completion
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4"),
    openrouter.WithTemperature(0.7),
    openrouter.WithMaxTokens(1000),
)

// Streaming chat completion
stream, err := client.ChatCompleteStream(ctx, messages,
    openrouter.WithModel("openai/gpt-4"),
)
```

### Legacy Completions

```go
// Non-streaming completion
response, err := client.Complete(ctx, "Once upon a time",
    openrouter.WithCompletionModel("openai/gpt-3.5-turbo-instruct"),
    openrouter.WithCompletionMaxTokens(100),
)

// Streaming completion
stream, err := client.CompleteStream(ctx, "Once upon a time",
    openrouter.WithCompletionModel("openai/gpt-3.5-turbo-instruct"),
)
```

## Message Helpers

The library provides helper functions for creating messages:

```go
// System message
msg := openrouter.CreateSystemMessage("You are a helpful assistant")

// User message
msg := openrouter.CreateUserMessage("Hello!")

// Assistant message
msg := openrouter.CreateAssistantMessage("Hi there! How can I help?")

// Tool message
msg := openrouter.CreateToolMessage("Function result", "tool-call-id")

// Multi-modal message (text + image)
msg := openrouter.CreateMultiModalMessage(
    "user",
    "What's in this image?",
    "https://example.com/image.jpg",
)
```

## Request Options

### Chat Completion Options

```go
openrouter.WithModel(model)              // Model to use
openrouter.WithTemperature(0.7)          // Randomness (0-2)
openrouter.WithTopP(0.9)                 // Nucleus sampling
openrouter.WithTopK(40)                  // Top-K sampling
openrouter.WithMaxTokens(1000)           // Maximum response length
openrouter.WithStop("\\n", "END")        // Stop sequences
openrouter.WithFrequencyPenalty(0.5)     // Reduce repetition
openrouter.WithPresencePenalty(0.5)      // Encourage new topics
openrouter.WithSeed(42)                  // Reproducible outputs
openrouter.WithTools(tools...)           // Function calling
openrouter.WithResponseFormat(format)     // JSON mode
openrouter.WithLogProbs(5)               // Token probabilities
```

### Completion Options

```go
openrouter.WithCompletionModel(model)
openrouter.WithCompletionTemperature(0.7)
openrouter.WithCompletionMaxTokens(100)
openrouter.WithCompletionStop("\\n")
openrouter.WithCompletionLogProbs(5)
openrouter.WithCompletionEcho(true)
openrouter.WithCompletionN(3)
openrouter.WithCompletionBestOf(5)
```

## Rate Limiting

The library includes built-in rate limiting support:

```go
// Create a rate limiter (10 requests/sec, burst of 5)
limiter := openrouter.NewRateLimiter(10.0, 5)
defer limiter.Close()

// Wait before making request
err := limiter.Wait(context.Background())
if err != nil {
    return err
}

response, err := client.ChatComplete(ctx, messages, opts...)
```

## Advanced Features

### Function/Tool Calling

```go
tools := []openrouter.Tool{
    {
        Type: "function",
        Function: openrouter.Function{
            Name:        "get_weather",
            Description: "Get current weather",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "location": map[string]interface{}{
                        "type": "string",
                        "description": "City and state",
                    },
                },
                "required": []string{"location"},
            },
        },
    },
}

response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4"),
    openrouter.WithTools(tools...),
    openrouter.WithToolChoice("auto"),
)

// Check for tool calls in response
if len(response.Choices[0].Message.ToolCalls) > 0 {
    toolCall := response.Choices[0].Message.ToolCalls[0]
    fmt.Printf("Function: %s\n", toolCall.Function.Name)
    fmt.Printf("Arguments: %s\n", toolCall.Function.Arguments)
}
```

### JSON Response Format

```go
format := openrouter.ResponseFormat{
    Type: "json_object",
    JSONSchema: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "answer": map[string]interface{}{"type": "string"},
            "confidence": map[string]interface{}{"type": "number"},
        },
    },
}

response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4"),
    openrouter.WithResponseFormat(format),
)
```

### Provider Preferences

```go
provider := openrouter.Provider{
    Order:            []string{"OpenAI", "Anthropic"},
    RequireParameters: true,
    DataCollection:   "deny",
    AllowFallbacks:   true,
}

response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4"),
    openrouter.WithProvider(provider),
)
```

## Context and Cancellation

All API methods support context for cancellation and timeouts:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

response, err := client.ChatComplete(ctx, messages, opts...)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Cancel after some condition
go func() {
    time.Sleep(5 * time.Second)
    cancel()
}()

stream, err := client.ChatCompleteStream(ctx, messages, opts...)
```

## Environment Variables

The library respects the following environment variables:

- `OPENROUTER_API_KEY` - Your OpenRouter API key (can be passed to NewClient instead)

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...

# Run specific test
go test -run TestChatComplete
```

## Examples

The `examples/` directory contains comprehensive examples:

- **basic/** - Simple usage examples for common tasks
- **streaming/** - Real-time streaming response handling
- **advanced/** - Advanced features like function calling, rate limiting, and custom configuration

## Live API Testing

Test the library against the live API:

```bash
# Install the test tool
go install github.com/hra42/openrouter-go/cmd/openrouter-test@latest

# Run tests
export OPENROUTER_API_KEY="your-key"
openrouter-test -test all
```

## Best Practices

1. **Always use context** - Pass context to support cancellation and timeouts
2. **Handle errors properly** - Check error types for appropriate handling
3. **Set reasonable timeouts** - Use WithTimeout to prevent hanging requests
4. **Implement retries** - Use WithRetry for resilience against temporary failures
5. **Monitor rate limits** - Use the rate limiter for high-volume applications
6. **Close streams** - Always defer stream.Close() for streaming responses
7. **Validate inputs** - The library validates required fields but check your data

## Troubleshooting

### Common Issues

**Rate Limiting**
```go
if reqErr, ok := err.(*openrouter.RequestError); ok && reqErr.IsRateLimitError() {
    // Wait and retry
    time.Sleep(time.Second * 5)
}
```

**Timeout Errors**
```go
client := openrouter.NewClient(apiKey,
    openrouter.WithTimeout(60 * time.Second), // Increase timeout
)
```

**Streaming Issues**
```go
// Always check stream errors
if err := stream.Err(); err != nil {
    log.Printf("Stream error: %v", err)
}
```

## Support

For issues, questions, or contributions, visit: https://github.com/hra42/openrouter-go