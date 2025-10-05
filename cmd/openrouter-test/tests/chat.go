package tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunChatTest tests basic chat completion
func RunChatTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
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
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	if verbose || true { // Always show some output
		fmt.Printf("   Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Tokens: %d prompt, %d completion, %d total\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	return true
}

// RunStreamTest tests streaming chat completion
func RunStreamTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
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
		printError("Failed to create stream", err)
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
		fmt.Printf("\n")
		printError("Stream error", err)
		return false
	}

	fmt.Printf("\n")
	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))
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

// RunCompletionTest tests legacy completion API
func RunCompletionTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
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
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	if verbose || true {
		fmt.Printf("   Prompt: %s\n", prompt)
		fmt.Printf("   Completion: %s\n", strings.TrimSpace(resp.Choices[0].Text))
		fmt.Printf("   Tokens: %d total\n", resp.Usage.TotalTokens)
	}

	return true
}

// RunErrorTest tests error handling
func RunErrorTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Error Handling\n")

	// Test with invalid model
	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Test message"),
	}

	_, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("invalid/nonexistent-model-xyz"),
	)

	if err == nil {
		printError("Expected error but got none", nil)
		return false
	}

	// Verify it's our custom error type
	if reqErr, ok := err.(*openrouter.RequestError); ok {
		printSuccess("Got expected error")
		printVerbose(verbose, "Error type: RequestError")
		printVerbose(verbose, "Status code: %d", reqErr.StatusCode)
		printVerbose(verbose, "Message: %s", reqErr.Message)
		return true
	}

	printError("Got unexpected error type", err)
	return false
}
