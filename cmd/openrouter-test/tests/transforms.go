package tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunTransformsTest tests message transforms functionality including middle-out compression,
// transform disabling, and streaming with transforms.
func RunTransformsTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
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
		printError("Chat with transforms failed", err)
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
		printError("Chat without transforms failed", err)
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
		printError("Failed to create stream with transforms", err)
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
		printError("Stream error", err)
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

	printSuccess("Message transforms tests completed")
	return true
}
