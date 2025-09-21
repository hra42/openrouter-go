// Package main demonstrates the use of message transforms in the OpenRouter Go client.
// Transforms help manage long conversations that exceed model context windows.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	// Create client
	client := openrouter.NewClient(
		openrouter.WithAPIKey(apiKey),
		openrouter.WithAppName("OpenRouter-Go-Transforms-Example"),
	)

	ctx := context.Background()

	// Example 1: Middle-out transform for long conversations
	fmt.Println("=== Example 1: Middle-out Transform ===")
	demonstrateMiddleOut(ctx, client)

	// Example 2: Disabling transforms
	fmt.Println("\n=== Example 2: Disabling Transforms ===")
	demonstrateDisabledTransforms(ctx, client)

	// Example 3: Transforms with streaming
	fmt.Println("\n=== Example 3: Transforms with Streaming ===")
	demonstrateStreamingTransforms(ctx, client)
}

func demonstrateMiddleOut(ctx context.Context, client *openrouter.Client) {
	// Create a long conversation that might exceed context for smaller models
	messages := []openrouter.Message{
		openrouter.CreateSystemMessage("You are a helpful assistant."),
		openrouter.CreateUserMessage("Let me tell you a very long story..."),
		openrouter.CreateAssistantMessage("I'm listening to your story."),
		openrouter.CreateUserMessage("Once upon a time, there was a kingdom far, far away..."),
		openrouter.CreateAssistantMessage("That sounds interesting, please continue."),
		openrouter.CreateUserMessage("The kingdom had many adventures and tales..."),
		openrouter.CreateAssistantMessage("What happened next?"),
		openrouter.CreateUserMessage("Many things happened over many years..."),
		openrouter.CreateAssistantMessage("I see."),
		// Add the important question at the end
		openrouter.CreateUserMessage("Forget the story. What is the capital of France? Reply with just the city name."),
	}

	// Use middle-out transform to handle the conversation
	// If it exceeds context, middle messages will be compressed
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("meta-llama/llama-3.1-8b-instruct"),
		openrouter.WithMaxTokens(20),
		openrouter.WithTransforms("middle-out"), // Enable middle-out compression
	)

	if err != nil {
		log.Printf("Error with middle-out transform: %v", err)
		return
	}

	fmt.Printf("Response with middle-out: %s\n", resp.Choices[0].Message.Content)
	fmt.Printf("(Middle content may have been compressed to fit context)\n")
}

func demonstrateDisabledTransforms(ctx context.Context, client *openrouter.Client) {
	// Simple message that doesn't need transforms
	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What is 2+2? Reply with just the number."),
	}

	// Explicitly disable transforms (empty array)
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
		openrouter.WithMaxTokens(10),
		openrouter.WithTransforms(), // Empty array disables transforms
	)

	if err != nil {
		log.Printf("Error with disabled transforms: %v", err)
		return
	}

	fmt.Printf("Response without transforms: %s\n", resp.Choices[0].Message.Content)
}

func demonstrateStreamingTransforms(ctx context.Context, client *openrouter.Client) {
	// Create a conversation for streaming
	messages := []openrouter.Message{
		openrouter.CreateSystemMessage("You are a concise assistant."),
		// Add some context that might be compressed
		openrouter.CreateUserMessage("I need help with something."),
		openrouter.CreateAssistantMessage("I'm here to help."),
		openrouter.CreateUserMessage("Let me provide some background first..."),
		openrouter.CreateAssistantMessage("Go ahead."),
		// The actual request
		openrouter.CreateUserMessage("List three primary colors. Be brief."),
	}

	// Stream with transforms enabled
	stream, err := client.ChatCompleteStream(ctx, messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
		openrouter.WithMaxTokens(50),
		openrouter.WithTransforms("middle-out"),
	)

	if err != nil {
		log.Printf("Error creating stream with transforms: %v", err)
		return
	}
	defer stream.Close()

	fmt.Print("Streaming response: ")
	var fullResponse strings.Builder

	for event := range stream.Events() {
		for _, choice := range event.Choices {
			if choice.Delta != nil {
				if content, ok := choice.Delta.Content.(string); ok {
					fullResponse.WriteString(content)
					fmt.Print(content)
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		log.Printf("\nStream error: %v", err)
		return
	}

	fmt.Println("\n(Transform applied if needed for context management)")
}
