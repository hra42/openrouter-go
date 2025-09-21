# openrouter-go

A zero-dependency Go package providing complete bindings for the OpenRouter API, supporting all available endpoints with full streaming capabilities.

## Features

- ✅ Complete API coverage (chat completions and legacy completions)
- ✅ Full streaming support with Server-Sent Events (SSE)
- ✅ Zero external dependencies
- ✅ Go 1.25.1 support
- ✅ Comprehensive error handling and retry logic
- ✅ Context-aware cancellation
- ✅ Thread-safe client operations
- ✅ Extensive configuration options via functional options pattern
- ✅ Per-request Zero Data Retention (ZDR) enforcement

## Installation

```bash
go get github.com/hra42/openrouter-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/hra42/openrouter-go"
)

func main() {
    client := openrouter.NewClient("your-api-key")

    response, err := client.ChatComplete(context.Background(), []openrouter.Message{
        {Role: "user", Content: "Hello, how are you?"},
    }, openrouter.WithModel("openai/gpt-4"))

    if err != nil {
        panic(err)
    }

    fmt.Println(response.Choices[0].Message.Content)
}
```

## Development Roadmap

### Phase 1: Foundation ✅
**Status:** Complete
- ✅ Set up Go module with appropriate package structure
- ✅ Define all request/response structs based on API documentation
- ✅ Create base client struct with configuration options
- ✅ Implement error types and custom error handling
- ✅ Design functional options pattern for optional parameters

### Phase 2: HTTP Communication Layer ✅
**Status:** Complete
- ✅ Implement HTTP request construction and execution
- ✅ Add authentication header management
- ✅ Create retry logic with exponential backoff
- ✅ Handle rate limiting and API errors gracefully
- ✅ Support custom HTTP client injection for testing

### Phase 3: Core API Implementation ✅
**Status:** Complete
- ✅ Implement ChatComplete method for chat completions endpoint
- ✅ Implement Complete method for legacy completions endpoint
- ✅ Add all optional parameters support (temperature, top_p, etc.)
- ✅ Ensure proper request validation
- ✅ Handle both streaming and non-streaming responses

### Phase 4: Streaming Support ✅
**Status:** Complete
- ✅ Build SSE parser from scratch (no external dependencies)
- ✅ Create streaming response types with channel-based API
- ✅ Implement proper connection management and cleanup
- ✅ Add context cancellation support for streams
- ✅ Handle streaming errors and reconnection logic

### Phase 5: Production Readiness ✅
**Status:** Complete
- ✅ Write comprehensive unit tests with mocked responses
- ✅ Add integration tests (with build tags)
- ✅ Create detailed usage examples for common scenarios
- ✅ Benchmark performance-critical paths
- ✅ Complete API documentation with godoc comments

## API Design

### Client Initialization

```go
// Basic initialization
client := openrouter.NewClient("api-key")

// With options
client := openrouter.NewClient("api-key",
    openrouter.WithBaseURL("https://custom.openrouter.ai"),
    openrouter.WithHTTPClient(customHTTPClient),
    openrouter.WithTimeout(60 * time.Second),
    openrouter.WithRetry(3, time.Second),
    openrouter.WithAppName("MyApp"),
    openrouter.WithReferer("https://myapp.com"),
)
```

### Chat Completions

```go
// Non-streaming
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("anthropic/claude-3-opus"),
    openrouter.WithTemperature(0.7),
    openrouter.WithMaxTokens(1000),
)

// Streaming
stream, err := client.ChatCompleteStream(ctx, messages,
    openrouter.WithModel("anthropic/claude-3-opus"),
)

for event := range stream.Events() {
    fmt.Print(event.Choices[0].Delta.Content)
}

if err := stream.Err(); err != nil {
    // Handle streaming error
}

// With Zero Data Retention (ZDR)
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("anthropic/claude-3-opus"),
    openrouter.WithZDR(true), // Enforce ZDR for this request
)
```

### Legacy Completions

```go
// Non-streaming
response, err := client.Complete(ctx, "Once upon a time",
    openrouter.WithModel("openai/gpt-3.5-turbo-instruct"),
    openrouter.WithMaxTokens(100),
)

// Streaming
stream, err := client.CompleteStream(ctx, "Once upon a time",
    openrouter.WithModel("openai/gpt-3.5-turbo-instruct"),
)

// With Zero Data Retention (ZDR)
response, err := client.Complete(ctx, "Once upon a time",
    openrouter.WithModel("openai/gpt-3.5-turbo-instruct"),
    openrouter.WithCompletionZDR(true), // Enforce ZDR for this request
)
```

## Package Structure

```
openrouter-go/
├── client.go           # Main client implementation
├── completions.go      # Completion endpoint methods
├── chat.go            # Chat completion endpoint methods
├── models.go          # Request/response type definitions
├── options.go         # Functional options for configuration
├── stream.go          # SSE streaming implementation
├── errors.go          # Custom error types
├── retry.go           # Retry and backoff logic
├── examples/
│   ├── basic/         # Basic usage examples
│   ├── streaming/     # Streaming examples
│   └── advanced/      # Advanced configuration examples
└── internal/
    └── sse/           # Internal SSE parser implementation
```

## Requirements

- Go 1.25.1
- No external dependencies

## Documentation

For detailed API documentation and usage examples, see [DOCUMENTATION.md](DOCUMENTATION.md).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or distribute this software, either in source code form or as a compiled binary, for any purpose, commercial or non-commercial, and by any means.

In jurisdictions that recognize copyright laws, the author or authors of this software dedicate any and all copyright interest in the software to the public domain. We make this dedication for the benefit of the public at large and to the detriment of our heirs and successors. We intend this dedication to be an overt act of relinquishment in perpetuity of all present and future rights to this software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <https://unlicense.org>

## Status

✅ **Production Ready** - All 5 phases complete! The library is now ready for production use with:
- ✅ Full foundation with all types and error handling
- ✅ Robust HTTP communication with retry logic
- ✅ Complete API implementation for chat and completions
- ✅ Zero-dependency SSE streaming with reconnection support
- ✅ Comprehensive test coverage and documentation
- ✅ Production-ready examples for all use cases

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test -run TestChatComplete
```

### Provider Routing

The library supports comprehensive provider routing options to control how your requests are handled across different providers.

#### Basic Provider Routing

```go
// Specify provider order
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("meta-llama/llama-3.1-70b-instruct"),
    openrouter.WithProviderOrder("together", "openai", "anthropic"),
)

// Disable fallbacks (only use specified providers)
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("mistralai/mixtral-8x7b-instruct"),
    openrouter.WithProviderOrder("together"),
    openrouter.WithAllowFallbacks(false),
)

// Sort providers by throughput or price
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("meta-llama/llama-3.1-70b-instruct"),
    openrouter.WithProviderSort("throughput"), // or "price", "latency"
)
```

#### Model Suffixes

```go
// Use :nitro suffix for throughput optimization
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("meta-llama/llama-3.1-70b-instruct:nitro"),
)

// Use :floor suffix for lowest price
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("meta-llama/llama-3.1-70b-instruct:floor"),
)
```

#### Provider Filtering

```go
// Only use specific providers
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithOnlyProviders("azure", "openai"),
)

// Ignore specific providers
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("meta-llama/llama-3.3-70b-instruct"),
    openrouter.WithIgnoreProviders("deepinfra"),
)

// Filter by quantization levels
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("meta-llama/llama-3.1-8b-instruct"),
    openrouter.WithQuantizations("fp8", "fp16"),
)
```

#### Price Constraints

```go
// Set maximum pricing constraints
maxPrice := openrouter.MaxPrice{
    Prompt: 1.0,     // Max $1 per million prompt tokens
    Completion: 2.0, // Max $2 per million completion tokens
}
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("meta-llama/llama-3.1-70b-instruct"),
    openrouter.WithMaxPrice(maxPrice),
    openrouter.WithProviderSort("throughput"), // Use fastest provider under price limit
)
```

#### Data Policies

```go
// Require providers that don't collect data
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("anthropic/claude-3-opus"),
    openrouter.WithDataCollection("deny"), // or "allow"
)

// Require providers that support all parameters
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithRequireParameters(true),
    openrouter.WithResponseFormat(openrouter.ResponseFormat{Type: "json_object"}),
)
```

### Zero Data Retention (ZDR)

The library supports per-request Zero Data Retention enforcement. When enabled, requests will only be routed to endpoints with Zero Data Retention policies.

```go
// For chat completions
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("anthropic/claude-3-opus"),
    openrouter.WithZDR(true), // Enforce ZDR for this specific request
)

// For legacy completions
response, err := client.Complete(ctx, prompt,
    openrouter.WithModel("openai/gpt-3.5-turbo-instruct"),
    openrouter.WithCompletionZDR(true), // Enforce ZDR for this specific request
)

// With custom provider configuration
provider := openrouter.Provider{
    ZDR: &[]bool{true}[0], // Enable ZDR
}
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("anthropic/claude-3-opus"),
    openrouter.WithProvider(provider),
)
```

Note: The request-level `zdr` parameter operates as an "OR" with your account-wide ZDR setting. If either is enabled, ZDR enforcement will be applied.

## Examples

The `examples/` directory contains comprehensive examples:

- **basic/** - Simple usage examples for common tasks
- **streaming/** - Real-time streaming response handling
- **advanced/** - Advanced features like function calling, rate limiting, and custom configuration

To run an example:

```bash
# Set your API key
export OPENROUTER_API_KEY="your-api-key"

# Run basic examples
go run examples/basic/main.go

# Run streaming examples
go run examples/streaming/main.go

# Run advanced examples
go run examples/advanced/main.go
```
