package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// Example 1: Simple chat completion
	fmt.Println("=== Example 1: Simple Chat Completion ===")
	simpleChat(client)

	// Example 2: Chat with system message
	fmt.Println("\n=== Example 2: Chat with System Message ===")
	chatWithSystem(client)

	// Example 3: Legacy completion
	fmt.Println("\n=== Example 3: Legacy Completion ===")
	legacyCompletion(client)

	// Example 4: Chat with options
	fmt.Println("\n=== Example 4: Chat with Options ===")
	chatWithOptions(client)

	// Example 5: Multi-turn conversation
	fmt.Println("\n=== Example 5: Multi-turn Conversation ===")
	multiTurnConversation(client)

	// Example 6: Error handling
	fmt.Println("\n=== Example 6: Error Handling ===")
	errorHandling(client)
}

func simpleChat(client *openrouter.Client) {
	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What is the capital of France?"),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)
}

func chatWithSystem(client *openrouter.Client) {
	messages := []openrouter.Message{
		openrouter.CreateSystemMessage("You are a helpful assistant who speaks like a pirate."),
		openrouter.CreateUserMessage("Tell me about the weather."),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func legacyCompletion(client *openrouter.Client) {
	prompt := "Once upon a time in a magical forest,"

	resp, err := client.Complete(
		context.Background(),
		prompt,
		openrouter.WithCompletionModel("openai/gpt-3.5-turbo-instruct"),
		openrouter.WithCompletionMaxTokens(50),
		openrouter.WithCompletionTemperature(0.8),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Prompt: %s\n", prompt)
	fmt.Printf("Completion: %s\n", resp.Choices[0].Text)
}

func chatWithOptions(client *openrouter.Client) {
	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Write a haiku about programming."),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("anthropic/claude-3-haiku"),
		openrouter.WithTemperature(0.9),
		openrouter.WithMaxTokens(100),
		openrouter.WithTopP(0.95),
		openrouter.WithFrequencyPenalty(0.2),
		openrouter.WithPresencePenalty(0.1),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Haiku:\n%s\n", resp.Choices[0].Message.Content)
}

func multiTurnConversation(client *openrouter.Client) {
	// Start a conversation
	messages := []openrouter.Message{
		openrouter.CreateSystemMessage("You are a helpful math tutor."),
		openrouter.CreateUserMessage("What is 15 + 27?"),
	}

	// First turn
	resp1, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("User: What is 15 + 27?\n")
	fmt.Printf("Assistant: %s\n", resp1.Choices[0].Message.Content)

	// Add assistant's response to conversation
	messages = append(messages, resp1.Choices[0].Message)

	// Second turn
	messages = append(messages, openrouter.CreateUserMessage("Now multiply that by 2."))

	resp2, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("User: Now multiply that by 2.\n")
	fmt.Printf("Assistant: %s\n", resp2.Choices[0].Message.Content)
}

func errorHandling(client *openrouter.Client) {
	// Try with invalid model
	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Hello"),
	}

	_, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("invalid/model-name"),
	)

	if err != nil {
		// Check error type
		if openrouter.IsRequestError(err) {
			reqErr := err.(*openrouter.RequestError)
			fmt.Printf("Request error (status %d): %s\n", reqErr.StatusCode, reqErr.Message)

			if reqErr.IsNotFoundError() {
				fmt.Println("Model not found!")
			}
		} else if openrouter.IsValidationError(err) {
			fmt.Printf("Validation error: %v\n", err)
		} else {
			fmt.Printf("Other error: %v\n", err)
		}
	}
}
