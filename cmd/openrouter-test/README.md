# OpenRouter Go Client - E2E Test Suite

This is a comprehensive end-to-end test suite for the OpenRouter Go client library. It tests all API endpoints with live API calls.

## Setup

Set your OpenRouter API key:

```bash
export OPENROUTER_API_KEY="your-api-key"
```

## Running Tests

### Run All Tests

```bash
go run cmd/openrouter-test/main.go -test all
```

### Run Specific Tests

```bash
# Chat completion
go run cmd/openrouter-test/main.go -test chat

# Streaming
go run cmd/openrouter-test/main.go -test stream

# Image tests
go run cmd/openrouter-test/main.go -test image          # URL-based image
go run cmd/openrouter-test/main.go -test multiimage     # Multiple images
go run cmd/openrouter-test/main.go -test imagedetail    # Image with detail level
go run cmd/openrouter-test/main.go -test contentbuilder # ContentBuilder API
go run cmd/openrouter-test/main.go -test base64image    # Local image (base64)

# Other tests
go run cmd/openrouter-test/main.go -test structured     # Structured outputs
go run cmd/openrouter-test/main.go -test tools          # Tool calling
go run cmd/openrouter-test/main.go -test websearch      # Web search
go run cmd/openrouter-test/main.go -test models         # List models
```

### Custom Model

Use a different model for testing:

```bash
go run cmd/openrouter-test/main.go -test chat -model anthropic/claude-3-haiku
```

For vision tests, use a vision-capable model:

```bash
go run cmd/openrouter-test/main.go -test image -model google/gemini-2.0-flash-thinking-exp:free
```

### Verbose Output

```bash
go run cmd/openrouter-test/main.go -test chat -v
```

## Image Tests

The test suite includes comprehensive image input tests:

### URL-Based Images
Tests sending images via public URLs to vision models.

### Multiple Images
Tests sending multiple images in a single request.

### Detail Levels
Tests the detail parameter (low/high/auto) for controlling image analysis quality.

### ContentBuilder
Tests the flexible builder API for constructing complex multimodal messages.

### Base64 Local Image
Tests encoding and sending local image files. Uses the included `test-image.png` file (colorful test tubes in a laboratory).

**Important:** The base64 image test must be run from the `cmd/openrouter-test` directory:

```bash
cd cmd/openrouter-test
go run . -test base64image -model google/gemini-2.0-flash-thinking-exp:free
```

## Test Image

The `test-image.png` file contains a photograph of five test tubes with colorful liquids (red, orange, yellow, green, and blue) in a laboratory setting. This image is used to test base64 encoding and local file upload functionality.

## Available Tests

| Test Name | Description |
|-----------|-------------|
| `all` | Run all tests |
| `chat` | Basic chat completion |
| `stream` | Streaming chat completion |
| `completion` | Legacy completion API |
| `error` | Error handling |
| `provider` | Provider routing |
| `zdr` | Zero Data Retention |
| `suffix` | Model suffix handling (`:nitro`, `:floor`, `:online`) |
| `price` | Price constraints |
| `structured` | Structured outputs with JSON schema |
| `tools` | Tool/function calling |
| `transforms` | Message transforms |
| `websearch` | Web search integration |
| `image` | Single image input (URL) |
| `multiimage` | Multiple images |
| `imagedetail` | Image with detail level |
| `contentbuilder` | ContentBuilder API |
| `base64image` | Base64-encoded local image |
| `models` | List available models |
| `endpoints` | Model endpoint information |
| `providers` | List providers |
| `credits` | Get credit balance |
| `activity` | Get usage activity |
| `key` | Get API key info |
| `listkeys` | List API keys |
| `createkey` | Create API key |
| `updatekey` | Update API key |
| `deletekey` | Delete API key |

## Notes

- Vision tests (image, multiimage, imagedetail, contentbuilder, base64image) require vision-capable models
- Some tests may be skipped if the selected model doesn't support the required features
- Rate limits may apply depending on your OpenRouter account tier
- The base64image test requires the test-image.png file in the working directory
