# openrouter-go

[![Go Report Card](https://goreportcard.com/badge/github.com/hra42/openrouter-go)](https://goreportcard.com/report/github.com/hra42/openrouter-go)
[![GoDoc](https://godoc.org/github.com/hra42/openrouter-go?status.svg)](https://godoc.org/github.com/hra42/openrouter-go)
[![License: Unlicense](https://img.shields.io/badge/license-Unlicense-blue.svg)](http://unlicense.org/)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hra42/openrouter-go)](https://github.com/hra42/openrouter-go)

A zero-dependency Go package providing complete bindings for the OpenRouter API, supporting all available endpoints with full streaming capabilities.

## Features

- ✅ Complete API coverage (chat completions, legacy completions, models, model endpoints, and providers)
- ✅ Full streaming support with Server-Sent Events (SSE)
- ✅ Zero external dependencies
- ✅ Go 1.25.1 support
- ✅ Comprehensive error handling and retry logic
- ✅ Context-aware cancellation
- ✅ Thread-safe client operations
- ✅ Extensive configuration options via functional options pattern
- ✅ Per-request Zero Data Retention (ZDR) enforcement
- ✅ Structured outputs with JSON schema validation
- ✅ Tool/Function calling support with streaming
- ✅ MCP (Model Context Protocol) tool conversion utilities
- ✅ Message transforms for automatic context window management
- ✅ Web Search plugin for real-time web data integration
- ✅ Image inputs (multimodal) with URL and base64 support
- ✅ Audio inputs with base64 encoding support (WAV, MP3)
- ✅ PDF inputs with configurable parsing engines and file annotation reuse
- ✅ Model listing and discovery with category filtering
- ✅ Model endpoint inspection with pricing and uptime details
- ✅ Provider listing with policy information
- ✅ Credit balance and usage tracking
- ✅ Activity analytics for usage monitoring and cost tracking
- ✅ API key information retrieval with usage and rate limit details
- ✅ API key management with listing, filtering, and creation capabilities
- ⚠️ **[BETA]** Responses API with reasoning, tool calling, web search, and streaming

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
    client := openrouter.NewClient(
        openrouter.WithAPIKey("your-api-key"),
    )

    messages := []openrouter.Message{
        {Role: "user", Content: "Hello, how are you?"},
    }

    response, err := client.ChatComplete(context.Background(),
        openrouter.WithModel("openai/gpt-4o"),
        openrouter.WithMessages(messages),
    )

    if err != nil {
        panic(err)
    }

    fmt.Println(response.Choices[0].Message.Content)
}
```

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

### Listing Available Models

```go
// List all available models
response, err := client.ListModels(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, model := range response.Data {
    fmt.Printf("%s - %s\n", model.ID, model.Name)
    fmt.Printf("  Context: %.0f tokens\n", *model.ContextLength)
    fmt.Printf("  Pricing: $%s/M prompt, $%s/M completion\n",
        model.Pricing.Prompt, model.Pricing.Completion)
}

// Filter models by category (e.g., "programming")
response, err := client.ListModels(ctx, &openrouter.ListModelsOptions{
    Category: "programming",
})
```

### Listing Model Endpoints

Get detailed information about the specific endpoints (providers) available for a model:

```go
// List all endpoints for a specific model
response, err := client.ListModelEndpoints(ctx, "openai", "gpt-4")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Model: %s\n", response.Data.Name)
fmt.Printf("Total endpoints: %d\n\n", len(response.Data.Endpoints))

// Examine each provider endpoint
for _, endpoint := range response.Data.Endpoints {
    fmt.Printf("Provider: %s\n", endpoint.ProviderName)
    fmt.Printf("  Status: %s\n", endpoint.Status)
    fmt.Printf("  Context Length: %.0f tokens\n", endpoint.ContextLength)
    fmt.Printf("  Pricing - Prompt: $%s/M, Completion: $%s/M\n",
        endpoint.Pricing.Prompt, endpoint.Pricing.Completion)

    if endpoint.UptimeLast30m != nil {
        fmt.Printf("  Uptime (30m): %.2f%%\n", *endpoint.UptimeLast30m*100)
    }

    if endpoint.Quantization != nil {
        fmt.Printf("  Quantization: %s\n", *endpoint.Quantization)
    }

    fmt.Printf("  Supported Parameters: %v\n\n", endpoint.SupportedParameters)
}
```

This endpoint is useful for:
- Comparing pricing across different providers for the same model
- Checking provider availability and uptime
- Finding endpoints with specific quantization levels
- Discovering which parameters are supported by each provider


### Listing Available Providers

Get information about all providers available through OpenRouter:

```go
// List all providers
response, err := client.ListProviders(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total providers: %d\n\n", len(response.Data))

// Display provider information
for _, provider := range response.Data {
    fmt.Printf("Provider: %s (%s)\n", provider.Name, provider.Slug)

    if provider.PrivacyPolicyURL != nil {
        fmt.Printf("  Privacy Policy: %s\n", *provider.PrivacyPolicyURL)
    }

    if provider.TermsOfServiceURL != nil {
        fmt.Printf("  Terms of Service: %s\n", *provider.TermsOfServiceURL)
    }

    if provider.StatusPageURL != nil {
        fmt.Printf("  Status Page: %s\n", *provider.StatusPageURL)
    }
}
```

This endpoint is useful for:
- Reviewing provider policies and terms
- Finding provider status pages for uptime monitoring
- Understanding which providers are available
- Checking provider compliance information

### Getting Credit Balance

Retrieve your current credit balance and usage for the authenticated user:

```go
// Get credit balance
response, err := client.GetCredits(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Credits: $%.2f\n", response.Data.TotalCredits)
fmt.Printf("Total Usage: $%.2f\n", response.Data.TotalUsage)

// Calculate remaining balance
remaining := response.Data.TotalCredits - response.Data.TotalUsage
fmt.Printf("Remaining: $%.2f\n", remaining)

// Check usage percentage
if response.Data.TotalCredits > 0 {
    usagePercent := (response.Data.TotalUsage / response.Data.TotalCredits) * 100
    fmt.Printf("Usage: %.2f%%\n", usagePercent)
}
```

This endpoint is useful for:
- Monitoring credit consumption in real-time
- Setting up alerts for low balance
- Tracking API usage costs
- Budget management and forecasting

### Getting Activity Data

Retrieve daily user activity data grouped by model endpoint for the last 30 (completed) UTC days:

```go
// Get all activity data
response, err := client.GetActivity(ctx, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total activity records: %d\n\n", len(response.Data))

// Calculate summary statistics
totalUsage := 0.0
totalRequests := 0.0

for _, data := range response.Data {
    totalUsage += data.Usage
    totalRequests += data.Requests
}

fmt.Printf("Total usage: $%.4f\n", totalUsage)
fmt.Printf("Total requests: %.0f\n", totalRequests)

// Filter by specific date
yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
dateActivity, err := client.GetActivity(ctx, &openrouter.ActivityOptions{
    Date: yesterday,
})
if err != nil {
    log.Fatal(err)
}

// Display activity for specific date
for _, data := range dateActivity.Data {
    fmt.Printf("Date: %s\n", data.Date)
    fmt.Printf("Model: %s\n", data.Model)
    fmt.Printf("Provider: %s\n", data.ProviderName)
    fmt.Printf("Requests: %.0f\n", data.Requests)
    fmt.Printf("Usage: $%.4f\n", data.Usage)
    fmt.Printf("Tokens: %.0f prompt, %.0f completion, %.0f reasoning\n",
        data.PromptTokens, data.CompletionTokens, data.ReasoningTokens)
}
```

**Important**: This endpoint requires a provisioning key (not a regular inference API key). Create one at: https://openrouter.ai/settings/provisioning-keys

This endpoint is useful for:
- Daily usage analytics and cost tracking
- Model performance comparison
- Provider usage distribution analysis
- Historical cost analysis and forecasting
- BYOK (Bring Your Own Key) usage tracking

### Getting API Key Information

Retrieve information about your current API key including usage, limits, and rate limits:

```go
// Get API key information
response, err := client.GetKey(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("API Key Label: %s\n", response.Data.Label)

// Display limit information
if response.Data.Limit != nil {
    fmt.Printf("Credit Limit: $%.2f\n", *response.Data.Limit)
} else {
    fmt.Printf("Credit Limit: Unlimited\n")
}

fmt.Printf("Usage: $%.4f\n", response.Data.Usage)

// Display remaining balance
if response.Data.LimitRemaining != nil {
    fmt.Printf("Remaining: $%.4f\n", *response.Data.LimitRemaining)

    // Calculate usage percentage
    if response.Data.Limit != nil && *response.Data.Limit > 0 {
        usagePercent := (response.Data.Usage / *response.Data.Limit) * 100
        fmt.Printf("Usage: %.2f%%\n", usagePercent)
    }
}

// Display key type
fmt.Printf("Free Tier: %v\n", response.Data.IsFreeTier)
fmt.Printf("Provisioning Key: %v\n", response.Data.IsProvisioningKey)

// Display rate limit if available
if response.Data.RateLimit != nil {
    fmt.Printf("Rate Limit: %.0f requests per %s\n",
        response.Data.RateLimit.Requests,
        response.Data.RateLimit.Interval)
}
```

This endpoint is useful for:
- Monitoring API key usage and limits
- Checking remaining credits
- Understanding rate limit constraints
- Identifying key type (free tier vs paid, inference vs provisioning)
- Building usage alerts and notifications

### Listing All API Keys

Retrieve a list of all API keys associated with your account. Requires a Provisioning API key (not a regular inference API key):

```go
// List all API keys
response, err := client.ListKeys(ctx, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total API keys: %d\n\n", len(response.Data))

// Display key information
for i, key := range response.Data {
    status := "Active"
    if key.Disabled {
        status = "Disabled"
    }

    fmt.Printf("%d. %s [%s]\n", i+1, key.Label, status)
    fmt.Printf("   Name: %s\n", key.Name)
    fmt.Printf("   Limit: $%.2f\n", key.Limit)
    fmt.Printf("   Created: %s\n", key.CreatedAt)
    fmt.Printf("   Updated: %s\n", key.UpdatedAt)
}

// Example with pagination and filtering
offset := 10
includeDisabled := true
filteredKeys, err := client.ListKeys(ctx, &openrouter.ListKeysOptions{
    Offset:          &offset,
    IncludeDisabled: &includeDisabled,
})
```

**Important**: This endpoint requires a provisioning key (not a regular inference API key). Create one at: https://openrouter.ai/settings/provisioning-keys

This endpoint is useful for:
- Managing multiple API keys programmatically
- Auditing key usage and creation dates
- Identifying and managing disabled keys
- Implementing key rotation strategies
- Building API key management dashboards

### Getting API Key by Hash

Retrieve details about a specific API key by its hash. Requires a Provisioning API key:

```go
// Get key details by hash (hash obtained from ListKeys or key creation)
hash := "abc123hash"
keyDetails, err := client.GetKeyByHash(ctx, hash)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Label: %s\n", keyDetails.Data.Label)
fmt.Printf("Name: %s\n", keyDetails.Data.Name)
fmt.Printf("Limit: $%.2f\n", keyDetails.Data.Limit)
fmt.Printf("Disabled: %v\n", keyDetails.Data.Disabled)
fmt.Printf("Created: %s\n", keyDetails.Data.CreatedAt)
fmt.Printf("Updated: %s\n", keyDetails.Data.UpdatedAt)

// Example: Get hash from list and retrieve details
keys, err := client.ListKeys(ctx, nil)
if err != nil {
    log.Fatal(err)
}

if len(keys.Data) > 0 {
    firstHash := keys.Data[0].Hash
    details, err := client.GetKeyByHash(ctx, firstHash)
    // ...
}
```

**Important**: This endpoint requires a provisioning key (not a regular inference API key). Create one at: https://openrouter.ai/settings/provisioning-keys

This endpoint is useful for:
- Inspecting individual API key details
- Verifying key status and configuration
- Monitoring specific key usage patterns
- Building key detail views in dashboards
- Auditing key configuration changes

### Creating API Keys

Create new API keys programmatically with custom limits and settings. Requires a Provisioning API key:

```go
// Create an API key with a credit limit
limit := 100.0
keyResp, err := client.CreateKey(ctx, &openrouter.CreateKeyRequest{
    Name:  "Production API Key",
    Limit: &limit,
})
if err != nil {
    log.Fatal(err)
}

// ⚠️ IMPORTANT: Save this value immediately!
// This is the ONLY time the full API key will be returned
fmt.Printf("New API Key: %s\n", keyResp.Key)
fmt.Printf("Label: %s\n", keyResp.Data.Label)
fmt.Printf("Limit: $%.2f\n", keyResp.Data.Limit)

// Create a key with BYOK limit inclusion
includeBYOK := true
keyResp2, err := client.CreateKey(ctx, &openrouter.CreateKeyRequest{
    Name:               "BYOK Key",
    Limit:              &limit,
    IncludeBYOKInLimit: &includeBYOK,
})

// Create a key without a specific limit (uses account limit)
keyResp3, err := client.CreateKey(ctx, &openrouter.CreateKeyRequest{
    Name: "Unlimited Key",
})
```

**Critical Security Note**: The `Key` field in the response contains the actual API key value. This is the **ONLY** time this value will ever be returned. Store it securely immediately!

**Important**: This endpoint requires a provisioning key (not a regular inference API key). Create one at: https://openrouter.ai/settings/provisioning-keys

This endpoint is useful for:
- Automated API key provisioning
- Implementing key rotation workflows
- Creating keys with custom credit limits
- Setting up BYOK (Bring Your Own Key) configurations
- Building self-service key management systems

### Deleting API Keys

Delete an API key by its hash. Requires a Provisioning API key:

```go
// Delete a key by hash (hash obtained from ListKeys or key creation)
hash := "abc123hash"
result, err := client.DeleteKey(ctx, hash)
if err != nil {
    log.Fatal(err)
}

if result.Data.Success {
    fmt.Println("API key successfully deleted")
}

// Example: Delete a key created in the same session
keyResp, err := client.CreateKey(ctx, &openrouter.CreateKeyRequest{
    Name: "Temporary Key",
})
if err != nil {
    log.Fatal(err)
}

// Later... delete it
deleteResult, err := client.DeleteKey(ctx, keyResp.Data.Hash)
if err != nil {
    log.Fatal(err)
}
```

**⚠️ WARNING**: This operation is **irreversible**! Once deleted, the API key cannot be recovered and any applications using it will immediately lose access.

**Important**: This endpoint requires a provisioning key (not a regular inference API key). Create one at: https://openrouter.ai/settings/provisioning-keys

This endpoint is useful for:
- Automated key rotation and cleanup
- Removing compromised or unused keys
- Implementing temporary key workflows
- Building key lifecycle management systems
- Programmatic key revocation

### Updating API Keys

Update an existing API key's properties (name, limit, disabled status) by its hash. Requires a Provisioning API key:

```go
// Update just the key name
hash := "abc123hash"
newName := "Updated Production Key"
result, err := client.UpdateKey(ctx, hash, &openrouter.UpdateKeyRequest{
    Name: &newName,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated key: %s\n", result.Data.Label)

// Disable a key
disabled := true
result, err := client.UpdateKey(ctx, hash, &openrouter.UpdateKeyRequest{
    Disabled: &disabled,
})

// Update credit limit
newLimit := 200.0
result, err := client.UpdateKey(ctx, hash, &openrouter.UpdateKeyRequest{
    Limit: &newLimit,
})

// Update multiple fields at once
result, err := client.UpdateKey(ctx, hash, &openrouter.UpdateKeyRequest{
    Name:               &newName,
    Limit:              &newLimit,
    IncludeBYOKInLimit: &[]bool{true}[0],
})
```

**Important**: This endpoint requires a provisioning key (not a regular inference API key). Create one at: https://openrouter.ai/settings/provisioning-keys

All fields in `UpdateKeyRequest` are optional - only include the fields you want to update:
- **Name**: Display name for the API key
- **Disabled**: Set to true to disable the key (prevents usage)
- **Limit**: Credit limit in dollars
- **IncludeBYOKInLimit**: Whether BYOK (Bring Your Own Key) usage counts toward the limit

This endpoint is useful for:
- Rotating key names for better organization
- Adjusting credit limits based on usage patterns
- Temporarily disabling keys without deletion
- Managing BYOK limit policies
- Implementing dynamic key management workflows

## Code Quality

The library is built with a focus on code quality and maintainability:

- **Named Constants**: Magic numbers have been extracted as named constants for better readability and maintainability
  - `defaultJitterFactor` (0.25): Default jitter factor for retry backoff (±25%)
  - `maxReconnectBackoff` (10s): Maximum backoff duration for stream reconnection attempts
  - `defaultMaxDelay` (30s): Default maximum delay for retry backoff
  - `defaultMultiplier` (2.0): Default multiplier for exponential backoff

- **Generic Options Pattern**: Uses Go 1.18+ generics to reduce code duplication in functional options
  - Type-safe option setters with `RequestConfig` interface constraint
  - Shared implementation for common fields across ChatCompletion and Completion requests
  - Eliminates ~400 lines of duplicate code while maintaining type safety

- **Comprehensive Testing**: Extensive unit test coverage with table-driven tests
- **Race Detection**: All code is tested for race conditions
- **Thread Safety**: Client is safe for concurrent use across goroutines
- **Error Handling**: Rich error types with detailed context

## Package Structure

```
openrouter-go/
├── client.go            # Main client implementation
├── completions.go       # Completion endpoint methods
├── chat.go              # Chat completion endpoint methods
├── models_endpoint.go   # Models listing endpoint methods
├── model_endpoints.go   # Model endpoints inspection methods
├── providers_endpoint.go # Providers listing endpoint methods
├── credits_endpoint.go  # Credits balance endpoint methods
├── activity_endpoint.go # Activity analytics endpoint methods
├── key_endpoint.go      # API key information endpoint methods
├── models.go            # Request/response type definitions
├── options.go           # Functional options for configuration
├── stream.go            # SSE streaming with generic Stream[T] implementation
├── errors.go            # Custom error types
├── retry.go             # Retry and backoff logic with named constants
├── mcp.go               # MCP tool conversion utilities
├── responses.go         # [BETA] Responses API endpoint methods
├── responses_models.go  # [BETA] Responses API type definitions
├── responses_options.go # [BETA] Responses API functional options
├── examples/
│   ├── basic/             # Basic usage examples
│   ├── streaming/         # Streaming examples
│   ├── structured-output/ # Structured outputs with JSON schema
│   ├── tool-calling/      # Tool/function calling examples
│   ├── mcp-tools/         # MCP tool conversion examples
│   ├── web_search/        # Web search plugin examples
│   ├── list-models/       # Model listing examples
│   ├── model-endpoints/   # Model endpoints inspection examples
│   ├── list-providers/    # Provider listing examples
│   ├── get-credits/       # Credit balance tracking examples
│   ├── activity/          # Activity analytics examples
│   ├── key/               # API key information examples
│   ├── list-keys/         # API key listing examples
│   ├── create-key/        # API key creation examples
│   ├── responses/         # [BETA] Responses API examples
│   └── advanced/          # Advanced configuration examples
└── internal/
    └── sse/               # Internal SSE parser implementation
```

## App Attribution

Get your app featured in OpenRouter rankings and analytics by including attribution headers:

```go
client := openrouter.NewClient(
    openrouter.WithAPIKey("your-api-key"),
    // Your app's URL (primary identifier)
    openrouter.WithReferer("https://myapp.com"),
    // Your app's display name
    openrouter.WithAppName("My AI Assistant"),
)
```

### Benefits

When you use app attribution, your app will:
- Appear in [OpenRouter's public rankings](https://openrouter.ai/rankings)
- Be featured on individual model pages in the "Apps" tab
- Get detailed analytics at `openrouter.ai/apps?url=<your-app-url>`
- Gain visibility in the OpenRouter developer community

### Localhost Development

For localhost development, always include a title:

```go
client := openrouter.NewClient(
    openrouter.WithAPIKey("your-api-key"),
    openrouter.WithReferer("http://localhost:3000"),
    openrouter.WithAppName("Development App"), // Required for localhost
)
```

See the [app attribution example](examples/app-attribution/main.go) for more details.

## Requirements

- Go 1.25.1
- No external dependencies

## Status

✅ **Production Ready** - All 5 phases complete! The library is now ready for production use with:
- ✅ Full foundation with all types and error handling
- ✅ Robust HTTP communication with retry logic
- ✅ Complete API implementation for chat and completions
- ✅ Zero-dependency SSE streaming with reconnection support
- ✅ Comprehensive test coverage and documentation
- ✅ Production-ready examples for all use cases

## Testing

### Unit Tests

Run the unit test suite:

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

### E2E Tests

The project includes a comprehensive end-to-end test suite in `cmd/openrouter-test/` that tests against the live OpenRouter API. The test suite is organized into logical modules:

**Test Structure:**
```
cmd/openrouter-test/
├── main.go              # Entry point and CLI
└── tests/
    ├── helpers.go       # Shared utilities
    ├── chat.go          # Chat, streaming, completion tests
    ├── routing.go       # Provider routing, ZDR, model suffixes
    ├── structured.go    # Structured output tests
    ├── tools.go         # Tool/function calling tests
    ├── transforms.go    # Message transforms tests
    ├── search.go        # Web search tests
    ├── models.go        # Models, endpoints, providers tests
    ├── account.go       # Credits, activity tests
    └── apikeys.go       # API key management tests
```

**Running E2E Tests:**

```bash
# Set your API key
export OPENROUTER_API_KEY="your-api-key"

# Run all tests (excluding web search)
go run cmd/openrouter-test/main.go -test all

# Run specific test categories
go run cmd/openrouter-test/main.go -test chat
go run cmd/openrouter-test/main.go -test streaming
go run cmd/openrouter-test/main.go -test tools
go run cmd/openrouter-test/main.go -test websearch  # Run separately on demand

# Run with custom model
go run cmd/openrouter-test/main.go -test all -model anthropic/claude-3-haiku

# Run with verbose output
go run cmd/openrouter-test/main.go -test chat -v

# Available tests:
# all, chat, stream, completion, error, provider, zdr, suffix,
# price, structured, tools, transforms, websearch, models,
# endpoints, providers, credits, activity, key, listkeys,
# createkey, updatekey, deletekey
```

### Message Transforms

The library supports message transforms to automatically handle prompts that exceed a model's context window. This feature uses "middle-out" compression to remove content from the middle of long prompts where models typically pay less attention.

#### Basic Transform Usage

```go
// Enable middle-out compression for chat completions
response, err := client.ChatComplete(ctx,
    openrouter.WithModel("meta-llama/llama-3.1-8b-instruct"),
    openrouter.WithMessages(messages),
    openrouter.WithTransforms("middle-out"), // Auto-compress if exceeds context
)

// Enable for legacy completions
response, err := client.Complete(ctx, prompt,
    openrouter.WithModel("openai/gpt-3.5-turbo-instruct"),
    openrouter.WithCompletionTransforms("middle-out"),
)
```

#### How It Works

When `middle-out` transform is enabled:
1. OpenRouter finds models with at least half of your required tokens (input + completion)
2. If your prompt exceeds the model's context, content is removed from the middle
3. For models with message count limits (e.g. Anthropic's Claude), messages are compressed to stay within limits

#### Default Behavior

All OpenRouter endpoints with 8K (8,192 tokens) or less context length automatically use `middle-out` by default. To disable:

```go
// Explicitly disable transforms for smaller models
response, err := client.ChatComplete(ctx,
    openrouter.WithModel("some-8k-model"),
    openrouter.WithMessages(messages),
    openrouter.WithTransforms(), // Empty array disables transforms
)
```

#### When to Use

Message transforms are useful when:
- Perfect recall is not required
- You want automatic fallback for long conversations
- Working with models that have smaller context windows
- Handling variable-length user inputs that might exceed limits

#### Important Notes

- Middle content is compressed because LLMs pay less attention to the middle of sequences
- The transform handles both token limits and message count limits
- Without transforms, requests exceeding limits will fail with an error
- Consider using models with larger context windows if perfect recall is critical

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

### Image Inputs (Multimodal)

The library provides comprehensive support for sending images to vision models through the OpenRouter API. You can send images via URLs or base64-encoded data.

#### Single Image with URL

```go
// Send a single image with text
messages := []openrouter.Message{
    openrouter.CreateUserMessageWithImage(
        "What's in this image?",
        "https://example.com/image.jpg",
    ),
}

response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("google/gemini-2.0-flash-thinking-exp:free"),
)
```

#### Multiple Images

```go
// Send multiple images in a single request
messages := []openrouter.Message{
    openrouter.CreateUserMessageWithImages(
        "Compare these images. What are the similarities?",
        "https://example.com/image1.jpg",
        "https://example.com/image2.jpg",
        "https://example.com/image3.jpg",
    ),
}

response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("google/gemini-2.0-flash-thinking-exp:free"),
)
```

#### Image with Detail Level

Some models support detail level parameters for controlling image analysis quality:

```go
// Request high-detail analysis (more expensive, more detailed)
messages := []openrouter.Message{
    openrouter.CreateUserMessageWithImageDetail(
        "Describe this image in detail.",
        "https://example.com/image.jpg",
        "high", // Options: "low", "high", or "auto"
    ),
}

response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("google/gemini-2.0-flash-thinking-exp:free"),
)
```

Detail level options:
- `"low"` - Faster and cheaper, suitable for general understanding
- `"high"` - More detailed analysis at higher cost
- `"auto"` - Let the model decide based on image size (default)

#### Base64-Encoded Images

For local files or private images:

```go
// Automatically encode and send a local image
message, err := openrouter.CreateUserMessageWithBase64Image(
    "What's in this image?",
    "path/to/image.jpg",
)
if err != nil {
    log.Fatal(err)
}

messages := []openrouter.Message{message}

response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("google/gemini-2.0-flash-thinking-exp:free"),
)
```

Multiple local images:

```go
message, err := openrouter.CreateUserMessageWithBase64Images(
    "Compare these images",
    "path/to/image1.jpg",
    "path/to/image2.png",
)
```

Manual base64 encoding:

```go
// Encode image file to base64 data URL
dataURL, err := openrouter.EncodeImageToBase64("path/to/image.jpg")
if err != nil {
    log.Fatal(err)
}

// Or encode bytes directly
imageBytes := []byte{...}
dataURL := openrouter.EncodeImageBytesToBase64(imageBytes, "image/jpeg")
```

#### Content Builder for Complex Messages

For messages with interleaved text and images:

```go
content := openrouter.NewContentBuilder().
    AddText("Here's the first image:").
    AddImage("https://example.com/image1.jpg").
    AddText("And here's the second with high detail:").
    AddImageWithDetail("https://example.com/image2.jpg", "high").
    AddText("What are the differences?")

messages := []openrouter.Message{
    content.BuildMessage("user"),
}

response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("google/gemini-2.0-flash-thinking-exp:free"),
)
```

#### Supported Image Formats

- PNG (image/png)
- JPEG (image/jpeg)
- WebP (image/webp)
- GIF (image/gif)

#### Model Support

Most modern vision models support image inputs, including:
- Google Gemini models (gemini-2.0-flash-thinking-exp, etc.)
- OpenAI GPT-4 Vision models
- Anthropic Claude 3 models
- And many others

Check the [OpenRouter models page](https://openrouter.ai/models) for the latest list of vision-capable models.

#### Best Practices

- Use URLs when possible - they're more efficient than base64 encoding
- Image URLs must be publicly accessible
- The number of images per request varies by model and provider
- Some providers may have size limits on images
- Pricing may vary based on image size and detail level
- For production use, consider the OpenRouter documentation's recommendations about image placement in messages

See the [image-inputs example](examples/image-inputs) for more comprehensive examples.

### PDF Inputs (File Support)

The library provides comprehensive support for sending PDF files to models through the OpenRouter API. PDF files can be sent via URLs or base64-encoded data. This feature works with **any model** on OpenRouter, with automatic fallback to PDF parsing when models don't have native file support.

#### Basic PDF from URL

```go
// Send a PDF via URL
messages := []openrouter.Message{
    openrouter.CreateUserMessageWithPDF(
        "What are the main points in this document?",
        "https://bitcoin.org/bitcoin.pdf",
        "bitcoin.pdf",
    ),
}

response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("anthropic/claude-sonnet-4"),
)
```

#### Local PDF Files

```go
// Automatically encode and send a local PDF
message, err := openrouter.CreateUserMessageWithBase64PDF(
    "Summarize this document",
    "path/to/document.pdf",
    "document.pdf",
)
if err != nil {
    log.Fatal(err)
}

response, err := client.ChatComplete(ctx, []openrouter.Message{message},
    openrouter.WithModel("google/gemma-3-27b-it"),
)
```

#### PDF Parsing Engines

OpenRouter provides three PDF processing engines with different cost/quality tradeoffs:

```go
message := openrouter.CreateUserMessageWithPDF(
    "Extract key concepts from this document",
    "https://example.com/document.pdf",
    "document.pdf",
)

// Configure the PDF parsing engine
plugin := openrouter.CreateFileParserPlugin("pdf-text")

response, err := client.ChatComplete(ctx, []openrouter.Message{message},
    openrouter.WithModel("google/gemma-3-27b-it"),
    openrouter.WithPlugins(plugin),
)
```

Available engines:
- `"pdf-text"` - Free, best for well-structured PDFs with clear text content
- `"mistral-ocr"` - $0.0004 per 1K pages, best for scanned documents with images/OCR needs
- `"native"` - Uses model's native file support (charged as input tokens)
- `""` (empty) - Auto-selects native support first, then defaults to `pdf-text`

#### Reusing File Annotations

File annotations allow you to avoid re-parsing the same PDF in follow-up requests, saving processing time and costs:

```go
// First request with PDF
firstMessage := openrouter.CreateUserMessageWithPDF(
    "What are the main concepts in this paper?",
    "https://example.com/document.pdf",
    "document.pdf",
)

resp1, err := client.ChatComplete(ctx, []openrouter.Message{firstMessage},
    openrouter.WithModel("google/gemma-3-27b-it"),
)

// Follow-up request - include the assistant's response with annotations
followUpMessages := []openrouter.Message{
    firstMessage,
    resp1.Choices[0].Message, // Contains file annotations
    openrouter.CreateUserMessage("Can you elaborate on the first point?"),
}

resp2, err := client.ChatComplete(ctx, followUpMessages,
    openrouter.WithModel("google/gemma-3-27b-it"),
)
// PDF is NOT re-parsed - saves processing time and costs!
```

#### Multiple Files

You can send multiple files (PDFs, images, etc.) in a single request:

```go
files := []openrouter.File{
    {
        Filename: "document1.pdf",
        FileData: "https://example.com/doc1.pdf",
    },
    {
        Filename: "document2.pdf",
        FileData: "https://example.com/doc2.pdf",
    },
    {
        Filename: "chart.png",
        FileData: "https://example.com/chart.png",
    },
}

message := openrouter.CreateUserMessageWithFiles(
    "Compare these documents and analyze the chart",
    files,
)

response, err := client.ChatComplete(ctx, []openrouter.Message{message},
    openrouter.WithModel("anthropic/claude-sonnet-4"),
)
```

#### Content Builder with PDFs

For complex messages with PDFs, images, and text:

```go
content := openrouter.NewContentBuilder().
    AddText("Analyze this document:").
    AddPDF("https://example.com/document.pdf", "document.pdf").
    AddText("And compare with this image:").
    AddImage("https://example.com/chart.png").
    Build()

message := openrouter.Message{
    Role:    "user",
    Content: content,
}

response, err := client.ChatComplete(ctx, []openrouter.Message{message},
    openrouter.WithModel("anthropic/claude-sonnet-4"),
)
```

#### Manual Base64 Encoding

```go
// Encode PDF file to base64 data URL
dataURL, err := openrouter.EncodePDFToBase64("path/to/document.pdf")
if err != nil {
    log.Fatal(err)
}

// Or encode bytes directly
pdfBytes := []byte{...}
dataURL := openrouter.EncodePDFBytesToBase64(pdfBytes)
```

#### Best Practices

- **Use URLs when possible** - More efficient than base64 encoding
- **Use `pdf-text` for digital PDFs** - It's free and works well for most documents
- **Reuse file annotations** - Include the assistant's response with annotations in follow-up requests
- **Use `mistral-ocr` only for scanned documents** - More expensive but necessary for image-based PDFs
- **Check model support** - Some models have native file support which may be more cost-effective

#### Supported File Types

While this section focuses on PDFs, the file input API supports:
- PDFs (application/pdf)
- Images (image/png, image/jpeg, etc.)
- And potentially other file types as OpenRouter expands support

See the [pdf-inputs example](examples/pdf-inputs) for more comprehensive examples.

### Text File Inputs

Send text-based files (code, configuration, documentation) to models using inline content. Text files are sent directly as UTF-8 text, not base64-encoded, for efficiency and token optimization.

#### Single Text File

```go
// Send a code file for review
message, err := openrouter.CreateUserMessageWithTextFile(
    "Review this code for bugs:",
    "/path/to/code.py",
)
if err != nil {
    log.Fatal(err)
}

response, err := client.ChatComplete(ctx, []openrouter.Message{message},
    openrouter.WithModel("anthropic/claude-sonnet-4"),
)
```

#### Multiple Text Files

```go
// Compare or analyze multiple files
message, err := openrouter.CreateUserMessageWithTextFiles(
    "Compare these configuration files and identify differences:",
    "/path/to/config1.yaml",
    "/path/to/config2.yaml",
)
if err != nil {
    log.Fatal(err)
}

response, err := client.ChatComplete(ctx, []openrouter.Message{message},
    openrouter.WithModel("openai/gpt-4"),
)
```

#### Content Builder with Text Files

```go
// Build complex messages with multiple text files
builder := openrouter.NewContentBuilder()
builder.AddText("I have several files to review:")

builder, err := builder.AddTextFile("/path/to/main.go")
if err != nil {
    log.Fatal(err)
}

builder.AddText("And here's the test file:")
builder, err = builder.AddTextFile("/path/to/main_test.go")
if err != nil {
    log.Fatal(err)
}

builder.AddText("Do these files work together correctly?")

message := builder.BuildMessage("user")
response, err := client.ChatComplete(ctx, []openrouter.Message{message},
    openrouter.WithModel("anthropic/claude-sonnet-4"),
)
```

#### Direct Text Content (Without File I/O)

```go
// Send text content directly without reading from a file
message := openrouter.CreateUserMessageWithTextContent(
    "Analyze this JSON structure:",
    `{
  "name": "example",
  "version": "1.0.0",
  "dependencies": {}
}`,
    "package.json",
)

response, err := client.ChatComplete(ctx, []openrouter.Message{message},
    openrouter.WithModel("openai/gpt-4o-mini"),
)
```

#### Supported File Formats

**Common text formats**
- `.txt` - Plain text
- `.md` - Markdown
- `.json` - JSON files
- `.csv` - CSV files

**Code files**
- `.js`, `.jsx` - JavaScript
- `.ts`, `.tsx` - TypeScript
- `.py` - Python
- `.go` - Go
- `.java` - Java
- `.rs` - Rust
- `.c`, `.cpp`, `.h` - C/C++
- `.rb` - Ruby
- `.php` - PHP
- `.swift` - Swift
- `.kt` - Kotlin

**Configuration files**
- `.yaml`, `.yml` - YAML
- `.toml` - TOML
- `.xml` - XML
- `.ini` - INI
- `.env` - Environment files

And many more! See [text-file-inputs example](examples/text-file-inputs) for complete details.

#### Key Features

- **Inline delivery** - Text sent directly (not base64), optimizing tokens
- **UTF-8 validation** - All files must contain valid UTF-8 text
- **Filename context** - Files are sent with filename headers for context
- **All models supported** - Works with any text-capable model on OpenRouter
- **Format validation** - Unsupported formats rejected with clear errors
- **Builder integration** - Seamlessly mix text files with other content types

### Structured Outputs

The library supports structured outputs for compatible models, ensuring responses follow a specific JSON Schema format. This feature is useful when you need consistent, well-formatted responses that can be reliably parsed by your application.

#### Basic Structured Output

```go
// Define a JSON schema for the expected response
weatherSchema := map[string]interface{}{
    "type": "object",
    "properties": map[string]interface{}{
        "location": map[string]interface{}{
            "type": "string",
            "description": "City or location name",
        },
        "temperature": map[string]interface{}{
            "type": "number",
            "description": "Temperature in Celsius",
        },
        "conditions": map[string]interface{}{
            "type": "string",
            "description": "Weather conditions",
        },
    },
    "required": []string{"location", "temperature", "conditions"},
    "additionalProperties": false,
}

// Use structured output with chat completion
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithJSONSchema("weather", true, weatherSchema),
    openrouter.WithRequireParameters(true), // Ensure model supports structured outputs
)

// The response will be valid JSON matching your schema
var weatherData map[string]interface{}
json.Unmarshal([]byte(response.Choices[0].Message.Content.(string)), &weatherData)
```

#### Simplified JSON Mode

```go
// For simpler cases, use JSON mode without a strict schema
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithJSONMode(), // Returns JSON without enforcing a schema
)
```

#### Streaming with Structured Output

```go
// Structured outputs work with streaming too
stream, err := client.ChatCompleteStream(ctx, messages,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithJSONSchema("response", true, schema),
)

var fullContent string
for event := range stream.Events() {
    if len(event.Choices) > 0 && event.Choices[0].Delta != nil {
        if content, ok := event.Choices[0].Delta.Content.(string); ok {
            fullContent += content
        }
    }
}

// Parse the complete JSON response
var result map[string]interface{}
json.Unmarshal([]byte(fullContent), &result)
```

#### Model Support

Not all models support structured outputs. To ensure compatibility:

1. Check the [models page](https://openrouter.ai/models?supported_parameters=structured_outputs) for support
2. Use `WithRequireParameters(true)` to route only to compatible providers
3. Models known to support structured outputs include:
   - OpenAI models (GPT-4o and later)
   - Many Fireworks-provided models

#### Best Practices

- Always set `strict: true` in your JSON schema for exact compliance
- Include clear descriptions in schema properties to guide the model
- Use `WithRequireParameters(true)` to ensure routing to compatible providers
- Test your schemas with the specific models you plan to use
- Handle parsing errors gracefully as a fallback

### Tool/Function Calling

The library provides full support for tool/function calling, allowing models to use external tools and functions during generation. This feature enables building powerful AI agents and assistants.

#### Basic Tool Calling

```go
// Define a tool
tools := []openrouter.Tool{
    {
        Type: "function",
        Function: openrouter.Function{
            Name:        "get_weather",
            Description: "Get the current weather for a location",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "location": map[string]interface{}{
                        "type":        "string",
                        "description": "City name or zip code",
                    },
                    "unit": map[string]interface{}{
                        "type":        "string",
                        "enum":        []string{"celsius", "fahrenheit"},
                        "description": "Temperature unit",
                    },
                },
                "required": []string{"location"},
            },
        },
    },
}

// Make a request with tools
messages := []openrouter.Message{
    {Role: "user", Content: "What's the weather in San Francisco?"},
}

response, err := client.ChatComplete(ctx,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithMessages(messages),
    openrouter.WithTools(tools),
)

// Check for tool calls in the response
if len(response.Choices[0].Message.ToolCalls) > 0 {
    // Process tool calls
    for _, toolCall := range response.Choices[0].Message.ToolCalls {
        // Parse arguments
        var args map[string]interface{}
        json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

        // Execute the tool (your implementation)
        result := executeWeatherTool(args)

        // Add tool result to messages
        messages = append(messages, response.Choices[0].Message)
        messages = append(messages, openrouter.Message{
            Role:       "tool",
            Content:    result,
            ToolCallID: toolCall.ID,
        })
    }

    // Get final response with tool results
    finalResponse, _ := client.ChatComplete(ctx,
        openrouter.WithModel("openai/gpt-4o"),
        openrouter.WithMessages(messages),
        openrouter.WithTools(tools),
    )
}
```

#### Tool Choice Control

```go
// Let the model decide (default)
response, _ := client.ChatComplete(ctx,
    openrouter.WithMessages(messages),
    openrouter.WithTools(tools),
    openrouter.WithToolChoice("auto"),
)

// Disable tool usage
response, _ := client.ChatComplete(ctx,
    openrouter.WithMessages(messages),
    openrouter.WithTools(tools),
    openrouter.WithToolChoice("none"),
)

// Force specific tool usage
response, _ := client.ChatComplete(ctx,
    openrouter.WithMessages(messages),
    openrouter.WithTools(tools),
    openrouter.WithToolChoice(map[string]interface{}{
        "type": "function",
        "function": map[string]interface{}{
            "name": "get_weather",
        },
    }),
)
```

#### Parallel Tool Calls

Control whether multiple tools can be called simultaneously:

```go
// Disable parallel tool calls (sequential only)
parallelCalls := false
response, _ := client.ChatComplete(ctx,
    openrouter.WithMessages(messages),
    openrouter.WithTools(tools),
    openrouter.WithParallelToolCalls(&parallelCalls),
)
```

#### Streaming with Tool Calls

Tool calls are fully supported in streaming mode:

```go
stream, err := client.ChatCompleteStream(ctx,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithMessages(messages),
    openrouter.WithTools(tools),
)

var toolCalls []openrouter.ToolCall
for event := range stream.Events() {
    // Parse streaming data
    var data map[string]interface{}
    json.Unmarshal([]byte(event.Data), &data)

    if choices, ok := data["choices"].([]interface{}); ok && len(choices) > 0 {
        choice := choices[0].(map[string]interface{})

        // Check for tool calls in delta
        if delta, ok := choice["delta"].(map[string]interface{}); ok {
            if toolCallsDelta, ok := delta["tool_calls"].([]interface{}); ok {
                // Accumulate tool call information
                // See examples/tool-calling/streaming.go for complete implementation
            }
        }

        // Check finish reason
        if finishReason, ok := choice["finish_reason"].(string); ok {
            if finishReason == "tool_calls" {
                // Process accumulated tool calls
            }
        }
    }
}
```

#### Multi-Tool Workflows

Design tools that work well together:

```go
tools := []openrouter.Tool{
    {
        Type: "function",
        Function: openrouter.Function{
            Name:        "search_products",
            Description: "Search for products in the catalog",
            // Parameters...
        },
    },
    {
        Type: "function",
        Function: openrouter.Function{
            Name:        "check_inventory",
            Description: "Check inventory for a product",
            // Parameters...
        },
    },
    {
        Type: "function",
        Function: openrouter.Function{
            Name:        "place_order",
            Description: "Place an order for a product",
            // Parameters...
        },
    },
}

// The model can chain these tools naturally:
// search → check inventory → place order
```

#### Model Support

Tool calling is supported by many models. You can find compatible models by filtering on [openrouter.ai/models?supported_parameters=tools](https://openrouter.ai/models?supported_parameters=tools).

Popular models with tool support include:
- OpenAI GPT-4o and GPT-4o-mini
- Anthropic Claude 3.5 Sonnet
- Google Gemini models
- Many open-source models via various providers

#### Best Practices for Tool Calling

- **Clear Descriptions**: Provide detailed descriptions for tools and parameters
- **Error Handling**: Always validate tool arguments before execution
- **Tool Results**: Return structured, informative results from tools
- **Context Preservation**: Maintain full conversation history including tool calls
- **Streaming**: Handle tool calls appropriately when streaming responses
- **Testing**: Test tool interactions with different models as behavior may vary

### MCP Tool Conversion

The library provides utilities for converting MCP (Model Context Protocol) tool definitions to OpenRouter's OpenAI-compatible format. This enables seamless integration with MCP servers and clients.

#### Converting MCP Tools

```go
// Define MCP tools (as received from an MCP server)
mcpTools := []openrouter.MCPTool{
    {
        Name:        "read_file",
        Description: "Read the contents of a file from the filesystem",
        InputSchema: &openrouter.MCPInputSchema{
            Type: "object",
            Properties: map[string]interface{}{
                "path": map[string]interface{}{
                    "type":        "string",
                    "description": "The path to the file to read",
                },
            },
            Required: []string{"path"},
        },
    },
}

// Convert to OpenRouter format
tools := openrouter.ConvertMCPTools(mcpTools)

// Use with chat completion
response, err := client.ChatComplete(ctx, messages,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithTools(tools...),
)
```

#### Parsing MCP Tools from JSON

```go
// Parse MCP tools from JSON (as received from an MCP server)
mcpToolsJSON := []byte(`[
    {
        "name": "get_weather",
        "description": "Get the current weather for a location",
        "inputSchema": {
            "type": "object",
            "properties": {
                "location": {"type": "string", "description": "City name"}
            },
            "required": ["location"]
        }
    }
]`)

mcpTools, err := openrouter.ParseMCPToolsFromJSON(mcpToolsJSON)
if err != nil {
    log.Fatal(err)
}

// Convert and use
tools := openrouter.ConvertMCPTools(mcpTools)
```

#### Handling MCP Tool Results

```go
// Convert MCP tool result to string for tool response
mcpResult := openrouter.MCPToolResult{
    Content: []openrouter.MCPContent{
        {Type: "text", Text: "File contents: Hello, World!"},
    },
}

resultStr := openrouter.ConvertToolResultToMCP(mcpResult)

// Use in tool response message
messages = append(messages, openrouter.Message{
    Role:       "tool",
    Content:    resultStr,
    ToolCallID: toolCall.ID,
})
```

#### MCP Types

The library provides the following MCP types:

- **MCPTool**: Represents an MCP tool definition (name, description, inputSchema)
- **MCPInputSchema**: JSON Schema for tool parameters (type, properties, required)
- **MCPToolResult**: Result from tool execution (content array, isError)
- **MCPContent**: Content item in responses (type, text, data, mimeType)

#### Functions

- `ConvertMCPTool(mcpTool MCPTool) Tool` - Convert single MCP tool
- `ConvertMCPTools(mcpTools []MCPTool) []Tool` - Convert multiple tools
- `ConvertToolResultToMCP(result MCPToolResult) string` - Convert result to string
- `ParseMCPToolFromJSON(data []byte) (MCPTool, error)` - Parse single tool from JSON
- `ParseMCPToolsFromJSON(data []byte) ([]MCPTool, error)` - Parse multiple tools from JSON

See the [mcp-tools example](examples/mcp-tools) for complete usage examples.

### Web Search

The library supports OpenRouter's web search feature for augmenting model responses with real-time web data. Web search can be enabled using the `:online` model suffix or by configuring the web plugin.

#### Quick Start with :online Suffix

```go
// Simple web search using :online suffix
response, err := client.ChatComplete(ctx,
    openrouter.WithModel("openai/gpt-4o:online"),
    openrouter.WithMessages([]openrouter.Message{
        {Role: "user", Content: "What are the latest AI developments this week?"},
    }),
)
```

#### Using the Web Plugin

```go
// Configure web search with the plugin
webPlugin := openrouter.NewWebPlugin() // Uses defaults: auto engine, 5 results

response, err := client.ChatComplete(ctx,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithPlugins(webPlugin),
    openrouter.WithMessages(messages),
)

// Custom web plugin configuration
webPlugin := openrouter.NewWebPluginWithOptions(
    openrouter.WebSearchEngineExa,    // Force Exa search
    10,                                // Get 10 results
    "Recent web results for context:", // Custom prompt
)

response, err := client.ChatComplete(ctx,
    openrouter.WithModel("anthropic/claude-3.5-sonnet"),
    openrouter.WithPlugins(webPlugin),
    openrouter.WithMessages(messages),
)
```

#### Search Engine Options

- **Native**: Uses the provider's built-in web search (OpenAI, Anthropic)
- **Exa**: Uses Exa's neural search API (works with all models)
- **Auto** (default): Automatically selects the best available engine

```go
// Force native search for supported models
webPlugin := openrouter.Plugin{
    ID:     "web",
    Engine: string(openrouter.WebSearchEngineNative),
}

// Force Exa search for all models
webPlugin := openrouter.Plugin{
    ID:         "web",
    Engine:     string(openrouter.WebSearchEngineExa),
    MaxResults: 3,
}
```

#### Search Context Size (Native Only)

For models with native search support, control the search context depth:

```go
response, err := client.ChatComplete(ctx,
    openrouter.WithModel("openai/gpt-4o"),
    openrouter.WithPlugins(openrouter.NewWebPlugin()),
    openrouter.WithWebSearchOptions(&openrouter.WebSearchOptions{
        SearchContextSize: string(openrouter.WebSearchContextHigh), // low, medium, high
    }),
    openrouter.WithMessages(messages),
)
```

#### Parsing Search Annotations

Web search results are included in the response annotations:

```go
response, err := client.ChatComplete(ctx,
    openrouter.WithModel("openai/gpt-4o:online"),
    openrouter.WithMessages(messages),
)

// Extract URL citations from the response
citations := openrouter.ParseAnnotations(response.Choices[0].Message.Annotations)
for _, citation := range citations {
    fmt.Printf("Source: %s\n", citation.Title)
    fmt.Printf("URL: %s\n", citation.URL)
    fmt.Printf("Content: %s\n\n", citation.Content)
}
```

#### Pricing

- **Exa Search**: $4 per 1000 results (default 5 results = $0.02 per request)
- **Native Search (OpenAI)**:
  - GPT-4o models: $30-50 per 1000 requests depending on context size
  - GPT-4o-mini models: $25-30 per 1000 requests
- **Native Search (Perplexity)**:
  - Sonar models: $5-12 per 1000 requests
  - SonarPro models: $6-14 per 1000 requests

#### Best Practices

- Use `:online` suffix for simple cases with default settings
- Configure the web plugin for fine-grained control over search behavior
- Consider search costs when choosing between native and Exa engines
- Parse annotations to display sources and improve transparency
- Use higher search context for research tasks, lower for quick facts

### Responses API [BETA]

> ⚠️ **WARNING: BETA API - EXPECT BREAKING CHANGES**
>
> The Responses API is in beta and may have breaking changes at any time. Do not rely on this API for production workloads. The API structure, parameters, and behavior may change without notice as OpenRouter continues to develop this feature.
>
> For stable production use, consider using the Chat Completions API instead.

The library supports OpenRouter's new Responses API, which provides an OpenAI-compatible stateless API with enhanced capabilities including reasoning, tool calling, web search integration, and streaming.

#### Basic Usage

```go
// Simple string input
resp, err := client.CreateResponse(ctx, "What is 2+2?",
    openrouter.WithResponsesModel("openai/gpt-4o-mini"),
    openrouter.WithResponsesMaxOutputTokens(100),
)
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.GetTextContent())
```

#### Structured Input

```go
// Use structured input for multi-turn conversations
input := []openrouter.ResponsesInputItem{
    openrouter.CreateResponsesSystemMessage("You are a helpful assistant."),
    openrouter.CreateResponsesUserMessage("What is the capital of France?"),
}

resp, err := client.CreateResponse(ctx, input,
    openrouter.WithResponsesModel("openai/gpt-4o-mini"),
    openrouter.WithResponsesMaxOutputTokens(200),
)
```

#### Reasoning

```go
// Enable reasoning for complex problems
resp, err := client.CreateResponse(ctx, "Solve this step by step: 15 * 17",
    openrouter.WithResponsesModel("openai/o4-mini"),
    openrouter.WithResponsesMaxOutputTokens(500),
    openrouter.WithResponsesReasoningEffort(openrouter.ReasoningEffortMedium),
)

// Check for reasoning summary
if summary := resp.GetReasoningSummary(); len(summary) > 0 {
    for i, step := range summary {
        fmt.Printf("%d. %s\n", i+1, step)
    }
}
```

Reasoning effort levels: `ReasoningEffortMinimal`, `ReasoningEffortLow`, `ReasoningEffortMedium`, `ReasoningEffortHigh`

#### Tool Calling

```go
// Define tools using the flat ResponsesTool structure
// Note: Responses API uses a different tool format than Chat Completions API
weatherTool := openrouter.CreateResponsesTool(
    "get_weather",
    "Get weather for a location",
    map[string]any{
        "type": "object",
        "properties": map[string]any{
            "location": map[string]any{
                "type":        "string",
                "description": "The city and state, e.g. San Francisco, CA",
            },
        },
        "required": []string{"location"},
    },
)

resp, err := client.CreateResponse(ctx, "What's the weather in Tokyo?",
    openrouter.WithResponsesModel("openai/gpt-4o-mini"),
    openrouter.WithResponsesTools(weatherTool),
)

// Check for function calls
calls := resp.GetFunctionCalls()
if len(calls) > 0 {
    for _, call := range calls {
        fmt.Printf("Function: %s, Args: %s\n", call.Name, call.Arguments)
    }
}
```

#### Web Search

```go
// Enable web search for real-time information
resp, err := client.CreateResponse(ctx, "What are the latest AI news?",
    openrouter.WithResponsesModel("openai/gpt-4o-mini"),
    openrouter.WithResponsesWebSearch(3), // Get up to 3 search results
)

// Check for citations
annotations := resp.GetAnnotations()
for _, ann := range annotations {
    if ann.Type == "url_citation" {
        fmt.Printf("Source: %s\n", ann.URL)
    }
}
```

#### Streaming

```go
stream, err := client.CreateResponseStream(ctx, "Write a haiku about programming.",
    openrouter.WithResponsesModel("openai/gpt-4o-mini"),
    openrouter.WithResponsesMaxOutputTokens(100),
)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

var lastContent string
for event := range stream.Events() {
    content := event.GetTextContent()
    // Only print the new delta (GetTextContent returns cumulative content)
    if len(content) > len(lastContent) {
        fmt.Print(content[len(lastContent):])
        lastContent = content
    }
}

if err := stream.Err(); err != nil {
    log.Fatal(err)
}
```

#### Key Differences from Chat API

| Feature | Chat Completions | Responses API |
|---------|-----------------|---------------|
| Endpoint | `/chat/completions` | `/responses` |
| Input | Array of messages | String or structured array |
| Response | `choices` array | `output` array with typed items |
| Reasoning | Not available | Configurable effort levels |
| Web Search | Via `:online` suffix | Via `plugins` parameter |

See the [responses example](examples/responses) for complete usage examples.

## Examples

The `examples/` directory contains comprehensive examples:

- **basic/** - Simple usage examples for common tasks
- **streaming/** - Real-time streaming response handling
- **list-models/** - List and discover available models with filtering
- **model-endpoints/** - Inspect model endpoints with pricing and provider details
- **list-providers/** - List available providers with policy information
- **structured-output/** - JSON schema validation and structured responses
- **tool-calling/** - Complete tool/function calling examples with streaming
- **transforms/** - Message transforms for context window management
- **web_search/** - Web search plugin examples with various configurations
- **responses/** - **[BETA]** Responses API examples with reasoning, tools, and streaming
- **advanced/** - Advanced features like rate limiting and custom configuration

To run an example:

```bash
# Set your API key
export OPENROUTER_API_KEY="your-api-key"

# Run basic examples
go run examples/basic/main.go

# Run streaming examples
go run examples/streaming/main.go

# Run list models examples
go run examples/list-models/main.go

# Run model endpoints examples
go run examples/model-endpoints/main.go

# Run list providers examples
go run examples/list-providers/main.go

# Run advanced examples
go run examples/advanced/main.go

# Run structured output examples
go run examples/structured-output/main.go

# Run tool calling examples
go run examples/tool-calling/main.go

# Run streaming tool calling example
go run examples/tool-calling/streaming.go

# Run transforms examples
go run examples/transforms/main.go

# Run web search examples
go run examples/web_search/main.go

# Run responses API examples [BETA]
go run examples/responses/main.go
```

## Documentation

For detailed API documentation and usage examples, see [DOCUMENTATION.md](DOCUMENTATION.md).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
