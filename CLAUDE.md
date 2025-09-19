# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the openrouter-go library - a zero-dependency Go package providing complete bindings for the OpenRouter API. The project is currently in initial development phase with no implementation code yet written.

## Development Phases

The implementation follows a structured roadmap:
1. **Foundation**: Core types and project structure
2. **HTTP Communication**: HTTP client with retry logic and error handling
3. **Core API**: ChatComplete and Complete endpoints
4. **Streaming**: SSE streaming implementation from scratch
5. **Production**: Testing, benchmarks, and documentation

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

## Key Design Decisions

- **Zero dependencies**: All functionality including SSE streaming must be implemented without external packages
- **Functional options pattern**: Use for all optional parameters (WithModel, WithTemperature, etc.)
- **Channel-based streaming API**: Streaming responses use channels for Events() and Err()
- **Context-aware**: All API methods accept context for cancellation
- **Custom HTTP client injection**: Support for testing via WithHTTPClient option

## API Endpoints to Implement

1. **Chat Completions**: `/api/v1/chat/completions` (ChatComplete and ChatCompleteStream)
2. **Legacy Completions**: `/api/v1/completions` (Complete and CompleteStream)

## Testing Commands

Since the project is in initial development, testing commands will be:
```bash
go test ./...           # Run all tests
go test -v ./...        # Run tests with verbose output
go test -race ./...     # Run tests with race detection
go test -cover ./...    # Run tests with coverage
```

## Development Requirements

- Go 1.25.1
- No external dependencies allowed
- Follow standard Go project layout
- Use godoc-style comments for all exported types and methods