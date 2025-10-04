# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build and Test
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test -run TestChatComplete

# Run tests in verbose mode
go test -v ./...

# Build the package
go build ./...

# Format code
go fmt ./...

# Run go vet for static analysis
go vet ./...

# Install dependencies
go mod download

# Update dependencies
go mod tidy
```

### Running Examples
```bash
# Set API key
export OPENROUTER_API_KEY="your-api-key"

# Run examples
go run examples/basic/main.go
go run examples/streaming/main.go
go run examples/list-models/main.go
go run examples/advanced/main.go
go run examples/structured-output/main.go
go run examples/tool-calling/main.go
go run examples/web_search/main.go
go run examples/transforms/main.go
go run examples/app-attribution/main.go
```

### Running E2E Tests
```bash
# Set API key
export OPENROUTER_API_KEY="your-api-key"

# Run all tests
go run cmd/openrouter-test/main.go -test all

# Run specific test
go run cmd/openrouter-test/main.go -test models
go run cmd/openrouter-test/main.go -test chat
go run cmd/openrouter-test/main.go -test streaming

# Run with verbose output
go run cmd/openrouter-test/main.go -test models -v
```

## Architecture Overview

This is a zero-dependency Go client library for the OpenRouter API that follows idiomatic Go patterns.

### Core Components

**Client (`client.go`)**: The main client struct that manages HTTP communication, authentication, and request routing. Uses functional options pattern for configuration.

**Request/Response Types (`models.go`)**: Defines all data structures for API requests and responses, including messages, tools, and completion parameters.

**Functional Options (`options.go`)**: Implements the options pattern for flexible configuration of both client and request parameters. This allows clean API usage like `WithModel()`, `WithMessages()`, etc.

**Streaming (`stream.go`)**: Implements Server-Sent Events (SSE) parsing for streaming responses without external dependencies. Handles reconnection and error recovery.

**Error Handling (`errors.go`)**: Custom error types that preserve OpenRouter API error details including rate limits, retry information, and moderation flags.

**Retry Logic (`retry.go`)**: Implements exponential backoff with jitter for handling transient failures and rate limits.

### API Endpoints

**Chat Completions (`chat.go`)**: Modern chat API supporting messages, tools, structured outputs, and streaming.

**Legacy Completions (`completions.go`)**: Support for older prompt-based completion API.

**Models Listing (`models_endpoint.go`)**: Retrieve available models with filtering by category and detailed model information.

**Web Search (`web_search.go`)**: Integration with OpenRouter's web search plugin for augmented responses.

### Key Design Patterns

1. **Functional Options**: Both client creation and API calls use functional options for clean, extensible configuration
2. **Context Support**: All API methods accept context.Context for proper cancellation and timeout handling
3. **Zero Dependencies**: Entire library uses only Go standard library
4. **Thread Safety**: Client is safe for concurrent use across goroutines
5. **Streaming**: Custom SSE parser handles streaming responses with proper cleanup

### Testing Approach

- Unit tests for all public APIs using httptest for mocking HTTP responses
- Table-driven tests for comprehensive coverage of different scenarios
- Race detection tests to ensure thread safety
- Integration test utility in `cmd/openrouter-test/` for live API testing

**IMPORTANT**: When adding a new endpoint, always add a corresponding e2e test to `cmd/openrouter-test/main.go`:
1. Add the test name to the `-test` flag description
2. Add a case handler in the switch statement
3. Implement a `runXxxTest()` function with comprehensive validation
4. Add the test to the `runAllTests()` function array

**CRITICAL E2E Test Guidelines**:
- **ALWAYS use the `model` parameter** for ANY test that makes actual API calls to chat/completion endpoints
- **NEVER hardcode model names** in test functions (e.g., "openai/gpt-4o", "meta-llama/llama-3.1-8b-instruct")
- Test function signatures that make API calls MUST accept `model string` as a parameter
- Pass the `model` parameter through the switch statement and `runAllTests()` function
- The ONLY exceptions are:
  - `runModelEndpointsTest` - uses hardcoded models for metadata inspection (not API calls)
  - `runErrorTest` - uses "invalid/nonexistent-model-xyz" for error testing
  - `runModelsTest` - doesn't make completion calls
- When using model suffixes, concatenate properly: `model+":nitro"`, `model+":floor"`, `model+":online"`

### Important Notes

- Always check error returns from API calls
- Streaming responses require proper cleanup by reading all events or canceling context
- The library automatically handles rate limiting with configurable retry logic
- Tool calls and structured outputs require compatible models (check OpenRouter docs)
- Web search can use either native provider search or Exa, with different pricing models