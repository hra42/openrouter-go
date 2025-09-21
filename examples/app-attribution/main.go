package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is not set")
	}

	// Create client with app attribution
	// These headers enable your app to appear in OpenRouter's rankings and analytics
	client := openrouter.NewClient(
		openrouter.WithAPIKey(apiKey),
		// Set your app's URL - this identifies your app in rankings
		openrouter.WithReferer("https://myapp.com"),
		// Set your app's display name for rankings and analytics
		openrouter.WithAppName("My AI Assistant"),
	)

	// Example 1: Basic chat completion with attribution
	basicExample(client)

	// Example 2: Localhost development with attribution
	localhostExample(apiKey)

	// Example 3: Using custom headers for additional attribution
	customHeadersExample(apiKey)
}

func basicExample(client *openrouter.Client) {
	fmt.Println("=== Basic Chat Completion with App Attribution ===")

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What are the benefits of app attribution in OpenRouter?"),
	}

	resp, err := client.ChatComplete(context.Background(), messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
	)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Println("\nResponse:")
		fmt.Println(resp.Choices[0].Message.Content)
	}

	// Your app is now tracked and will appear in:
	// - https://openrouter.ai/rankings (main rankings page)
	// - Model-specific "Apps" tabs (e.g., https://openrouter.ai/models/openai/gpt-3.5-turbo)
	// - App analytics at https://openrouter.ai/apps?url=https://myapp.com
	fmt.Println("\n✓ Request attributed to 'My AI Assistant' (https://myapp.com)")
}

func localhostExample(apiKey string) {
	fmt.Println("\n=== Localhost Development with App Attribution ===")

	// For localhost development, always include a title header
	// Without a title, localhost apps won't be tracked
	client := openrouter.NewClient(
		openrouter.WithAPIKey(apiKey),
		openrouter.WithReferer("http://localhost:3000"),
		// Title is essential for localhost tracking
		openrouter.WithAppName("My Development App"),
	)

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Hello from localhost!"),
	}

	resp, err := client.ChatComplete(context.Background(), messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
		openrouter.WithMaxTokens(50),
	)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Println("\nResponse:", resp.Choices[0].Message.Content)
	}

	fmt.Println("\n✓ Localhost request attributed to 'My Development App'")
}

func customHeadersExample(apiKey string) {
	fmt.Println("\n=== Custom Headers for Enhanced Attribution ===")

	// You can combine app attribution with other custom headers
	client := openrouter.NewClient(
		openrouter.WithAPIKey(apiKey),
		openrouter.WithReferer("https://enterprise.example.com"),
		openrouter.WithAppName("Enterprise AI Platform"),
		// Add custom headers for additional tracking or routing
		openrouter.WithHeader("X-Custom-User-ID", "user-12345"),
		openrouter.WithHeader("X-Custom-Session", "session-abc"),
	)

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Demonstrate custom header usage"),
	}

	resp, err := client.ChatComplete(context.Background(), messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
		openrouter.WithMaxTokens(50),
		// You can also use metadata for request-specific tracking
		openrouter.WithMetadata(map[string]interface{}{
			"department": "engineering",
			"project":    "ai-integration",
		}),
	)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Println("\nResponse:", resp.Choices[0].Message.Content)
	}

	fmt.Println("\n✓ Request attributed with custom headers and metadata")
	fmt.Println("  - App: Enterprise AI Platform")
	fmt.Println("  - URL: https://enterprise.example.com")
	fmt.Println("  - Additional headers and metadata included")
}