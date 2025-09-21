# OpenRouter Go Client - Live API Test Tool

A command-line tool to test the openrouter-go library against the live OpenRouter API.

## Installation

```bash
go install github.com/hra42/openrouter-go/cmd/openrouter-test@latest
```

Or build from source:

```bash
go build -o openrouter-test ./cmd/openrouter-test
```

## Usage

### Prerequisites

You need an OpenRouter API key. Get one at https://openrouter.ai/

### Basic Usage

```bash
# Set API key via environment variable
export OPENROUTER_API_KEY="your-api-key"
openrouter-test

# Or pass API key directly
openrouter-test -key "your-api-key"
```

### Available Tests

Run all tests:
```bash
openrouter-test -test all
```

Run specific tests:
```bash
# Test chat completion
openrouter-test -test chat

# Test streaming
openrouter-test -test stream

# Test legacy completion
openrouter-test -test completion

# Test error handling
openrouter-test -test error

# Test provider routing
openrouter-test -test provider

# Test Zero Data Retention (ZDR)
openrouter-test -test zdr

# Test model suffixes (nitro/floor)
openrouter-test -test suffix

# Test price constraints
openrouter-test -test price

# Test structured output
openrouter-test -test structured

# Test tool/function calling
openrouter-test -test tools

# Test message transforms
openrouter-test -test transforms

# Test web search
openrouter-test -test websearch
```

### Options

```
-key string
    OpenRouter API key (or set OPENROUTER_API_KEY env var)

-model string
    Model to use (default: "openai/gpt-3.5-turbo")

-test string
    Test to run: all, chat, stream, completion, error, provider, zdr, suffix,
    price, structured, tools, transforms, websearch (default: "all")

-max-tokens int
    Maximum tokens for response (default: 100)

-timeout duration
    Request timeout (default: 30s)

-v
    Verbose output
```

### Examples

Test with a specific model:
```bash
openrouter-test -model "anthropic/claude-3-haiku" -test chat
```

Test streaming with verbose output:
```bash
openrouter-test -test stream -v
```

Quick test with minimal tokens:
```bash
openrouter-test -test chat -max-tokens 20
```

Test with custom timeout:
```bash
openrouter-test -timeout 60s -test all
```

## Test Descriptions

### Chat Completion Test
Tests the standard chat completion endpoint with a simple math question.

### Streaming Test
Tests SSE streaming by asking the model to count from 1 to 5.

### Legacy Completion Test
Tests the legacy completion endpoint (requires instruct model support).

### Error Handling Test
Deliberately triggers an error to test error handling capabilities.

### Provider Routing Test
Tests provider order, fallbacks, and routing constraints.

### ZDR Test
Tests Zero Data Retention enforcement for privacy-sensitive requests.

### Model Suffix Test
Tests :nitro (throughput) and :floor (lowest price) model suffixes.

### Price Constraint Test
Tests maximum price constraints and data collection policies.

### Structured Output Test
Tests JSON schema validation and structured response formats.

### Tool/Function Calling Test
Tests function calling capabilities including parallel calls and streaming.

### Message Transforms Test
Tests middle-out compression for managing context window limits.

### Web Search Test
Tests web search integration using :online suffix, web plugins, and different search engines (Native/Exa).

## Exit Codes

- `0`: All tests passed
- `1`: One or more tests failed or error occurred

## Output

The tool provides colored output with:
- 🔄 Test in progress
- ✅ Test passed
- ❌ Test failed
- ⚠️ Test skipped (e.g., model not available)
- 📊 Summary statistics

## Troubleshooting

### "Model not found" errors
Some models require specific permissions or paid access. Try using `openai/gpt-3.5-turbo` which is widely available.

### Rate limiting
If you encounter rate limits, the tool will automatically retry with exponential backoff.

### Timeout errors
Increase the timeout with `-timeout 60s` for slower models or connections.