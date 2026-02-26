# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Per-request timeout and retry overrides via `WithRequestTimeout()` and `WithRequestRetryConfig()`
- Missing sampling options: `WithMinP()`, `WithTopA()`
- Configurable stream internals: reconnection retries, max backoff, channel buffer size
- Circuit breaker pattern for stream failure protection
- Model-feature compatibility helpers: `SupportsTools()`, `SupportsVision()`, etc.
- Provider routing helpers for feature-based provider selection
- Request validation with actionable error messages
- Cost estimation utility via `EstimateCost()`
- Performance benchmarks for SSE parser, chunking, and embeddings
- Network failure and retry tests for improved coverage
- CI pipeline with Go 1.22/1.23/1.24 matrix and golangci-lint
- This CHANGELOG

### Removed
- Deprecated `Provider.IgnoreProviders` field (use `Provider.Ignore` with `WithIgnoreProviders()`)
- Deprecated `Provider.QuantizationFallback` field (use `Provider.Quantizations` with `WithQuantizations()`)

## [0.1.0] - Initial Release

### Features
- **Client**: Configurable HTTP client with functional options pattern, custom headers, and API key authentication
- **Chat Completions**: Full chat API support including messages, tools, structured outputs, and function calling
- **Streaming**: Zero-dependency SSE parser with reconnection and error recovery
- **Embeddings**: Single and batch embedding creation with provider options and model listing
- **Chunking**: Text chunking by characters, words, sentences, paragraphs, sections, and tokens with overlap support
- **Models API**: List and filter available models with detailed model information
- **Responses API** (Beta): OpenAI-compatible Responses API with reasoning, tool calling, web search, and streaming
- **Broadcast Webhook**: OTLP JSON trace payload parsing for OpenRouter's Broadcast feature
- **Web Search**: Integration with OpenRouter's web search plugin
- **MCP Tool Conversion**: Utilities for converting MCP tool definitions to OpenAI-compatible format
- **Retry Logic**: Exponential backoff with jitter, configurable retry policies, and rate limiting
- **Error Handling**: Custom error types preserving API error details, rate limits, and moderation flags
- **Zero Dependencies**: Entire library uses only Go standard library
- **Thread Safety**: Client safe for concurrent use across goroutines
