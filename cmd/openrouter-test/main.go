package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Command-line flags
	var (
		apiKey    = flag.String("key", os.Getenv("OPENROUTER_API_KEY"), "OpenRouter API key (or set OPENROUTER_API_KEY env var)")
		model     = flag.String("model", "openai/gpt-3.5-turbo", "Model to use")
		test      = flag.String("test", "all", "Test to run: all, chat, stream, completion, error, provider, zdr, suffix, price, structured, tools, transforms, websearch")
		verbose   = flag.Bool("v", false, "Verbose output")
		timeout   = flag.Duration("timeout", 30*time.Second, "Request timeout")
		maxTokens = flag.Int("max-tokens", 100, "Maximum tokens for response")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "OpenRouter Go Client - Live API Test Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -key YOUR_KEY -test chat\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -test stream -model anthropic/claude-3-haiku\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  export OPENROUTER_API_KEY=YOUR_KEY && %s -test all\n", os.Args[0])
	}

	flag.Parse()

	if *apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: API key is required. Set via -key flag or OPENROUTER_API_KEY environment variable\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Create client
	client := openrouter.NewClient(
		openrouter.WithAPIKey(*apiKey),
		openrouter.WithTimeout(*timeout),
		openrouter.WithAppName("OpenRouter-Go-Test"),
		openrouter.WithRetry(3, time.Second),
	)

	fmt.Printf("🚀 OpenRouter Go Client - Live API Test\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Model: %s\n", *model)
	fmt.Printf("Test: %s\n", *test)
	fmt.Printf("Max Tokens: %d\n", *maxTokens)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	var success, failed int
	ctx := context.Background()

	// Run tests based on selection
	switch strings.ToLower(*test) {
	case "all":
		success, failed = runAllTests(ctx, client, *model, *maxTokens, *verbose)
	case "chat":
		if runChatTest(ctx, client, *model, *maxTokens, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "stream":
		if runStreamTest(ctx, client, *model, *maxTokens, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "completion":
		if runCompletionTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "error":
		if runErrorTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "provider":
		if runProviderRoutingTest(ctx, client, *model, *maxTokens, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "zdr":
		if runZDRTest(ctx, client, *model, *maxTokens, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "suffix":
		if runModelSuffixTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "price":
		if runPriceConstraintTest(ctx, client, *model, *maxTokens, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "structured":
		if runStructuredOutputTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "tools":
		if runToolCallingTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "transforms":
		if runTransformsTest(ctx, client, *model, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "websearch":
		if runWebSearchTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown test: %s\n", *test)
		flag.Usage()
		os.Exit(1)
	}

	// Print summary
	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📊 Test Summary\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("✅ Passed: %d\n", success)
	fmt.Printf("❌ Failed: %d\n", failed)

	if failed > 0 {
		os.Exit(1)
	}
	fmt.Printf("\n🎉 All tests passed!\n")
}

func runAllTests(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) (success, failed int) {
	tests := []struct {
		name string
		fn   func() bool
	}{
		{"Chat Completion", func() bool { return runChatTest(ctx, client, model, maxTokens, verbose) }},
		{"Streaming", func() bool { return runStreamTest(ctx, client, model, maxTokens, verbose) }},
		{"Legacy Completion", func() bool { return runCompletionTest(ctx, client, verbose) }},
		{"Error Handling", func() bool { return runErrorTest(ctx, client, verbose) }},
		{"Provider Routing", func() bool { return runProviderRoutingTest(ctx, client, model, maxTokens, verbose) }},
		{"ZDR", func() bool { return runZDRTest(ctx, client, model, maxTokens, verbose) }},
		{"Model Suffixes", func() bool { return runModelSuffixTest(ctx, client, verbose) }},
		{"Price Constraints", func() bool { return runPriceConstraintTest(ctx, client, model, maxTokens, verbose) }},
		{"Structured Output", func() bool { return runStructuredOutputTest(ctx, client, verbose) }},
		{"Tool Calling", func() bool { return runToolCallingTest(ctx, client, verbose) }},
		{"Message Transforms", func() bool { return runTransformsTest(ctx, client, model, verbose) }},
		{"Web Search", func() bool { return runWebSearchTest(ctx, client, verbose) }},
	}

	for _, test := range tests {
		if test.fn() {
			success++
		} else {
			failed++
		}
		fmt.Println()
	}

	return success, failed
}

func runChatTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Chat Completion\n")

	messages := []openrouter.Message{
		openrouter.CreateSystemMessage("You are a helpful assistant. Keep responses brief."),
		openrouter.CreateUserMessage("What is 2+2? Reply with just the number."),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithTemperature(0.7),
	)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return false
	}

	fmt.Printf("✅ Success! (%.2fs)\n", elapsed.Seconds())

	if verbose || true { // Always show some output
		fmt.Printf("   Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Tokens: %d prompt, %d completion, %d total\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	return true
}

func runStreamTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Streaming Chat\n")

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Count from 1 to 5, one number per line."),
	}

	start := time.Now()
	stream, err := client.ChatCompleteStream(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
	)
	if err != nil {
		fmt.Printf("❌ Failed to create stream: %v\n", err)
		return false
	}
	defer stream.Close()

	fmt.Printf("   Streaming: ")
	var fullResponse strings.Builder
	eventCount := 0

	for event := range stream.Events() {
		eventCount++
		for _, choice := range event.Choices {
			if choice.Delta != nil {
				if content, ok := choice.Delta.Content.(string); ok {
					fullResponse.WriteString(content)
					if verbose {
						fmt.Print(content)
					}
				}
			}
		}
	}

	elapsed := time.Since(start)

	if err := stream.Err(); err != nil {
		fmt.Printf("\n❌ Stream error: %v\n", err)
		return false
	}

	fmt.Printf("\n✅ Success! (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("   Events received: %d\n", eventCount)
	if !verbose {
		response := strings.TrimSpace(fullResponse.String())
		if len(response) > 100 {
			response = response[:100] + "..."
		}
		fmt.Printf("   Response: %s\n", response)
	}

	return true
}

func runCompletionTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Legacy Completion\n")

	// Only certain models support legacy completion
	completionModel := "openai/gpt-3.5-turbo-instruct"
	prompt := "The capital of France is"

	start := time.Now()
	resp, err := client.Complete(ctx, prompt,
		openrouter.WithCompletionModel(completionModel),
		openrouter.WithCompletionMaxTokens(10),
		openrouter.WithCompletionTemperature(0.5),
	)
	elapsed := time.Since(start)

	if err != nil {
		// Some accounts might not have access to instruct models
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 {
				fmt.Printf("⚠️  Skipped: Model %s not available\n", completionModel)
				return true // Don't fail the test
			}
		}
		fmt.Printf("❌ Failed: %v\n", err)
		return false
	}

	fmt.Printf("✅ Success! (%.2fs)\n", elapsed.Seconds())

	if verbose || true {
		fmt.Printf("   Prompt: %s\n", prompt)
		fmt.Printf("   Completion: %s\n", strings.TrimSpace(resp.Choices[0].Text))
		fmt.Printf("   Tokens: %d total\n", resp.Usage.TotalTokens)
	}

	return true
}

func runErrorTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Error Handling\n")

	// Test with invalid model to trigger error
	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Hello"),
	}

	_, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("invalid/nonexistent-model-xyz"),
		openrouter.WithMaxTokens(10),
	)

	if err == nil {
		fmt.Printf("❌ Expected error but got success\n")
		return false
	}

	// Check error handling
	if openrouter.IsRequestError(err) {
		reqErr := err.(*openrouter.RequestError)
		fmt.Printf("✅ Caught expected error\n")
		fmt.Printf("   Status: %d\n", reqErr.StatusCode)
		fmt.Printf("   Message: %s\n", reqErr.Message)

		if verbose {
			fmt.Printf("   Type: %s\n", reqErr.Type)
			fmt.Printf("   Is404: %v\n", reqErr.IsNotFoundError())
		}
		return true
	}

	fmt.Printf("❌ Unexpected error type: %T\n", err)
	return false
}

func runProviderRoutingTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Provider Routing\n")

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Say 'Hello from provider routing test' and nothing else."),
	}

	// Test 1: Provider order with fallbacks
	fmt.Printf("   Testing provider order...\n")
	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithProviderOrder("openai", "anthropic", "google"),
		openrouter.WithAllowFallbacks(true),
	)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("   ⚠️  Provider order test failed: %v\n", err)
		// Don't fail entirely as provider might not be available
	} else {
		fmt.Printf("   ✅ Provider order (%.2fs)\n", elapsed.Seconds())
		if verbose {
			fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		}
	}

	// Test 2: Ignore providers
	fmt.Printf("   Testing ignore providers...\n")
	start = time.Now()
	_, err = client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithIgnoreProviders("cohere"),
	)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ⚠️  Ignore providers test failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Ignore providers (%.2fs)\n", elapsed.Seconds())
	}

	// Test 3: Require parameters
	fmt.Printf("   Testing require parameters...\n")
	start = time.Now()
	_, err = client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithRequireParameters(true),
		openrouter.WithTemperature(0.5),
	)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ⚠️  Require parameters test failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Require parameters (%.2fs)\n", elapsed.Seconds())
	}

	fmt.Printf("✅ Provider routing tests completed\n")
	return true
}

func runZDRTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Zero Data Retention (ZDR)\n")

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What is 1+1? Reply with just the number."),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithZDR(true),
	)
	elapsed := time.Since(start)

	if err != nil {
		// ZDR might not be available for all models/providers
		fmt.Printf("⚠️  ZDR not available: %v\n", err)
		fmt.Printf("   (This is expected if the model/provider doesn't support ZDR)\n")
		return true // Don't fail the test
	}

	fmt.Printf("✅ Success with ZDR enabled! (%.2fs)\n", elapsed.Seconds())

	if verbose {
		fmt.Printf("   Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		fmt.Printf("   Model: %s\n", resp.Model)
	}

	return true
}

func runModelSuffixTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Model Suffixes (Nitro/Floor)\n")

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Reply with 'OK' and nothing else."),
	}

	// Test 1: Nitro suffix for throughput optimization
	fmt.Printf("   Testing :nitro suffix...\n")
	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("meta-llama/llama-3.1-8b-instruct:nitro"),
		openrouter.WithMaxTokens(10),
	)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("   ⚠️  Nitro test failed: %v\n", err)
		// Continue with other tests
	} else {
		fmt.Printf("   ✅ Nitro suffix (%.2fs)\n", elapsed.Seconds())
		if verbose {
			fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		}
	}

	// Test 2: Floor suffix for lowest price
	fmt.Printf("   Testing :floor suffix...\n")
	start = time.Now()
	resp, err = client.ChatComplete(ctx, messages,
		openrouter.WithModel("meta-llama/llama-3.1-8b-instruct:floor"),
		openrouter.WithMaxTokens(10),
	)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ⚠️  Floor test failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Floor suffix (%.2fs)\n", elapsed.Seconds())
		if verbose {
			fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		}
	}

	fmt.Printf("✅ Model suffix tests completed\n")
	return true
}

func runPriceConstraintTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Price Constraints\n")

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Say 'Price test OK' and nothing else."),
	}

	// Test with max price constraints
	maxPrice := openrouter.MaxPrice{
		Prompt:     5.0,  // $5 per million prompt tokens
		Completion: 10.0, // $10 per million completion tokens
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithMaxPrice(maxPrice),
		openrouter.WithProviderSort("throughput"), // Get fastest provider within price limit
	)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return false
	}

	fmt.Printf("✅ Success with price constraints! (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("   Max price: $%.2f/M prompt, $%.2f/M completion\n", maxPrice.Prompt, maxPrice.Completion)

	if verbose {
		fmt.Printf("   Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		fmt.Printf("   Model used: %s\n", resp.Model)
		fmt.Printf("   Tokens: %d prompt, %d completion\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens)
	}

	// Test with data collection policy
	fmt.Printf("   Testing data collection policy...\n")
	start = time.Now()
	_, err = client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithDataCollection("deny"),
	)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ⚠️  Data collection 'deny' not available: %v\n", err)
	} else {
		fmt.Printf("   ✅ Data collection policy (%.2fs)\n", elapsed.Seconds())
	}

	return true
}

func runStructuredOutputTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Structured Output\n")

	// Test 1: Basic structured output with weather data
	fmt.Printf("   Testing weather data schema...\n")
	weatherSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"city": map[string]interface{}{
				"type":        "string",
				"description": "City name",
			},
			"temperature": map[string]interface{}{
				"type":        "number",
				"description": "Temperature in Celsius",
			},
			"conditions": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"sunny", "cloudy", "rainy", "snowy"},
				"description": "Weather conditions",
			},
			"humidity": map[string]interface{}{
				"type":        "integer",
				"minimum":     0,
				"maximum":     100,
				"description": "Humidity percentage",
			},
		},
		"required":             []string{"city", "temperature", "conditions", "humidity"},
		"additionalProperties": false,
	}

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Describe the weather in Paris. Make up realistic values."),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o"), // Use a model that supports structured outputs
		openrouter.WithJSONSchema("weather", true, weatherSchema),
		openrouter.WithMaxTokens(100),
		openrouter.WithRequireParameters(true), // Only use providers that support structured outputs
	)
	elapsed := time.Since(start)

	if err != nil {
		// Check if it's because the model doesn't support structured outputs
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 400 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Model doesn't support structured outputs: %v\n", err.Error())
				fmt.Printf("   Testing with simpler JSON mode instead...\n")

				// Try with simple JSON mode
				resp, err = client.ChatComplete(ctx, messages,
					openrouter.WithModel("openai/gpt-4o"),
					openrouter.WithJSONMode(),
					openrouter.WithMaxTokens(100),
				)

				if err != nil {
					fmt.Printf("   ❌ JSON mode also failed: %v\n", err)
					return false
				}
			} else {
				fmt.Printf("   ❌ Failed: %v\n", err)
				return false
			}
		} else {
			fmt.Printf("   ❌ Failed: %v\n", err)
			return false
		}
	}

	// Parse and validate the JSON response
	var weatherData map[string]interface{}
	content := resp.Choices[0].Message.Content.(string)
	if err := json.Unmarshal([]byte(content), &weatherData); err != nil {
		fmt.Printf("   ❌ Failed to parse JSON: %v\n", err)
		fmt.Printf("   Response: %s\n", content)
		return false
	}

	fmt.Printf("   ✅ Weather schema (%.2fs)\n", elapsed.Seconds())
	if verbose {
		prettyJSON, _ := json.MarshalIndent(weatherData, "      ", "  ")
		fmt.Printf("      Response:\n%s\n", string(prettyJSON))
	} else {
		fmt.Printf("      City: %v, Temp: %v°C, Conditions: %v\n",
			weatherData["city"], weatherData["temperature"], weatherData["conditions"])
	}

	// Test 2: Structured output with streaming
	fmt.Printf("   Testing structured output with streaming...\n")
	taskSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"tasks": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "string",
						},
						"priority": map[string]interface{}{
							"type": "string",
							"enum": []string{"low", "medium", "high"},
						},
					},
					"required":               []string{"name", "priority"},
					"additionalProperties": false,
				},
				"minItems": 3,
				"maxItems": 3,
			},
		},
		"required":             []string{"tasks"},
		"additionalProperties": false,
	}

	streamMessages := []openrouter.Message{
		openrouter.CreateUserMessage("Create a list of exactly 3 programming tasks with priorities."),
	}

	start = time.Now()
	stream, err := client.ChatCompleteStream(ctx, streamMessages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithJSONSchema("tasks", true, taskSchema),
		openrouter.WithMaxTokens(150),
	)

	if err != nil {
		// Fallback to non-streaming if streaming with structured output isn't supported
		fmt.Printf("   ⚠️  Streaming with structured output error: %v\n", err)
	} else {
		defer stream.Close()

		var fullContent strings.Builder
		eventCount := 0

		for event := range stream.Events() {
			eventCount++
			for _, choice := range event.Choices {
				if choice.Delta != nil {
					if content, ok := choice.Delta.Content.(string); ok {
						fullContent.WriteString(content)
					}
				}
			}
		}

		elapsed = time.Since(start)

		if err := stream.Err(); err != nil {
			fmt.Printf("   ❌ Stream error: %v\n", err)
			return false
		}

		// Validate the streamed JSON
		var taskData map[string]interface{}
		if err := json.Unmarshal([]byte(fullContent.String()), &taskData); err != nil {
			fmt.Printf("   ⚠️  Streamed content not valid JSON (this can happen with some models)\n")
		} else {
			fmt.Printf("   ✅ Streaming with schema (%.2fs, %d events)\n", elapsed.Seconds(), eventCount)
			if verbose {
				prettyJSON, _ := json.MarshalIndent(taskData, "      ", "  ")
				fmt.Printf("      Tasks:\n%s\n", string(prettyJSON))
			}
		}
	}

	// Test 3: Simple JSON mode without strict schema
	fmt.Printf("   Testing simple JSON mode...\n")
	jsonMessages := []openrouter.Message{
		openrouter.CreateSystemMessage("You must always respond with valid JSON."),
		openrouter.CreateUserMessage("List 3 benefits of Go programming language as a JSON object with a 'benefits' array."),
	}

	start = time.Now()
	resp, err = client.ChatComplete(ctx, jsonMessages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithJSONMode(),
		openrouter.WithMaxTokens(150),
	)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ❌ JSON mode failed: %v\n", err)
		return false
	}

	// Validate it's valid JSON
	var jsonData map[string]interface{}
	content = resp.Choices[0].Message.Content.(string)
	if err := json.Unmarshal([]byte(content), &jsonData); err != nil {
		fmt.Printf("   ❌ Response is not valid JSON: %v\n", err)
		return false
	}

	fmt.Printf("   ✅ JSON mode (%.2fs)\n", elapsed.Seconds())
	if verbose {
		prettyJSON, _ := json.MarshalIndent(jsonData, "      ", "  ")
		fmt.Printf("      Response:\n%s\n", string(prettyJSON))
	}

	fmt.Printf("✅ Structured output tests completed\n")
	return true
}

func runToolCallingTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Tool/Function Calling\n")

	// Test 1: Basic tool calling
	fmt.Printf("   Testing basic tool calling...\n")

	// Define a simple calculator tool
	tools := []openrouter.Tool{
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "calculate",
				Description: "Perform basic mathematical calculations",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"operation": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"add", "subtract", "multiply", "divide"},
							"description": "Mathematical operation to perform",
						},
						"a": map[string]interface{}{
							"type":        "number",
							"description": "First operand",
						},
						"b": map[string]interface{}{
							"type":        "number",
							"description": "Second operand",
						},
					},
					"required": []string{"operation", "a", "b"},
				},
			},
		},
	}

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What is 15 multiplied by 7?"),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o-mini"), // Use a model that supports tools
		openrouter.WithTools(tools...),
		openrouter.WithMaxTokens(100),
	)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("   ❌ Failed initial request: %v\n", err)
		return false
	}

	// Check if the model requested tool calls
	if len(resp.Choices) == 0 || len(resp.Choices[0].Message.ToolCalls) == 0 {
		fmt.Printf("   ❌ Model didn't request any tool calls\n")
		return false
	}

	fmt.Printf("   ✅ Tool request received (%.2fs)\n", elapsed.Seconds())

	// Process tool calls
	toolCall := resp.Choices[0].Message.ToolCalls[0]
	fmt.Printf("   Tool: %s\n", toolCall.Function.Name)
	fmt.Printf("   Arguments: %s\n", toolCall.Function.Arguments)

	// Parse arguments and simulate tool execution
	var args struct {
		Operation string  `json:"operation"`
		A         float64 `json:"a"`
		B         float64 `json:"b"`
	}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		fmt.Printf("   ❌ Failed to parse arguments: %v\n", err)
		return false
	}

	// Execute the tool (simulated)
	var result float64
	switch args.Operation {
	case "add":
		result = args.A + args.B
	case "subtract":
		result = args.A - args.B
	case "multiply":
		result = args.A * args.B
	case "divide":
		if args.B != 0 {
			result = args.A / args.B
		} else {
			result = 0
		}
	}

	toolResult := fmt.Sprintf(`{"result": %f}`, result)
	fmt.Printf("   Tool result: %s\n", toolResult)

	// Add tool response to messages
	messages = append(messages, resp.Choices[0].Message)
	messages = append(messages, openrouter.Message{
		Role:       "tool",
		Content:    toolResult,
		ToolCallID: toolCall.ID,
	})

	// Get final response
	start = time.Now()
	finalResp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o-mini"),
		openrouter.WithTools(tools...),
		openrouter.WithMaxTokens(100),
	)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ❌ Failed final request: %v\n", err)
		return false
	}

	fmt.Printf("   ✅ Final response received (%.2fs)\n", elapsed.Seconds())
	if verbose {
		fmt.Printf("   Final answer: %s\n", strings.TrimSpace(finalResp.Choices[0].Message.Content.(string)))
	}

	// Test 2: Multiple tools and tool choice
	fmt.Printf("   Testing multiple tools and tool choice...\n")

	multiTools := []openrouter.Tool{
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "get_weather",
				Description: "Get current weather for a location",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"location"},
				},
			},
		},
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "get_time",
				Description: "Get current time for a timezone",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"timezone": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"timezone"},
				},
			},
		},
	}

	// Test with auto tool choice (default)
	weatherMessages := []openrouter.Message{
		openrouter.CreateUserMessage("What's the weather in Tokyo?"),
	}

	resp, err = client.ChatComplete(ctx, weatherMessages,
		openrouter.WithModel("openai/gpt-4o-mini"),
		openrouter.WithTools(multiTools...),
		openrouter.WithToolChoice("auto"),
		openrouter.WithMaxTokens(100),
	)

	if err != nil {
		fmt.Printf("   ⚠️  Multi-tool test failed: %v\n", err)
	} else if len(resp.Choices[0].Message.ToolCalls) > 0 {
		fmt.Printf("   ✅ Auto tool choice worked\n")
		fmt.Printf("      Tool selected: %s\n", resp.Choices[0].Message.ToolCalls[0].Function.Name)
	}

	// Test 3: Parallel tool calls
	fmt.Printf("   Testing parallel tool calls control...\n")

	parallelMessages := []openrouter.Message{
		openrouter.CreateUserMessage("What's the weather in Paris and the time in New York?"),
	}

	// Test with parallel tool calls disabled
	parallelCalls := false
	resp, err = client.ChatComplete(ctx, parallelMessages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithTools(multiTools...),
		openrouter.WithParallelToolCalls(&parallelCalls),
		openrouter.WithMaxTokens(100),
	)

	if err != nil {
		fmt.Printf("   ⚠️  Parallel tool calls test failed: %v\n", err)
	} else {
		toolCount := len(resp.Choices[0].Message.ToolCalls)
		fmt.Printf("   ✅ Parallel control test completed\n")
		fmt.Printf("      Tools requested: %d (parallel_tool_calls=false)\n", toolCount)
		if verbose && toolCount > 0 {
			for i, tc := range resp.Choices[0].Message.ToolCalls {
				fmt.Printf("      Tool %d: %s\n", i+1, tc.Function.Name)
			}
		}
	}

	// Test 4: Streaming with tool calls
	fmt.Printf("   Testing streaming with tool calls...\n")

	streamMessages := []openrouter.Message{
		openrouter.CreateUserMessage("Calculate 25 plus 17"),
	}

	stream, err := client.ChatCompleteStream(ctx, streamMessages,
		openrouter.WithModel("openai/gpt-4o-mini"),
		openrouter.WithTools(tools...),
		openrouter.WithMaxTokens(100),
	)

	if err != nil {
		fmt.Printf("   ❌ Failed to create stream: %v\n", err)
		return false
	}
	defer stream.Close()

	hasToolCalls := false
	eventCount := 0

	for event := range stream.Events() {
		eventCount++
		// Check for tool calls in the streaming response
		for _, choice := range event.Choices {
			// Check if delta contains tool calls
			if choice.Delta != nil && len(choice.Delta.ToolCalls) > 0 {
				hasToolCalls = true
			}
			// Check finish reason
			if choice.FinishReason == "tool_calls" {
				hasToolCalls = true
			}
		}
	}

	if err := stream.Err(); err != nil {
		fmt.Printf("   ❌ Stream error: %v\n", err)
		return false
	}

	if hasToolCalls {
		fmt.Printf("   ✅ Streaming with tool calls worked (%d events)\n", eventCount)
	} else {
		fmt.Printf("   ⚠️  No tool calls in stream (model may have answered directly)\n")
	}

	fmt.Printf("✅ Tool calling tests completed\n")
	return true
}

func runTransformsTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Message Transforms\n")

	// Test 1: Test with middle-out transform enabled for chat completion
	fmt.Printf("   Testing middle-out transform for chat...\n")

	// Create a long conversation that might exceed context for small models
	messages := []openrouter.Message{
		openrouter.CreateSystemMessage("You are a helpful assistant. Keep responses brief."),
		openrouter.CreateUserMessage("Hello, I'm learning about transforms."),
		openrouter.CreateAssistantMessage("Hello! I'm here to help you learn about transforms. What would you like to know?"),
		openrouter.CreateUserMessage("Tell me about the concept."),
		openrouter.CreateAssistantMessage("Transforms in the OpenRouter context help manage long conversations by compressing content that exceeds model context windows."),
		openrouter.CreateUserMessage("How does middle-out compression work?"),
		openrouter.CreateAssistantMessage("Middle-out compression removes or truncates messages from the middle of the conversation, keeping the beginning and end intact, since LLMs pay less attention to the middle of sequences."),
		openrouter.CreateUserMessage("Now, ignoring everything above, just tell me: what is 2+2? Reply with only the number."),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(10),
		openrouter.WithTransforms("middle-out"), // Enable middle-out compression
	)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("   ❌ Chat with transforms failed: %v\n", err)
		return false
	}

	fmt.Printf("   ✅ Chat with middle-out transform (%.2fs)\n", elapsed.Seconds())
	if verbose {
		fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		fmt.Printf("      Model: %s\n", resp.Model)
		fmt.Printf("      Tokens: %d prompt, %d completion\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens)
	}

	// Test 2: Test with transforms explicitly disabled (empty array)
	fmt.Printf("   Testing with transforms disabled...\n")

	shortMessages := []openrouter.Message{
		openrouter.CreateUserMessage("Say 'OK' and nothing else."),
	}

	start = time.Now()
	resp, err = client.ChatComplete(ctx, shortMessages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(10),
		openrouter.WithTransforms(), // Empty array disables transforms
	)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ❌ Chat without transforms failed: %v\n", err)
		return false
	}

	fmt.Printf("   ✅ Chat with transforms disabled (%.2fs)\n", elapsed.Seconds())
	if verbose {
		fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
	}

	// Test 3: Test transforms with legacy completion endpoint
	fmt.Printf("   Testing transforms with completion endpoint...\n")

	// Only test if the model supports completion endpoint
	completionModel := "openai/gpt-3.5-turbo-instruct"
	longPrompt := strings.Repeat("This is a test sentence. ", 50) + "\n\nNow just say 'Done' and nothing else:"

	start = time.Now()
	compResp, err := client.Complete(ctx, longPrompt,
		openrouter.WithCompletionModel(completionModel),
		openrouter.WithCompletionMaxTokens(10),
		openrouter.WithCompletionTransforms("middle-out"),
	)
	elapsed = time.Since(start)

	if err != nil {
		// Some accounts might not have access to instruct models
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Skipped: Model %s not available\n", completionModel)
			} else {
				fmt.Printf("   ⚠️  Completion with transforms failed: %v\n", err)
			}
		} else {
			fmt.Printf("   ⚠️  Completion with transforms failed: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ Completion with transforms (%.2fs)\n", elapsed.Seconds())
		if verbose {
			fmt.Printf("      Response: %s\n", strings.TrimSpace(compResp.Choices[0].Text))
			fmt.Printf("      Tokens: %d total\n", compResp.Usage.TotalTokens)
		}
	}

	// Test 4: Test with streaming and transforms
	fmt.Printf("   Testing transforms with streaming...\n")

	streamMessages := []openrouter.Message{
		openrouter.CreateUserMessage("Count from 1 to 3, one number per line."),
	}

	stream, err := client.ChatCompleteStream(ctx, streamMessages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(20),
		openrouter.WithTransforms("middle-out"),
	)

	if err != nil {
		fmt.Printf("   ❌ Failed to create stream with transforms: %v\n", err)
		return false
	}
	defer stream.Close()

	var fullResponse strings.Builder
	eventCount := 0

	for event := range stream.Events() {
		eventCount++
		for _, choice := range event.Choices {
			if choice.Delta != nil {
				if content, ok := choice.Delta.Content.(string); ok {
					fullResponse.WriteString(content)
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		fmt.Printf("   ❌ Stream error: %v\n", err)
		return false
	}

	fmt.Printf("   ✅ Streaming with transforms (%d events)\n", eventCount)
	if verbose {
		response := strings.TrimSpace(fullResponse.String())
		if len(response) > 100 {
			response = response[:100] + "..."
		}
		fmt.Printf("      Response: %s\n", response)
	}

	// Test 5: Test behavior with small context model (8K or less should default to middle-out)
	fmt.Printf("   Testing default behavior with small context model...\n")

	// Try to find a model with 8K context or less
	smallContextModel := "meta-llama/llama-3.1-8b-instruct"

	resp, err = client.ChatComplete(ctx, shortMessages,
		openrouter.WithModel(smallContextModel),
		openrouter.WithMaxTokens(10),
		// Not specifying transforms - should default to middle-out for 8K models
	)

	if err != nil {
		// Model might not be available
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Small context model test skipped: %s not available\n", smallContextModel)
			} else {
				fmt.Printf("   ⚠️  Small context model test failed: %v\n", err)
			}
		} else {
			fmt.Printf("   ⚠️  Small context model test failed: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ Default transform behavior tested\n")
		if verbose {
			fmt.Printf("      Model: %s (defaults to middle-out if ≤8K context)\n", smallContextModel)
			fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		}
	}

	fmt.Printf("✅ Message transforms tests completed\n")
	return true
}

func runWebSearchTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Web Search\n")

	// Test 1: Using :online suffix
	fmt.Printf("   Testing :online model suffix...\n")

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What major tech company had the biggest news today? Reply in one sentence."),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o-mini:online"),
		openrouter.WithMaxTokens(100),
	)
	elapsed := time.Since(start)

	if err != nil {
		// Web search might not be available or model might not support it
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			fmt.Printf("   ⚠️  :online suffix not available: %v\n", reqErr.Message)
		} else {
			fmt.Printf("   ❌ Failed: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ :online suffix worked (%.2fs)\n", elapsed.Seconds())
		if verbose {
			fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		}

		// Check for annotations (web search citations)
		if len(resp.Choices[0].Message.Annotations) > 0 {
			fmt.Printf("      Found %d web citations\n", len(resp.Choices[0].Message.Annotations))
			citations := openrouter.ParseAnnotations(resp.Choices[0].Message.Annotations)
			for i, citation := range citations[:min(3, len(citations))] {
				fmt.Printf("      Citation %d: %s\n", i+1, citation.Title)
				if verbose {
					fmt.Printf("         URL: %s\n", citation.URL)
				}
			}
		}
	}

	// Test 2: Using web plugin with defaults
	fmt.Printf("   Testing web plugin with defaults...\n")

	webPlugin := openrouter.NewWebPlugin()

	start = time.Now()
	resp, err = client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithPlugins(webPlugin),
		openrouter.WithMaxTokens(100),
	)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ⚠️  Web plugin test failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Web plugin with defaults (%.2fs)\n", elapsed.Seconds())
		if verbose {
			response := strings.TrimSpace(resp.Choices[0].Message.Content.(string))
			if len(response) > 200 {
				response = response[:200] + "..."
			}
			fmt.Printf("      Response: %s\n", response)
		}
	}

	// Test 3: Custom web plugin configuration
	fmt.Printf("   Testing custom web plugin...\n")

	customPlugin := openrouter.NewWebPluginWithOptions(
		openrouter.WebSearchEngineAuto,
		3, // Only 3 results
		"Here are some relevant search results:",
	)

	weatherMessages := []openrouter.Message{
		openrouter.CreateUserMessage("What's the current temperature in Tokyo? Reply with just the temperature."),
	}

	start = time.Now()
	resp, err = client.ChatComplete(ctx, weatherMessages,
		openrouter.WithModel("anthropic/claude-3.5-sonnet"),
		openrouter.WithPlugins(customPlugin),
		openrouter.WithMaxTokens(50),
	)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ⚠️  Custom web plugin test failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Custom web plugin (%.2fs, max 3 results)\n", elapsed.Seconds())
		if verbose {
			fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		}
	}

	// Test 4: Force specific engine (Exa)
	fmt.Printf("   Testing forced Exa engine...\n")

	exaPlugin := openrouter.Plugin{
		ID:         "web",
		Engine:     string(openrouter.WebSearchEngineExa),
		MaxResults: 2,
	}

	techMessages := []openrouter.Message{
		openrouter.CreateUserMessage("Find recent research papers about transformer models. List just the titles."),
	}

	start = time.Now()
	resp, err = client.ChatComplete(ctx, techMessages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithPlugins(exaPlugin),
		openrouter.WithMaxTokens(150),
	)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ⚠️  Exa engine test failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Forced Exa engine (%.2fs)\n", elapsed.Seconds())
		if verbose {
			response := strings.TrimSpace(resp.Choices[0].Message.Content.(string))
			if len(response) > 300 {
				response = response[:300] + "..."
			}
			fmt.Printf("      Response: %s\n", response)
		}
	}

	// Test 5: Native search with context size (if supported)
	fmt.Printf("   Testing native search with context size...\n")

	nativePlugin := openrouter.Plugin{
		ID:     "web",
		Engine: string(openrouter.WebSearchEngineNative),
	}

	start = time.Now()
	resp, err = client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithPlugins(nativePlugin),
		openrouter.WithWebSearchOptions(&openrouter.WebSearchOptions{
			SearchContextSize: string(openrouter.WebSearchContextMedium),
		}),
		openrouter.WithMaxTokens(100),
	)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ⚠️  Native search test failed: %v\n", err)
		fmt.Printf("      (This is expected if the model doesn't support native search)\n")
	} else {
		fmt.Printf("   ✅ Native search with medium context (%.2fs)\n", elapsed.Seconds())
		if verbose {
			fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		}
	}

	// Test 6: Web search with streaming
	fmt.Printf("   Testing web search with streaming...\n")

	streamMessages := []openrouter.Message{
		openrouter.CreateUserMessage("What's the latest AI news today? Give me 3 bullet points."),
	}

	stream, err := client.ChatCompleteStream(ctx, streamMessages,
		openrouter.WithModel("openai/gpt-4o-mini:online"),
		openrouter.WithMaxTokens(150),
	)

	if err != nil {
		fmt.Printf("   ⚠️  Streaming with web search failed: %v\n", err)
	} else {
		defer stream.Close()

		var fullResponse strings.Builder
		eventCount := 0
		hasAnnotations := false

		for event := range stream.Events() {
			eventCount++
			for _, choice := range event.Choices {
				if choice.Delta != nil {
					if content, ok := choice.Delta.Content.(string); ok {
						fullResponse.WriteString(content)
					}
					// Check for annotations in delta
					if len(choice.Delta.Annotations) > 0 {
						hasAnnotations = true
					}
				}
			}
		}

		if err := stream.Err(); err != nil {
			fmt.Printf("   ❌ Stream error: %v\n", err)
		} else {
			fmt.Printf("   ✅ Streaming with web search (%d events)\n", eventCount)
			if hasAnnotations {
				fmt.Printf("      Web citations found in stream\n")
			}
			if verbose {
				response := strings.TrimSpace(fullResponse.String())
				if len(response) > 200 {
					response = response[:200] + "..."
				}
				fmt.Printf("      Response: %s\n", response)
			}
		}
	}

	// Test 7: Test helper function for online model
	fmt.Printf("   Testing WithOnlineModel helper...\n")

	onlineModel := openrouter.WithOnlineModel("meta-llama/llama-3.1-8b-instruct")
	if onlineModel == "meta-llama/llama-3.1-8b-instruct:online" {
		fmt.Printf("   ✅ WithOnlineModel helper works correctly\n")
	} else {
		fmt.Printf("   ❌ WithOnlineModel helper failed\n")
		return false
	}

	fmt.Printf("✅ Web search tests completed\n")
	return true
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}