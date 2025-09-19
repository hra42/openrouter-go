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

### Phase 1: Foundation
**Scope:** Establish project structure and core types
- Set up Go module with appropriate package structure
- Define all request/response structs based on API documentation
- Create base client struct with configuration options
- Implement error types and custom error handling
- Design functional options pattern for optional parameters

### Phase 2: HTTP Communication Layer
**Scope:** Build robust HTTP client functionality
- Implement HTTP request construction and execution
- Add authentication header management
- Create retry logic with exponential backoff
- Handle rate limiting and API errors gracefully
- Support custom HTTP client injection for testing

### Phase 3: Core API Implementation
**Scope:** Implement both API endpoints
- Implement ChatComplete method for chat completions endpoint
- Implement Complete method for legacy completions endpoint
- Add all optional parameters support (temperature, top_p, etc.)
- Ensure proper request validation
- Handle both streaming and non-streaming responses

### Phase 4: Streaming Support
**Scope:** Full SSE streaming implementation
- Build SSE parser from scratch (no external dependencies)
- Create streaming response types with channel-based API
- Implement proper connection management and cleanup
- Add context cancellation support for streams
- Handle streaming errors and reconnection logic

### Phase 5: Production Readiness
**Scope:** Testing, documentation, and polish
- Write comprehensive unit tests with mocked responses
- Add integration tests (with build tags)
- Create detailed usage examples for common scenarios
- Benchmark performance-critical paths
- Complete API documentation with godoc comments

## API Design

### Client Initialization

```go
// Basic initialization
client := openrouter.NewClient("api-key")

// With options
client := openrouter.NewClient("api-key",
    openrouter.WithBaseURL("https://custom.openrouter.ai"),
    openrouter.WithHTTPClient(customHTTPClient),
    openrouter.WithRetryCount(3),
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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or distribute this software, either in source code form or as a compiled binary, for any purpose, commercial or non-commercial, and by any means.

In jurisdictions that recognize copyright laws, the author or authors of this software dedicate any and all copyright interest in the software to the public domain. We make this dedication for the benefit of the public at large and to the detriment of our heirs and successors. We intend this dedication to be an overt act of relinquishment in perpetuity of all present and future rights to this software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <https://unlicense.org>

## Status

🚧 **In Development** - This package is currently being built following the roadmap above.
