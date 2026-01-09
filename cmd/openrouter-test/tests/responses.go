package tests

// WARNING: Responses API tests are for a BETA API that may have breaking changes.
// These tests may fail or behave unexpectedly as the API evolves.

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunResponsesBasicTest tests basic Responses API functionality
// NOTE: This tests a BETA API - expect potential failures due to API changes
func RunResponsesBasicTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Responses API Basic\n")

	// Test 1: Simple string input
	fmt.Printf("   Testing simple string input...\n")
	start := time.Now()
	resp, err := client.CreateResponse(ctx, "What is 2+2? Reply with just the number.",
		openrouter.WithResponsesModel(model),
		openrouter.WithResponsesMaxOutputTokens(maxTokens),
		openrouter.WithResponsesTemperature(0.7),
	)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Simple string input succeeded (%.2fs)", elapsed.Seconds()))

	content := resp.GetTextContent()
	if content == "" {
		printError("No text content in response", nil)
		return false
	}

	if verbose {
		fmt.Printf("   Response: %s\n", truncateString(content, 100))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Tokens: %d input, %d output, %d total\n",
			resp.Usage.InputTokens, resp.Usage.OutputTokens, resp.Usage.TotalTokens)
	}

	// Test 2: Structured input
	fmt.Printf("\n   Testing structured input...\n")
	input := []openrouter.ResponsesInputItem{
		openrouter.CreateResponsesSystemMessage("You are a helpful assistant. Keep responses brief."),
		openrouter.CreateResponsesUserMessage("What is the capital of France? Reply in one word."),
	}

	start = time.Now()
	resp, err = client.CreateResponse(ctx, input,
		openrouter.WithResponsesModel(model),
		openrouter.WithResponsesMaxOutputTokens(maxTokens),
	)
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed", err)
		return false
	}

	content = resp.GetTextContent()
	if content == "" {
		printError("No text content in response", nil)
		return false
	}

	printSuccess(fmt.Sprintf("Structured input succeeded (%.2fs)", elapsed.Seconds()))

	if verbose {
		fmt.Printf("   Response: %s\n", truncateString(content, 100))
	}

	printSuccess("All Responses API basic tests passed")
	return true
}

// RunResponsesReasoningTest tests Responses API reasoning functionality
func RunResponsesReasoningTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Responses API Reasoning\n")

	// Note: Reasoning is primarily supported by specific models like openai/o4-mini
	// This test may need a model that supports reasoning

	start := time.Now()
	resp, err := client.CreateResponse(ctx, "What is 15 multiplied by 17? Think step by step.",
		openrouter.WithResponsesModel(model),
		openrouter.WithResponsesMaxOutputTokens(maxTokens),
		openrouter.WithResponsesReasoningEffort(openrouter.ReasoningEffortMedium),
	)
	elapsed := time.Since(start)

	if err != nil {
		// Some models may not support reasoning - check if it's a model capability error
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if strings.Contains(reqErr.Message, "reasoning") || strings.Contains(reqErr.Message, "parameter") {
				fmt.Printf("⚠️  Skipped: Model %s may not support reasoning\n", model)
				return true
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Reasoning request succeeded (%.2fs)", elapsed.Seconds()))

	content := resp.GetTextContent()
	if content == "" {
		printError("No text content in response", nil)
		return false
	}

	if verbose {
		fmt.Printf("   Response: %s\n", truncateString(content, 200))

		// Check for reasoning summary if available
		summary := resp.GetReasoningSummary()
		if len(summary) > 0 {
			fmt.Printf("   Reasoning summary:\n")
			for i, step := range summary {
				fmt.Printf("     %d. %s\n", i+1, truncateString(step, 100))
			}
		}

		if resp.Usage.ReasoningTokens > 0 {
			fmt.Printf("   Reasoning tokens: %d\n", resp.Usage.ReasoningTokens)
		}
	}

	printSuccess("Responses API reasoning test passed")
	return true
}

// RunResponsesToolsTest tests Responses API tool calling functionality
func RunResponsesToolsTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Responses API Tools\n")

	weatherTool := openrouter.Tool{
		Type: "function",
		Function: openrouter.Function{
			Name:        "get_weather",
			Description: "Get the current weather in a given location",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "The city and state, e.g. San Francisco, CA",
					},
				},
				"required": []string{"location"},
			},
		},
	}

	start := time.Now()
	resp, err := client.CreateResponse(ctx, "What's the weather like in Paris today?",
		openrouter.WithResponsesModel(model),
		openrouter.WithResponsesMaxOutputTokens(maxTokens),
		openrouter.WithResponsesTools(weatherTool),
	)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Tool calling request succeeded (%.2fs)", elapsed.Seconds()))

	// Check if we got function calls or a regular response
	calls := resp.GetFunctionCalls()
	content := resp.GetTextContent()

	if len(calls) > 0 {
		fmt.Printf("   Function calls received: %d\n", len(calls))
		for i, call := range calls {
			if verbose {
				fmt.Printf("   Call %d: %s(%s)\n", i+1, call.Name, truncateString(call.Arguments, 50))
			}
		}
	} else if content != "" {
		fmt.Printf("   Model responded with text instead of tool call\n")
		if verbose {
			fmt.Printf("   Response: %s\n", truncateString(content, 100))
		}
	}

	printSuccess("Responses API tools test passed")
	return true
}

// RunResponsesWebSearchTest tests Responses API web search functionality
func RunResponsesWebSearchTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Responses API Web Search\n")

	start := time.Now()
	resp, err := client.CreateResponse(ctx, "What is the current weather in New York? Cite your sources.",
		openrouter.WithResponsesModel(model),
		openrouter.WithResponsesMaxOutputTokens(maxTokens),
		openrouter.WithResponsesWebSearch(3),
	)
	elapsed := time.Since(start)

	if err != nil {
		// Web search may not be available for all models or may require specific configuration
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if strings.Contains(reqErr.Message, "plugin") || strings.Contains(reqErr.Message, "web") {
				fmt.Printf("⚠️  Skipped: Web search may not be available for model %s\n", model)
				return true
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Web search request succeeded (%.2fs)", elapsed.Seconds()))

	content := resp.GetTextContent()
	if content == "" {
		printError("No text content in response", nil)
		return false
	}

	annotations := resp.GetAnnotations()

	if verbose {
		fmt.Printf("   Response: %s\n", truncateString(content, 200))
		if len(annotations) > 0 {
			fmt.Printf("   Citations: %d\n", len(annotations))
			for i, ann := range annotations {
				if ann.Type == "url_citation" {
					fmt.Printf("   [%d] %s\n", i+1, truncateString(ann.URL, 60))
				}
			}
		} else {
			fmt.Printf("   No citations in response\n")
		}
	}

	printSuccess("Responses API web search test passed")
	return true
}

// RunResponsesStreamTest tests Responses API streaming functionality
func RunResponsesStreamTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Responses API Streaming\n")

	start := time.Now()
	stream, err := client.CreateResponseStream(ctx, "Count from 1 to 5, one number per line.",
		openrouter.WithResponsesModel(model),
		openrouter.WithResponsesMaxOutputTokens(maxTokens),
	)
	if err != nil {
		printError("Failed to create stream", err)
		return false
	}
	defer stream.Close()

	fmt.Printf("   Streaming: ")
	var fullResponse string
	eventCount := 0

	for event := range stream.Events() {
		eventCount++
		content := event.GetTextContent()
		if content != "" && content != fullResponse {
			// Only print the new content
			if len(content) > len(fullResponse) {
				newContent := content[len(fullResponse):]
				if verbose {
					fmt.Print(newContent)
				}
			}
			fullResponse = content
		}
	}

	elapsed := time.Since(start)

	if err := stream.Err(); err != nil {
		fmt.Printf("\n")
		printError("Stream error", err)
		return false
	}

	fmt.Printf("\n")
	printSuccess(fmt.Sprintf("Streaming succeeded (%.2fs)", elapsed.Seconds()))
	fmt.Printf("   Events received: %d\n", eventCount)

	if !verbose {
		if len(fullResponse) > 100 {
			fullResponse = fullResponse[:100] + "..."
		}
		fmt.Printf("   Response: %s\n", fullResponse)
	}

	printSuccess("Responses API streaming test passed")
	return true
}
