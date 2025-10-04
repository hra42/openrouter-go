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
		test      = flag.String("test", "all", "Test to run: all, chat, stream, completion, error, provider, zdr, suffix, price, structured, tools, transforms, websearch, models, endpoints, providers, credits")
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

	// Create client with app attribution
	client := openrouter.NewClient(
		openrouter.WithAPIKey(*apiKey),
		openrouter.WithTimeout(*timeout),
		openrouter.WithReferer("https://github.com/hra42/openrouter-go"),
		openrouter.WithAppName("OpenRouter-Go Test Suite"),
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
		if runCompletionTest(ctx, client, *model, *verbose) {
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
		if runModelSuffixTest(ctx, client, *model, *verbose) {
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
		if runStructuredOutputTest(ctx, client, *model, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "tools":
		if runToolCallingTest(ctx, client, *model, *verbose) {
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
		if runWebSearchTest(ctx, client, *model, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "models":
		if runModelsTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "endpoints":
		if runModelEndpointsTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "providers":
		if runProvidersTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "credits":
		if runCreditsTest(ctx, client, *verbose) {
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
		{"Legacy Completion", func() bool { return runCompletionTest(ctx, client, model, verbose) }},
		{"Error Handling", func() bool { return runErrorTest(ctx, client, verbose) }},
		{"Provider Routing", func() bool { return runProviderRoutingTest(ctx, client, model, maxTokens, verbose) }},
		{"ZDR", func() bool { return runZDRTest(ctx, client, model, maxTokens, verbose) }},
		{"Model Suffixes", func() bool { return runModelSuffixTest(ctx, client, model, verbose) }},
		{"Price Constraints", func() bool { return runPriceConstraintTest(ctx, client, model, maxTokens, verbose) }},
		{"Structured Output", func() bool { return runStructuredOutputTest(ctx, client, model, verbose) }},
		{"Tool Calling", func() bool { return runToolCallingTest(ctx, client, model, verbose) }},
		{"Message Transforms", func() bool { return runTransformsTest(ctx, client, model, verbose) }},
		{"Web Search", func() bool { return runWebSearchTest(ctx, client, model, verbose) }},
		{"List Models", func() bool { return runModelsTest(ctx, client, verbose) }},
		{"Model Endpoints", func() bool { return runModelEndpointsTest(ctx, client, verbose) }},
		{"List Providers", func() bool { return runProvidersTest(ctx, client, verbose) }},
		{"Get Credits", func() bool { return runCreditsTest(ctx, client, verbose) }},
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

func runCompletionTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Legacy Completion\n")

	prompt := "The capital of France is"

	start := time.Now()
	resp, err := client.Complete(ctx, prompt,
		openrouter.WithCompletionModel(model),
		openrouter.WithCompletionMaxTokens(10),
		openrouter.WithCompletionTemperature(0.5),
	)
	elapsed := time.Since(start)

	if err != nil {
		// Some models might not support legacy completion format
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support legacy completion\n", model)
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

func runModelSuffixTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Model Suffixes (Nitro/Floor)\n")

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Reply with 'OK' and nothing else."),
	}

	// Test 1: Nitro suffix for throughput optimization
	fmt.Printf("   Testing :nitro suffix...\n")
	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model+":nitro"),
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
		openrouter.WithModel(model+":floor"),
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

func runStructuredOutputTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
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
		openrouter.WithModel(model),
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
					openrouter.WithModel(model),
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
					"required":             []string{"name", "priority"},
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
		openrouter.WithModel(model),
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
		openrouter.WithModel(model),
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

func runToolCallingTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
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
		openrouter.WithModel(model),
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
		openrouter.WithModel(model),
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
		openrouter.WithModel(model),
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
		openrouter.WithModel(model),
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
		openrouter.WithModel(model),
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

	// Test with legacy completion endpoint
	longPrompt := strings.Repeat("This is a test sentence. ", 50) + "\n\nNow just say 'Done' and nothing else:"

	start = time.Now()
	compResp, err := client.Complete(ctx, longPrompt,
		openrouter.WithCompletionModel(model),
		openrouter.WithCompletionMaxTokens(10),
		openrouter.WithCompletionTransforms("middle-out"),
	)
	elapsed = time.Since(start)

	if err != nil {
		// Some models might not support legacy completion format
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Skipped: Model %s not available or doesn't support legacy completion\n", model)
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

	// Test 5: Test behavior with the provided model
	fmt.Printf("   Testing default behavior with provided model...\n")

	resp, err = client.ChatComplete(ctx, shortMessages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(10),
		// Not specifying transforms
	)

	if err != nil {
		// Model might not be available
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Model test skipped: %s not available\n", model)
			} else {
				fmt.Printf("   ⚠️  Model test failed: %v\n", err)
			}
		} else {
			fmt.Printf("   ⚠️  Model test failed: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ Default transform behavior tested\n")
		if verbose {
			fmt.Printf("      Model: %s\n", model)
			fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		}
	}

	fmt.Printf("✅ Message transforms tests completed\n")
	return true
}

func runWebSearchTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Web Search\n")

	// Test 1: Using :online suffix
	fmt.Printf("   Testing :online model suffix...\n")

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What major tech company had the biggest news today? Reply in one sentence."),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model+":online"),
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
		openrouter.WithModel(model),
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
		openrouter.WithModel(model),
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
		openrouter.WithModel(model),
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
		openrouter.WithModel(model),
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
		openrouter.WithModel(model+":online"),
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

	onlineModel := openrouter.WithOnlineModel(model)
	expectedOnline := model + ":online"
	if onlineModel == expectedOnline {
		fmt.Printf("   ✅ WithOnlineModel helper works correctly\n")
	} else {
		fmt.Printf("   ❌ WithOnlineModel helper failed\n")
		return false
	}

	fmt.Printf("✅ Web search tests completed\n")
	return true
}

func runModelsTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List Models\n")

	// Test 1: List all models
	fmt.Printf("   Testing list all models...\n")
	start := time.Now()
	resp, err := client.ListModels(ctx, nil)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Failed to list models: %v\n", err)
		return false
	}

	fmt.Printf("   ✅ Listed all models (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total models: %d\n", len(resp.Data))

	if len(resp.Data) == 0 {
		fmt.Printf("   ❌ No models returned\n")
		return false
	}

	// Display first few models
	if verbose {
		fmt.Printf("\n   First 5 models:\n")
		for i, model := range resp.Data {
			if i >= 5 {
				break
			}
			fmt.Printf("      %d. %s (%s)\n", i+1, model.Name, model.ID)
			if model.ContextLength != nil {
				fmt.Printf("         Context: %.0f tokens\n", *model.ContextLength)
			}
			fmt.Printf("         Pricing: $%s/M prompt, $%s/M completion\n",
				model.Pricing.Prompt, model.Pricing.Completion)
		}
	} else {
		// Show just a couple in non-verbose mode
		for i, model := range resp.Data {
			if i >= 2 {
				break
			}
			fmt.Printf("      Example: %s (%s)\n", model.Name, model.ID)
		}
	}

	// Test 2: Validate model structure
	fmt.Printf("\n   Validating model data structure...\n")
	firstModel := resp.Data[0]

	// Check required fields
	if firstModel.ID == "" {
		fmt.Printf("   ❌ Model missing ID\n")
		return false
	}
	if firstModel.Name == "" {
		fmt.Printf("   ❌ Model missing Name\n")
		return false
	}
	if firstModel.Description == "" {
		fmt.Printf("   ❌ Model missing Description\n")
		return false
	}

	// Check architecture
	if len(firstModel.Architecture.InputModalities) == 0 {
		fmt.Printf("   ❌ Model missing InputModalities\n")
		return false
	}
	if len(firstModel.Architecture.OutputModalities) == 0 {
		fmt.Printf("   ❌ Model missing OutputModalities\n")
		return false
	}
	if firstModel.Architecture.Tokenizer == "" {
		fmt.Printf("   ❌ Model missing Tokenizer\n")
		return false
	}

	// Check pricing
	if firstModel.Pricing.Prompt == "" {
		fmt.Printf("   ❌ Model missing Prompt pricing\n")
		return false
	}
	if firstModel.Pricing.Completion == "" {
		fmt.Printf("   ❌ Model missing Completion pricing\n")
		return false
	}

	fmt.Printf("   ✅ Model structure validation passed\n")

	if verbose {
		fmt.Printf("\n   First model details:\n")
		fmt.Printf("      ID: %s\n", firstModel.ID)
		fmt.Printf("      Name: %s\n", firstModel.Name)
		fmt.Printf("      Description: %s\n", truncateString(firstModel.Description, 80))
		if firstModel.ContextLength != nil {
			fmt.Printf("      Context Length: %.0f tokens\n", *firstModel.ContextLength)
		}
		fmt.Printf("      Input Modalities: %v\n", firstModel.Architecture.InputModalities)
		fmt.Printf("      Output Modalities: %v\n", firstModel.Architecture.OutputModalities)
		fmt.Printf("      Tokenizer: %s\n", firstModel.Architecture.Tokenizer)
		if firstModel.Architecture.InstructType != nil {
			fmt.Printf("      Instruct Type: %s\n", *firstModel.Architecture.InstructType)
		}
		fmt.Printf("      Is Moderated: %v\n", firstModel.TopProvider.IsModerated)
		if len(firstModel.SupportedParameters) > 0 {
			fmt.Printf("      Supported Parameters: %v\n", firstModel.SupportedParameters)
		}
	}

	// Test 3: Filter by category
	fmt.Printf("\n   Testing category filter (programming)...\n")
	start = time.Now()
	categoryResp, err := client.ListModels(ctx, &openrouter.ListModelsOptions{
		Category: "programming",
	})
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ❌ Failed to list models by category: %v\n", err)
		return false
	}

	fmt.Printf("   ✅ Listed programming models (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Programming models: %d\n", len(categoryResp.Data))

	if len(categoryResp.Data) == 0 {
		fmt.Printf("   ⚠️  No programming models found (this might be expected)\n")
	} else if verbose {
		fmt.Printf("\n   Top 3 programming models:\n")
		for i, model := range categoryResp.Data {
			if i >= 3 {
				break
			}
			fmt.Printf("      %d. %s (%s)\n", i+1, model.Name, model.ID)
		}
	}

	// Test 4: Check for specific well-known models
	fmt.Printf("\n   Checking for well-known models...\n")
	wellKnownModels := []string{
		"openai/gpt-4o",
		"anthropic/claude-3.5-sonnet",
		"google/gemini-pro",
		"meta-llama/llama-3.1-8b-instruct",
	}

	foundModels := make(map[string]bool)
	for _, model := range resp.Data {
		for _, knownModel := range wellKnownModels {
			if model.ID == knownModel {
				foundModels[knownModel] = true
			}
		}
	}

	foundCount := len(foundModels)
	fmt.Printf("   Found %d/%d well-known models\n", foundCount, len(wellKnownModels))

	if verbose {
		for _, knownModel := range wellKnownModels {
			status := "❌"
			if foundModels[knownModel] {
				status = "✅"
			}
			fmt.Printf("      %s %s\n", status, knownModel)
		}
	}

	// Test 5: Verify pricing information
	fmt.Printf("\n   Validating pricing information...\n")
	hasPricingInfo := 0
	for _, model := range resp.Data {
		if model.Pricing.Prompt != "" && model.Pricing.Completion != "" {
			hasPricingInfo++
		}
	}

	pricingPercent := (float64(hasPricingInfo) / float64(len(resp.Data))) * 100
	fmt.Printf("   %.1f%% of models have pricing info (%d/%d)\n",
		pricingPercent, hasPricingInfo, len(resp.Data))

	if pricingPercent < 90 {
		fmt.Printf("   ⚠️  Warning: Less than 90%% of models have pricing info\n")
	} else {
		fmt.Printf("   ✅ Pricing validation passed\n")
	}

	fmt.Printf("\n✅ List models tests completed\n")
	return true
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func runModelEndpointsTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List Model Endpoints\n")

	// Test 1: List endpoints for GPT-4
	fmt.Printf("   Testing endpoints for GPT-4...\n")
	start := time.Now()
	resp, err := client.ListModelEndpoints(ctx, "openai", "gpt-4")
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Failed to list model endpoints: %v\n", err)
		return false
	}

	fmt.Printf("   ✅ Listed GPT-4 endpoints (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Model: %s (%s)\n", resp.Data.Name, resp.Data.ID)
	fmt.Printf("      Total endpoints: %d\n", len(resp.Data.Endpoints))

	if len(resp.Data.Endpoints) == 0 {
		fmt.Printf("   ❌ No endpoints returned\n")
		return false
	}

	// Display endpoint details
	if verbose {
		fmt.Printf("\n   Model Details:\n")
		fmt.Printf("      Description: %s\n", truncateString(resp.Data.Description, 100))
		if resp.Data.Architecture.Tokenizer != nil {
			fmt.Printf("      Tokenizer: %s\n", *resp.Data.Architecture.Tokenizer)
		}
		if resp.Data.Architecture.InstructType != nil {
			fmt.Printf("      Instruct Type: %s\n", *resp.Data.Architecture.InstructType)
		}
		fmt.Printf("      Input Modalities: %v\n", resp.Data.Architecture.InputModalities)
		fmt.Printf("      Output Modalities: %v\n", resp.Data.Architecture.OutputModalities)

		fmt.Printf("\n   First 3 endpoints:\n")
		for i, endpoint := range resp.Data.Endpoints {
			if i >= 3 {
				break
			}
			fmt.Printf("      Endpoint %d:\n", i+1)
			fmt.Printf("         Provider: %s\n", endpoint.ProviderName)
			fmt.Printf("         Name: %s\n", endpoint.Name)
			fmt.Printf("         Status: %.0f\n", endpoint.Status)
			fmt.Printf("         Context Length: %.0f tokens\n", endpoint.ContextLength)
			if endpoint.MaxCompletionTokens != nil {
				fmt.Printf("         Max Completion Tokens: %.0f\n", *endpoint.MaxCompletionTokens)
			}
			if endpoint.Quantization != nil && *endpoint.Quantization != "" {
				fmt.Printf("         Quantization: %s\n", *endpoint.Quantization)
			}
			fmt.Printf("         Pricing - Prompt: $%s/M, Completion: $%s/M\n",
				endpoint.Pricing.Prompt, endpoint.Pricing.Completion)
			if endpoint.UptimeLast30m != nil {
				fmt.Printf("         Uptime (30m): %.2f%%\n", *endpoint.UptimeLast30m*100)
			}
			if len(endpoint.SupportedParameters) > 0 {
				fmt.Printf("         Supported Parameters: %d\n", len(endpoint.SupportedParameters))
			}
		}
	} else {
		// Non-verbose: just show a sample
		endpoint := resp.Data.Endpoints[0]
		fmt.Printf("      Example endpoint: %s (Provider: %s)\n",
			endpoint.Name, endpoint.ProviderName)
		fmt.Printf("      Pricing: $%s/M prompt, $%s/M completion\n",
			endpoint.Pricing.Prompt, endpoint.Pricing.Completion)
	}

	// Test 2: Validate endpoint structure
	fmt.Printf("\n   Validating endpoint data structure...\n")
	firstEndpoint := resp.Data.Endpoints[0]

	// Check required fields
	if firstEndpoint.Name == "" {
		fmt.Printf("   ❌ Endpoint missing Name\n")
		return false
	}
	if firstEndpoint.ProviderName == "" {
		fmt.Printf("   ❌ Endpoint missing ProviderName\n")
		return false
	}
	if firstEndpoint.ContextLength == 0 {
		fmt.Printf("   ❌ Endpoint missing ContextLength\n")
		return false
	}
	// Status can be 0, 1, or other numeric values, so we don't validate it's non-zero

	// Check pricing
	if firstEndpoint.Pricing.Prompt == "" {
		fmt.Printf("   ❌ Endpoint missing Prompt pricing\n")
		return false
	}
	if firstEndpoint.Pricing.Completion == "" {
		fmt.Printf("   ❌ Endpoint missing Completion pricing\n")
		return false
	}

	fmt.Printf("   ✅ Endpoint structure validation passed\n")

	// Test 3: List endpoints for Claude
	fmt.Printf("\n   Testing endpoints for Claude-3.5 Sonnet...\n")
	start = time.Now()
	claudeResp, err := client.ListModelEndpoints(ctx, "anthropic", "claude-3.5-sonnet")
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("   ❌ Failed to list Claude endpoints: %v\n", err)
		return false
	}

	fmt.Printf("   ✅ Listed Claude endpoints (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Model: %s\n", claudeResp.Data.Name)
	fmt.Printf("      Total endpoints: %d\n", len(claudeResp.Data.Endpoints))

	if verbose && len(claudeResp.Data.Endpoints) > 0 {
		fmt.Printf("      Example endpoint: %s\n", claudeResp.Data.Endpoints[0].ProviderName)
	}

	// Test 4: Test with invalid model (error handling)
	fmt.Printf("\n   Testing error handling with invalid model...\n")
	_, err = client.ListModelEndpoints(ctx, "invalid", "nonexistent-model")

	if err == nil {
		fmt.Printf("   ⚠️  Expected error but got success (model might exist)\n")
	} else {
		fmt.Printf("   ✅ Error handling works correctly\n")
		if verbose {
			fmt.Printf("      Error: %v\n", err)
		}
	}

	// Test 5: Test with empty parameters
	fmt.Printf("\n   Testing validation with empty parameters...\n")

	_, err = client.ListModelEndpoints(ctx, "", "gpt-4")
	if err == nil {
		fmt.Printf("   ❌ Should have errored with empty author\n")
		return false
	}
	fmt.Printf("   ✅ Empty author validation passed\n")

	_, err = client.ListModelEndpoints(ctx, "openai", "")
	if err == nil {
		fmt.Printf("   ❌ Should have errored with empty slug\n")
		return false
	}
	fmt.Printf("   ✅ Empty slug validation passed\n")

	// Test 6: Compare pricing across endpoints
	if verbose && len(resp.Data.Endpoints) > 1 {
		fmt.Printf("\n   Pricing comparison for %s:\n", resp.Data.Name)
		fmt.Printf("   %-30s %-15s %-15s\n", "Provider", "Prompt/M", "Completion/M")
		fmt.Printf("   %s\n", strings.Repeat("-", 60))
		for i, endpoint := range resp.Data.Endpoints {
			if i >= 5 {
				fmt.Printf("   ... and %d more endpoints\n", len(resp.Data.Endpoints)-5)
				break
			}
			fmt.Printf("   %-30s $%-14s $%-14s\n",
				endpoint.ProviderName,
				endpoint.Pricing.Prompt,
				endpoint.Pricing.Completion,
			)
		}
	}

	fmt.Printf("\n✅ Model endpoints tests completed\n")
	return true
}

func runProvidersTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List Providers\n")

	// Test: List all providers
	fmt.Printf("   Testing list all providers...\n")
	start := time.Now()
	resp, err := client.ListProviders(ctx)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Failed to list providers: %v\n", err)
		return false
	}

	fmt.Printf("   ✅ Listed all providers (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total providers: %d\n", len(resp.Data))

	if len(resp.Data) == 0 {
		fmt.Printf("   ❌ No providers returned\n")
		return false
	}

	// Display first few providers
	if verbose {
		fmt.Printf("\n   First 5 providers:\n")
		for i, provider := range resp.Data {
			if i >= 5 {
				break
			}
			fmt.Printf("      %d. %s (%s)\n", i+1, provider.Name, provider.Slug)
			if provider.PrivacyPolicyURL != nil {
				fmt.Printf("         Privacy Policy: %s\n", *provider.PrivacyPolicyURL)
			}
			if provider.TermsOfServiceURL != nil {
				fmt.Printf("         Terms of Service: %s\n", *provider.TermsOfServiceURL)
			}
			if provider.StatusPageURL != nil {
				fmt.Printf("         Status Page: %s\n", *provider.StatusPageURL)
			}
		}
	} else {
		// Show just a couple in non-verbose mode
		for i, provider := range resp.Data {
			if i >= 3 {
				break
			}
			fmt.Printf("      Example: %s (%s)\n", provider.Name, provider.Slug)
		}
	}

	// Validate provider structure
	fmt.Printf("\n   Validating provider data structure...\n")
	firstProvider := resp.Data[0]

	// Check required fields
	if firstProvider.Name == "" {
		fmt.Printf("   ❌ Provider missing Name\n")
		return false
	}
	if firstProvider.Slug == "" {
		fmt.Printf("   ❌ Provider missing Slug\n")
		return false
	}

	fmt.Printf("   ✅ Provider structure validation passed\n")

	if verbose {
		fmt.Printf("\n   First provider details:\n")
		fmt.Printf("      Name: %s\n", firstProvider.Name)
		fmt.Printf("      Slug: %s\n", firstProvider.Slug)
		if firstProvider.PrivacyPolicyURL != nil {
			fmt.Printf("      Privacy Policy URL: %s\n", *firstProvider.PrivacyPolicyURL)
		} else {
			fmt.Printf("      Privacy Policy URL: (not provided)\n")
		}
		if firstProvider.TermsOfServiceURL != nil {
			fmt.Printf("      Terms of Service URL: %s\n", *firstProvider.TermsOfServiceURL)
		} else {
			fmt.Printf("      Terms of Service URL: (not provided)\n")
		}
		if firstProvider.StatusPageURL != nil {
			fmt.Printf("      Status Page URL: %s\n", *firstProvider.StatusPageURL)
		} else {
			fmt.Printf("      Status Page URL: (not provided)\n")
		}
	}

	// Check for well-known providers
	fmt.Printf("\n   Checking for well-known providers...\n")
	wellKnownProviders := []string{
		"openai",
		"anthropic",
		"google",
		"meta",
	}

	foundProviders := make(map[string]bool)
	for _, provider := range resp.Data {
		for _, knownProvider := range wellKnownProviders {
			if provider.Slug == knownProvider {
				foundProviders[knownProvider] = true
			}
		}
	}

	foundCount := len(foundProviders)
	fmt.Printf("   Found %d/%d well-known providers\n", foundCount, len(wellKnownProviders))

	if verbose {
		for _, knownProvider := range wellKnownProviders {
			status := "❌"
			if foundProviders[knownProvider] {
				status = "✅"
			}
			fmt.Printf("      %s %s\n", status, knownProvider)
		}
	}

	// Verify policy URLs
	fmt.Printf("\n   Validating policy URLs...\n")
	hasPrivacyPolicy := 0
	hasTermsOfService := 0
	hasStatusPage := 0

	for _, provider := range resp.Data {
		if provider.PrivacyPolicyURL != nil && *provider.PrivacyPolicyURL != "" {
			hasPrivacyPolicy++
		}
		if provider.TermsOfServiceURL != nil && *provider.TermsOfServiceURL != "" {
			hasTermsOfService++
		}
		if provider.StatusPageURL != nil && *provider.StatusPageURL != "" {
			hasStatusPage++
		}
	}

	privacyPercent := (float64(hasPrivacyPolicy) / float64(len(resp.Data))) * 100
	termsPercent := (float64(hasTermsOfService) / float64(len(resp.Data))) * 100
	statusPercent := (float64(hasStatusPage) / float64(len(resp.Data))) * 100

	fmt.Printf("   %.1f%% have privacy policy (%d/%d)\n", privacyPercent, hasPrivacyPolicy, len(resp.Data))
	fmt.Printf("   %.1f%% have terms of service (%d/%d)\n", termsPercent, hasTermsOfService, len(resp.Data))
	fmt.Printf("   %.1f%% have status page (%d/%d)\n", statusPercent, hasStatusPage, len(resp.Data))

	fmt.Printf("\n✅ List providers tests completed\n")
	return true
}

func runCreditsTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Get Credits\n")

	// Test: Get credits for authenticated user
	fmt.Printf("   Testing get credits...\n")
	start := time.Now()
	resp, err := client.GetCredits(ctx)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Failed to get credits: %v\n", err)
		return false
	}

	fmt.Printf("   ✅ Retrieved credits (%.2fs)\n", elapsed.Seconds())

	// Display credits information
	fmt.Printf("      Total Credits: $%.2f\n", resp.Data.TotalCredits)
	fmt.Printf("      Total Usage: $%.2f\n", resp.Data.TotalUsage)

	remaining := resp.Data.TotalCredits - resp.Data.TotalUsage
	fmt.Printf("      Remaining: $%.2f\n", remaining)

	if remaining < 0 {
		fmt.Printf("      ⚠️  Warning: Usage exceeds credits (negative balance)\n")
	}

	// Validate response structure
	fmt.Printf("\n   Validating response structure...\n")

	// Check that values are non-negative (usage can exceed credits, but both should be >= 0)
	if resp.Data.TotalCredits < 0 {
		fmt.Printf("   ❌ Invalid TotalCredits value: %.2f (should be >= 0)\n", resp.Data.TotalCredits)
		return false
	}
	if resp.Data.TotalUsage < 0 {
		fmt.Printf("   ❌ Invalid TotalUsage value: %.2f (should be >= 0)\n", resp.Data.TotalUsage)
		return false
	}

	fmt.Printf("   ✅ Response structure validation passed\n")

	if verbose {
		fmt.Printf("\n   Credit details:\n")
		fmt.Printf("      Total Credits: $%.4f\n", resp.Data.TotalCredits)
		fmt.Printf("      Total Usage: $%.4f\n", resp.Data.TotalUsage)
		fmt.Printf("      Remaining: $%.4f\n", remaining)

		if resp.Data.TotalCredits > 0 {
			usagePercent := (resp.Data.TotalUsage / resp.Data.TotalCredits) * 100
			fmt.Printf("      Usage: %.2f%%\n", usagePercent)
		}
	}

	// Test case variations
	fmt.Printf("\n   Testing with different contexts...\n")

	// Test with custom timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = client.GetCredits(ctxWithTimeout)
	if err != nil {
		fmt.Printf("   ❌ Failed with custom timeout: %v\n", err)
		return false
	}
	fmt.Printf("   ✅ Custom timeout context works\n")

	fmt.Printf("\n✅ Get credits tests completed\n")
	return true
}
