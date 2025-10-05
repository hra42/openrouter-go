package tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunWebSearchTest tests the web search functionality including :online suffix,
// web plugin configurations, and streaming with web search.
func RunWebSearchTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
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
			printError("Stream error", err)
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
		printSuccess("WithOnlineModel helper works correctly")
	} else {
		printError("WithOnlineModel helper failed", nil)
		return false
	}

	printSuccess("Web search tests completed")
	return true
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
