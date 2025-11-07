//go:build tools
// +build tools

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

// TestLocalImage is a standalone function to test local image encoding
// This can be used to quickly verify the base64 image functionality works
func TestLocalImage() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Use the test image
	imagePath := "test-image.png"

	fmt.Printf("Testing base64 image encoding with: %s\n", imagePath)

	// Create message with base64-encoded image
	message, err := openrouter.CreateUserMessageWithBase64Image(
		"What's in this image? Describe it briefly.",
		imagePath,
	)
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}

	fmt.Println("Image encoded successfully!")
	fmt.Println("Sending request to OpenRouter...")

	messages := []openrouter.Message{message}

	// Use a vision-capable model
	resp, err := client.ChatComplete(context.Background(), messages,
		openrouter.WithModel("google/gemini-2.0-flash-thinking-exp:free"),
		openrouter.WithMaxTokens(200),
	)
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}

	fmt.Printf("\n✅ Success!\n")
	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
	fmt.Printf("\nUsage:\n")
	fmt.Printf("  Prompt tokens: %d\n", resp.Usage.PromptTokens)
	fmt.Printf("  Completion tokens: %d\n", resp.Usage.CompletionTokens)
	fmt.Printf("  Total tokens: %d\n", resp.Usage.TotalTokens)
}
