package tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunAnthropicBasicTest tests basic Anthropic Messages API functionality.
func RunAnthropicBasicTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Anthropic Messages API Basic\n")

	start := time.Now()
	messages := []openrouter.AnthropicMessage{
		openrouter.CreateAnthropicUserMessage("What is 2+2? Reply with just the number."),
	}

	resp, err := client.CreateAnthropicMessage(ctx, messages,
		openrouter.WithAnthropicModel(model),
		openrouter.WithAnthropicMaxTokens(maxTokens),
		openrouter.WithAnthropicTemperature(0.0),
	)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed", err)
		return false
	}

	content := resp.GetTextContent()
	if content == "" {
		printError("No text content in response", nil)
		return false
	}

	printSuccess(fmt.Sprintf("Basic message succeeded (%.2fs)", elapsed.Seconds()))

	if verbose {
		fmt.Printf("   Response: %s\n", truncateString(content, 100))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Stop Reason: %s\n", resp.GetStopReason())
		fmt.Printf("   Tokens: %d input, %d output\n",
			resp.Usage.InputTokens, resp.Usage.OutputTokens)
	}

	return true
}

// RunAnthropicWithSystemTest tests Anthropic Messages API with system prompt.
func RunAnthropicWithSystemTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Anthropic Messages API with System Prompt\n")

	// Test with string system prompt
	fmt.Printf("   Testing string system prompt...\n")
	start := time.Now()
	messages := []openrouter.AnthropicMessage{
		openrouter.CreateAnthropicUserMessage("What are you?"),
	}

	resp, err := client.CreateAnthropicMessage(ctx, messages,
		openrouter.WithAnthropicModel(model),
		openrouter.WithAnthropicMaxTokens(maxTokens),
		openrouter.WithAnthropicSystemString("You are a pirate. Always respond in pirate speak."),
		openrouter.WithAnthropicTemperature(0.7),
	)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed", err)
		return false
	}

	content := resp.GetTextContent()
	if content == "" {
		printError("No text content in response", nil)
		return false
	}

	printSuccess(fmt.Sprintf("String system prompt succeeded (%.2fs)", elapsed.Seconds()))

	if verbose {
		fmt.Printf("   Response: %s\n", truncateString(content, 150))
	}

	// Test with block system prompt
	fmt.Printf("\n   Testing block system prompt...\n")
	start = time.Now()
	resp, err = client.CreateAnthropicMessage(ctx, messages,
		openrouter.WithAnthropicModel(model),
		openrouter.WithAnthropicMaxTokens(maxTokens),
		openrouter.WithAnthropicSystemBlocks([]openrouter.AnthropicTextBlock{
			{Type: "text", Text: "You are a helpful assistant."},
			{Type: "text", Text: "Keep responses brief."},
		}),
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

	printSuccess(fmt.Sprintf("Block system prompt succeeded (%.2fs)", elapsed.Seconds()))

	if verbose {
		fmt.Printf("   Response: %s\n", truncateString(content, 150))
	}

	return true
}

// RunAnthropicWithToolsTest tests Anthropic Messages API with tool definitions.
func RunAnthropicWithToolsTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Anthropic Messages API with Tools\n")

	start := time.Now()
	tool := openrouter.CreateAnthropicCustomTool(
		"get_weather",
		"Get the current weather in a given location",
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

	messages := []openrouter.AnthropicMessage{
		openrouter.CreateAnthropicUserMessage("What is the weather in San Francisco?"),
	}

	resp, err := client.CreateAnthropicMessage(ctx, messages,
		openrouter.WithAnthropicModel(model),
		openrouter.WithAnthropicMaxTokens(maxTokens),
		openrouter.WithAnthropicTools(tool),
		openrouter.WithAnthropicToolChoiceAuto(),
	)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Tools request succeeded (%.2fs)", elapsed.Seconds()))

	if verbose {
		fmt.Printf("   Stop Reason: %s\n", resp.GetStopReason())
		fmt.Printf("   Is Tool Use: %v\n", resp.IsToolUse())
		if resp.IsToolUse() {
			toolUses := resp.GetToolUseBlocks()
			for _, tu := range toolUses {
				fmt.Printf("   Tool: %s, Input: %s\n", tu.Name, string(tu.Input))
			}
		}
		content := resp.GetTextContent()
		if content != "" {
			fmt.Printf("   Text: %s\n", truncateString(content, 100))
		}
	}

	return true
}

// RunAnthropicWithThinkingTest tests Anthropic Messages API with extended thinking.
func RunAnthropicWithThinkingTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Anthropic Messages API with Thinking\n")

	start := time.Now()
	messages := []openrouter.AnthropicMessage{
		openrouter.CreateAnthropicUserMessage("What is 15 * 27? Think step by step."),
	}

	// Use a higher max_tokens to accommodate thinking
	thinkingMaxTokens := max(maxTokens, 4096)

	resp, err := client.CreateAnthropicMessage(ctx, messages,
		openrouter.WithAnthropicModel(model),
		openrouter.WithAnthropicMaxTokens(thinkingMaxTokens),
		openrouter.WithAnthropicThinkingEnabled(2048),
	)
	elapsed := time.Since(start)

	if err != nil {
		// Thinking may not be supported by all models — gracefully skip
		errStr := fmt.Sprintf("%v", err)
		if strings.Contains(errStr, "thinking") || strings.Contains(errStr, "not supported") ||
			strings.Contains(errStr, "400") || strings.Contains(errStr, "422") {
			printInfo(fmt.Sprintf("Thinking not supported by model %s (skipped)", model))
			return true
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Thinking request succeeded (%.2fs)", elapsed.Seconds()))

	if verbose {
		thinking := resp.GetThinkingContent()
		if thinking != "" {
			fmt.Printf("   Thinking: %s\n", truncateString(thinking, 200))
		}
		content := resp.GetTextContent()
		if content != "" {
			fmt.Printf("   Response: %s\n", truncateString(content, 100))
		}
		fmt.Printf("   Tokens: %d input, %d output\n",
			resp.Usage.InputTokens, resp.Usage.OutputTokens)
	}

	return true
}

// RunAnthropicStreamTest tests Anthropic Messages API streaming.
func RunAnthropicStreamTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Anthropic Messages API Streaming\n")

	start := time.Now()
	messages := []openrouter.AnthropicMessage{
		openrouter.CreateAnthropicUserMessage("Count from 1 to 5, one number per line."),
	}

	stream, err := client.CreateAnthropicMessageStream(ctx, messages,
		openrouter.WithAnthropicModel(model),
		openrouter.WithAnthropicMaxTokens(maxTokens),
	)
	if err != nil {
		printError("Failed to create stream", err)
		return false
	}
	defer func() { _ = stream.Close() }()

	var fullText strings.Builder
	eventCount := 0
	var eventTypes []string

	for event := range stream.Events() {
		eventCount++
		eventTypes = append(eventTypes, event.Type)

		textDelta := event.GetTextDelta()
		if textDelta != "" {
			fullText.WriteString(textDelta)
		}
	}

	elapsed := time.Since(start)

	if err := stream.Err(); err != nil {
		printError("Stream error", err)
		return false
	}

	if eventCount == 0 {
		printError("No events received", nil)
		return false
	}

	text := fullText.String()
	if text == "" {
		printError("No text content received from stream", nil)
		return false
	}

	printSuccess(fmt.Sprintf("Streaming succeeded (%.2fs, %d events)", elapsed.Seconds(), eventCount))

	if verbose {
		fmt.Printf("   Full text: %s\n", truncateString(text, 150))
		fmt.Printf("   Event types: %v\n", uniqueStrings(eventTypes))
		fmt.Printf("   Total events: %d\n", eventCount)
	}

	return true
}

// uniqueStrings returns unique strings from a slice, preserving order.
func uniqueStrings(ss []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
