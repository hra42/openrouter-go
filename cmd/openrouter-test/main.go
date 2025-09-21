package main

import (
	"context"
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
		test      = flag.String("test", "all", "Test to run: all, chat, stream, completion, error")
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
	client := openrouter.NewClient(*apiKey,
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