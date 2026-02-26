# SDK Improvement Opportunities

## High Impact

### 1. Per-request timeout/retry overrides
Currently timeouts and retry config are client-wide only. Adding `WithTimeout()` and `WithRetryConfig()` as request-level options would give callers more control over individual API calls.

### 2. Missing sampling options
`WithMinP()` and `WithTopA()` are absent despite the fields existing in the request structs. Easy wins for API completeness.

### 3. Review `panic()` calls
There are ~13 instances of `panic`/`log.Fatal` in non-test code. A library should never panic; these should return errors instead.

### 4. Break up `models.go`
At ~35K lines, it's unwieldy. Splitting into logical groupings (chat types, embedding types, provider types, etc.) would improve navigability.

### 5. Cost estimation helpers
Given we already have model pricing data from `ListModels()`, a utility like `EstimateCost(model, inputTokens, outputTokens)` would be very useful.

## Medium Impact

### 6. Configurable stream internals
Stream reconnection retries (hardcoded to 3), max backoff, and channel buffer size (hardcoded to 10) should be configurable.

### 7. Circuit breaker for streams
Repeated stream failures currently just retry with backoff. A circuit breaker pattern would prevent hammering a degraded endpoint.

### 8. Model-feature compatibility helpers
Something like `SupportsTools(modelID)` or `SupportsVision(modelID)` to help developers avoid trial-and-error.

### 9. Provider routing helpers
Functions to help pick providers by feature (e.g., "which providers support this model with ZDR?").

### 10. Better validation
Validate mutually exclusive parameters and model-specific constraints before sending requests, with actionable error messages.

## Lower Impact / Polish

### 11. Linting CI config
No `golangci-lint` config visible; adding one would catch issues automatically.

### 12. Performance benchmarks
No `Benchmark*` tests exist. Adding them for streaming, chunking, and embedding would help track regressions.

### 13. Network failure tests
Unit tests don't cover timeout/connection-reset scenarios well.

### 14. Migration guide
No versioned changelog or upgrade path documentation.

### 15. Deprecated field cleanup
Old fields like `IgnoreProviders` and `QuantizationFallback` are still present for backwards compatibility but could use a deprecation timeline.
