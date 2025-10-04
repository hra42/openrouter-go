package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	// Create a new client
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Example 1: Basic streaming
	fmt.Println("=== Example 1: Basic Chat Streaming ===")
	basicStreaming(client)

	// Example 2: Streaming with context cancellation
	fmt.Println("\n=== Example 2: Streaming with Cancellation ===")
	streamingWithCancellation(client)

	// Example 3: Legacy completion streaming
	fmt.Println("\n=== Example 3: Legacy Completion Streaming ===")
	legacyStreaming(client)

	// Example 4: Collecting streamed responses
	fmt.Println("\n=== Example 4: Collecting Streamed Responses ===")
	collectingResponses(client)

	// Example 5: Streaming with error handling
	fmt.Println("\n=== Example 5: Streaming with Error Handling ===")
	streamingErrorHandling(client)
}

func basicStreaming(client *openrouter.Client) {
	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Write a short story about a robot learning to paint. Make it 3 paragraphs."),
	}

	stream, err := client.ChatCompleteStream(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
	)
	if err != nil {
		log.Printf("Error creating stream: %v", err)
		return
	}
	defer stream.Close()

	fmt.Println("Streaming response:")
	fmt.Println("---")

	// Read events as they arrive
	for event := range stream.Events() {
		// Each event contains a delta with partial content
		for _, choice := range event.Choices {
			if choice.Delta != nil {
				if content, ok := choice.Delta.Content.(string); ok {
					fmt.Print(content)
				}
			}
		}
	}

	// Check for any errors that occurred during streaming
	if err := stream.Err(); err != nil {
		log.Printf("\nStream error: %v", err)
	}

	fmt.Println("\n---")
}

func streamingWithCancellation(client *openrouter.Client) {
	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Count from 1 to 100 slowly, with a description for each number."),
	}

	// Create a context that will be cancelled after 2 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stream, err := client.ChatCompleteStream(
		ctx,
		messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
	)
	if err != nil {
		log.Printf("Error creating stream: %v", err)
		return
	}
	defer stream.Close()

	fmt.Println("Streaming (will cancel after 2 seconds):")
	fmt.Println("---")

	for event := range stream.Events() {
		for _, choice := range event.Choices {
			if choice.Delta != nil {
				if content, ok := choice.Delta.Content.(string); ok {
					fmt.Print(content)
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("\n[Stream cancelled due to timeout]")
		} else {
			log.Printf("\nStream error: %v", err)
		}
	}

	fmt.Println("\n---")
}

func legacyStreaming(client *openrouter.Client) {
	prompt := "The three most important programming concepts are:"

	stream, err := client.CompleteStream(
		context.Background(),
		prompt,
		openrouter.WithCompletionModel("openai/gpt-3.5-turbo-instruct"),
		openrouter.WithCompletionMaxTokens(100),
		openrouter.WithCompletionTemperature(0.7),
	)
	if err != nil {
		log.Printf("Error creating stream: %v", err)
		return
	}
	defer stream.Close()

	fmt.Printf("Prompt: %s\n", prompt)
	fmt.Println("Streaming completion:")
	fmt.Println("---")

	for event := range stream.Events() {
		for _, choice := range event.Choices {
			fmt.Print(choice.Text)
		}
	}

	if err := stream.Err(); err != nil {
		log.Printf("\nStream error: %v", err)
	}

	fmt.Println("\n---")
}

func collectingResponses(client *openrouter.Client) {
	messages := []openrouter.Message{
		openrouter.CreateUserMessage("List 5 interesting facts about space."),
	}

	stream, err := client.ChatCompleteStream(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
	)
	if err != nil {
		log.Printf("Error creating stream: %v", err)
		return
	}
	defer stream.Close()

	// Collect all responses
	var responses []openrouter.ChatCompletionResponse
	var fullContent strings.Builder

	fmt.Println("Collecting streamed responses...")

	for event := range stream.Events() {
		responses = append(responses, event)

		// Also build the full content
		for _, choice := range event.Choices {
			if choice.Delta != nil {
				if content, ok := choice.Delta.Content.(string); ok {
					fullContent.WriteString(content)
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		log.Printf("Stream error: %v", err)
		return
	}

	fmt.Printf("Received %d streaming events\n", len(responses))
	fmt.Println("Full response:")
	fmt.Println("---")
	fmt.Println(fullContent.String())
	fmt.Println("---")

	// You can also use the helper function
	concatenated := openrouter.ConcatenateChatStreamResponses(responses)
	fmt.Printf("Helper function result matches: %v\n", concatenated == fullContent.String())
}

func streamingErrorHandling(client *openrouter.Client) {
	// Try streaming with an invalid configuration
	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Hello"),
	}

	// Use a very low max_tokens to potentially trigger a finish reason
	stream, err := client.ChatCompleteStream(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
		openrouter.WithMaxTokens(5), // Very low token limit
	)
	if err != nil {
		log.Printf("Error creating stream: %v", err)

		// Check error type
		if openrouter.IsRequestError(err) {
			reqErr := err.(*openrouter.RequestError)
			fmt.Printf("Request error (status %d): %s\n", reqErr.StatusCode, reqErr.Message)
		}
		return
	}
	defer stream.Close()

	fmt.Println("Streaming with token limit:")
	fmt.Println("---")

	var finishReason string
	for event := range stream.Events() {
		for _, choice := range event.Choices {
			if choice.Delta != nil {
				if content, ok := choice.Delta.Content.(string); ok {
					fmt.Print(content)
				}
			}

			// Check finish reason
			if choice.FinishReason != "" {
				finishReason = choice.FinishReason
			}
		}
	}

	fmt.Println("\n---")

	if finishReason != "" {
		fmt.Printf("Stream finished with reason: %s\n", finishReason)
	}

	if err := stream.Err(); err != nil {
		if openrouter.IsStreamError(err) {
			streamErr := err.(*openrouter.StreamError)
			fmt.Printf("Stream error: %s\n", streamErr.Message)
			if streamErr.Unwrap() != nil {
				fmt.Printf("Underlying error: %v\n", streamErr.Unwrap())
			}
		} else {
			log.Printf("Other error: %v", err)
		}
	}
}
