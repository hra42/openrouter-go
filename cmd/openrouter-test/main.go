package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
	"github.com/hra42/openrouter-go/cmd/openrouter-test/tests"
)

func main() {
	// Command-line flags
	var (
		apiKey    = flag.String("key", os.Getenv("OPENROUTER_API_KEY"), "OpenRouter API key (or set OPENROUTER_API_KEY env var)")
		model     = flag.String("model", "openai/gpt-3.5-turbo", "Model to use")
		test      = flag.String("test", "all", "Test to run: all, chat, stream, completion, error, provider, zdr, suffix, price, structured, tools, transforms, websearch, models, endpoints, providers, credits, activity, key, listkeys, createkey, updatekey, deletekey")
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
		if tests.RunChatTest(ctx, client, *model, *maxTokens, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "stream":
		if tests.RunStreamTest(ctx, client, *model, *maxTokens, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "completion":
		if tests.RunCompletionTest(ctx, client, *model, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "error":
		if tests.RunErrorTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "provider":
		if tests.RunProviderRoutingTest(ctx, client, *model, *maxTokens, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "zdr":
		if tests.RunZDRTest(ctx, client, *model, *maxTokens, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "suffix":
		if tests.RunModelSuffixTest(ctx, client, *model, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "price":
		if tests.RunPriceConstraintTest(ctx, client, *model, *maxTokens, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "structured":
		if tests.RunStructuredOutputTest(ctx, client, *model, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "tools":
		if tests.RunToolCallingTest(ctx, client, *model, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "transforms":
		if tests.RunTransformsTest(ctx, client, *model, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "websearch":
		if tests.RunWebSearchTest(ctx, client, *model, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "models":
		if tests.RunModelsTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "endpoints":
		if tests.RunModelEndpointsTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "providers":
		if tests.RunProvidersTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "credits":
		if tests.RunCreditsTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "activity":
		if tests.RunActivityTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "key":
		if tests.RunKeyTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "listkeys":
		if tests.RunListKeysTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "createkey":
		if tests.RunCreateKeyTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "updatekey":
		if tests.RunUpdateKeyTest(ctx, client, *verbose) {
			success = 1
		} else {
			failed = 1
		}
	case "deletekey":
		if tests.RunDeleteKeyTest(ctx, client, *verbose) {
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
	testCases := []struct {
		name string
		fn   func() bool
	}{
		{"Chat Completion", func() bool { return tests.RunChatTest(ctx, client, model, maxTokens, verbose) }},
		{"Streaming", func() bool { return tests.RunStreamTest(ctx, client, model, maxTokens, verbose) }},
		{"Legacy Completion", func() bool { return tests.RunCompletionTest(ctx, client, model, verbose) }},
		{"Error Handling", func() bool { return tests.RunErrorTest(ctx, client, verbose) }},
		{"Provider Routing", func() bool { return tests.RunProviderRoutingTest(ctx, client, model, maxTokens, verbose) }},
		{"ZDR", func() bool { return tests.RunZDRTest(ctx, client, model, maxTokens, verbose) }},
		{"Model Suffixes", func() bool { return tests.RunModelSuffixTest(ctx, client, model, verbose) }},
		{"Price Constraints", func() bool { return tests.RunPriceConstraintTest(ctx, client, model, maxTokens, verbose) }},
		{"Structured Output", func() bool { return tests.RunStructuredOutputTest(ctx, client, model, verbose) }},
		{"Tool Calling", func() bool { return tests.RunToolCallingTest(ctx, client, model, verbose) }},
		{"Message Transforms", func() bool { return tests.RunTransformsTest(ctx, client, model, verbose) }},
		{"List Models", func() bool { return tests.RunModelsTest(ctx, client, verbose) }},
		{"Model Endpoints", func() bool { return tests.RunModelEndpointsTest(ctx, client, verbose) }},
		{"List Providers", func() bool { return tests.RunProvidersTest(ctx, client, verbose) }},
		{"Get Credits", func() bool { return tests.RunCreditsTest(ctx, client, verbose) }},
		{"Get Activity", func() bool { return tests.RunActivityTest(ctx, client, verbose) }},
		{"Get API Key Info", func() bool { return tests.RunKeyTest(ctx, client, verbose) }},
		{"List API Keys", func() bool { return tests.RunListKeysTest(ctx, client, verbose) }},
		{"Create API Key", func() bool { return tests.RunCreateKeyTest(ctx, client, verbose) }},
		{"Update API Key", func() bool { return tests.RunUpdateKeyTest(ctx, client, verbose) }},
		{"Delete API Key", func() bool { return tests.RunDeleteKeyTest(ctx, client, verbose) }},
	}

	for _, tc := range testCases {
		if tc.fn() {
			success++
		} else {
			failed++
		}
		fmt.Println()
	}

	return success, failed
}
