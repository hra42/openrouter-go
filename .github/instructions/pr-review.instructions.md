# Pull Request Review Instructions for openrouter-go

You are reviewing pull requests for **openrouter-go**, a zero-dependency Go client library for the OpenRouter API. Use these guidelines to provide thorough, constructive code reviews.

## Core Principles

1. **Zero Dependencies**: This library uses ONLY the Go standard library. REJECT any PR that adds external dependencies.
2. **Idiomatic Go**: Code must follow Go conventions and best practices.
3. **Thread Safety**: All public APIs must be safe for concurrent use.
4. **Comprehensive Testing**: All code changes require corresponding tests.
5. **Backward Compatibility**: Breaking changes to the public API require strong justification.

## Code Quality Standards

### Go Best Practices

- [ ] Code follows `gofmt` formatting standards
- [ ] Code passes `go vet` without warnings
- [ ] Variable and function names are clear and idiomatic (camelCase for private, PascalCase for public)
- [ ] Error handling is comprehensive - all errors are checked and properly handled
- [ ] Context is used appropriately for cancellation and timeouts
- [ ] No goroutine leaks - all goroutines have proper cleanup
- [ ] Mutexes and locks are used correctly to prevent race conditions
- [ ] defer statements are used appropriately for cleanup

### Project-Specific Patterns

- [ ] **Functional Options Pattern**: New configuration options use the functional options pattern
  - Options for client creation use `ClientOption` type
  - Options for requests use generic `RequestOption[T]` with `RequestConfig` constraint
  - All options are optional and have sensible defaults

- [ ] **Named Constants**: Magic numbers are extracted as named constants with descriptive names
  ```go
  const (
      defaultTimeout = 30 * time.Second
      maxRetries     = 3
  )
  ```

- [ ] **Generic Stream Implementation**: Streaming uses the generic `Stream[T]` type
  - ChatCompletions use `Stream[ChatCompletionChunk]`
  - Legacy completions use `Stream[CompletionChunk]`

- [ ] **Error Handling**: Use the custom `OpenRouterError` type for API errors
  - Preserve all error details from the API response
  - Include rate limit information when available

### Code Organization

- [ ] New types are defined in `models.go`
- [ ] Request/response handling is in appropriate endpoint files
- [ ] Functional options are in `options.go`
- [ ] Streaming logic uses `stream.go` infrastructure
- [ ] Internal utilities stay in `internal/` package

## Testing Requirements

### Unit Tests

- [ ] All new functions have corresponding unit tests
- [ ] Tests use table-driven test pattern where appropriate
- [ ] Tests use `httptest` to mock HTTP responses
- [ ] Edge cases and error conditions are tested
- [ ] Tests pass with `-race` flag (no race conditions)
- [ ] Test coverage for new code is substantial (aim for >80%)

Example test structure:
```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### E2E Tests (CRITICAL for New Endpoints)

When a new API endpoint is added, the PR **MUST** include:

- [ ] New test added to `cmd/openrouter-test/main.go`:
  - Test name added to `-test` flag description
  - Case handler in the switch statement
  - Implementation in appropriate test module file under `cmd/openrouter-test/tests/`
  - Test added to `runAllTests()` function

- [ ] E2E test validates:
  - Successful API calls with valid parameters
  - Proper error handling for invalid inputs
  - Response structure matches expected format
  - All response fields are populated correctly

- [ ] **Model parameter requirement**:
  - Tests that make actual API calls MUST accept `model string` parameter
  - NEVER hardcode model names in test functions (except `runModelEndpointsTest` and `runErrorTest`)
  - Pass model through from CLI flag to test function
  - For model suffixes, concatenate properly: `model+":nitro"`, `model+":online"`

## Documentation Requirements

- [ ] All exported functions, types, and methods have godoc comments
- [ ] Comments start with the name of the thing being documented
- [ ] Complex logic has inline comments explaining the "why"
- [ ] README.md is updated for new features
- [ ] CLAUDE.md is updated with relevant commands and architecture changes
- [ ] Examples are provided for significant new features in `examples/` directory

### Example Documentation

```go
// WithTimeout configures the HTTP client timeout for all requests.
// The default timeout is 30 seconds. Use context.Context for per-request timeouts.
func WithTimeout(timeout time.Duration) ClientOption {
    return func(c *Client) {
        c.httpClient.Timeout = timeout
    }
}
```

## API Design

### Public API Changes

- [ ] New public APIs follow existing naming conventions
- [ ] Breaking changes are clearly documented and justified
- [ ] New features work with both streaming and non-streaming modes
- [ ] Optional parameters use pointer types (`*string`, `*bool`, etc.)
- [ ] Required parameters use value types
- [ ] All public methods accept `context.Context` as first parameter

### Request/Response Types

- [ ] New request types implement `RequestConfig` interface if used with options
- [ ] JSON tags are lowercase with snake_case
- [ ] `omitempty` is used for optional fields
- [ ] Pointer types used for optional fields that need to distinguish nil from zero value
- [ ] Response types match OpenRouter API documentation

Example:
```go
type NewFeatureRequest struct {
    Model       string   `json:"model"`
    Required    string   `json:"required"`
    Optional    *string  `json:"optional,omitempty"`
    Defaults    int      `json:"defaults,omitempty"`
}
```

## Security Considerations

- [ ] No API keys or sensitive data in code or tests
- [ ] Input validation prevents injection attacks
- [ ] Error messages don't leak sensitive information
- [ ] Rate limiting and retry logic respect API guidelines
- [ ] HTTPS is enforced (no HTTP fallback)

## Performance

- [ ] No unnecessary allocations in hot paths
- [ ] Streaming responses don't buffer entire response in memory
- [ ] Goroutines are cleaned up properly
- [ ] Context cancellation is respected promptly
- [ ] HTTP connections are reused (using http.Client properly)

## Common Issues to Flag

### Critical Issues

- ❌ Adding external dependencies
- ❌ Race conditions (test with `go test -race`)
- ❌ Breaking public API without migration path
- ❌ Missing error handling
- ❌ Hardcoded credentials or API keys
- ❌ Missing tests for new functionality
- ❌ Goroutine leaks

### Should Fix

- ⚠️ Missing documentation for exported symbols
- ⚠️ Magic numbers without constants
- ⚠️ Duplicate code that could be refactored
- ⚠️ Inefficient algorithms or data structures
- ⚠️ Missing edge case testing
- ⚠️ Non-idiomatic Go code

### Nice to Have

- 💡 Additional examples for complex features
- 💡 More comprehensive error messages
- 💡 Performance optimizations
- 💡 Additional test coverage

## Review Checklist Summary

Before approving a PR, verify:

1. [ ] Zero dependencies maintained
2. [ ] All tests pass (`go test -race ./...`)
3. [ ] Code is formatted (`gofmt`, `go vet`)
4. [ ] New features have unit tests
5. [ ] New endpoints have E2E tests with model parameters
6. [ ] Documentation is complete and accurate
7. [ ] No breaking changes (or properly justified)
8. [ ] Thread-safe and context-aware
9. [ ] Follows project patterns (functional options, error handling)
10. [ ] Examples provided for significant features

## Helpful Commands for Reviewers

```bash
# Run all tests with race detection
go test -race ./...

# Check test coverage
go test -cover ./...

# Format code
go fmt ./...

# Run static analysis
go vet ./...

# Run E2E tests (requires API key)
export OPENROUTER_API_KEY="sk-..."
go run cmd/openrouter-test/main.go -test all

# Build without errors
go build ./...
```

## Tone and Communication

- Be constructive and educational in feedback
- Explain the "why" behind requested changes
- Recognize good work and clever solutions
- Suggest alternatives rather than just pointing out problems
- Link to relevant documentation or examples
- Prioritize critical issues over style preferences

Thank you for maintaining the quality and consistency of openrouter-go! 🚀
