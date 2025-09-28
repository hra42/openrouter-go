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
- ✅ Structured outputs with JSON schema validation
- ✅ Tool/Function calling support with streaming
- ✅ Message transforms for automatic context window management
- ✅ Web Search plugin for real-time web data integration

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
│   ├── basic/             # Basic usage examples
│   ├── streaming/         # Streaming examples
│   ├── structured-output/ # Structured outputs with JSON schema
│   ├── tool-calling/      # Tool/function calling examples
│   ├── web_search/        # Web search plugin examples
│   └── advanced/          # Advanced configuration examples
└── internal/
    └── sse/           # Internal SSE parser implementation
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

### CI/CD with Jenkins

The project includes a Jenkinsfile for continuous integration. The pipeline:
- Runs all unit tests
- Performs coverage analysis
- Executes race condition detection
- Runs integration tests with the OpenRouter API

To use the Jenkins pipeline, ensure you have configured the `openrouter-api-key` credential in your Jenkins instance.

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

## Examples

The `examples/` directory contains comprehensive examples:

- **basic/** - Simple usage examples for common tasks
- **streaming/** - Real-time streaming response handling
- **structured-output/** - JSON schema validation and structured responses
- **tool-calling/** - Complete tool/function calling examples with streaming
- **transforms/** - Message transforms for context window management
- **web_search/** - Web search plugin examples with various configurations
- **advanced/** - Advanced features like rate limiting and custom configuration

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
```

## Documentation

For detailed API documentation and usage examples, see [DOCUMENTATION.md](DOCUMENTATION.md).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
