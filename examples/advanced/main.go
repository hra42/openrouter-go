package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hra42/openrouter-go"
)

func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	// Example 1: Custom HTTP client with proxy
	fmt.Println("=== Example 1: Custom HTTP Client ===")
	customHTTPClient()

	// Example 2: Advanced client configuration
	fmt.Println("\n=== Example 2: Advanced Configuration ===")
	advancedConfiguration(apiKey)

	// Example 3: Function calling
	fmt.Println("\n=== Example 3: Function Calling ===")
	functionCalling(apiKey)

	// Example 4: JSON response format
	fmt.Println("\n=== Example 4: JSON Response Format ===")
	jsonResponseFormat(apiKey)

	// Example 5: Multi-modal messages
	fmt.Println("\n=== Example 5: Multi-modal Messages ===")
	multiModalMessages(apiKey)

	// Example 6: Provider preferences
	fmt.Println("\n=== Example 6: Provider Preferences ===")
	providerPreferences(apiKey)

	// Example 7: Rate limiting
	fmt.Println("\n=== Example 7: Rate Limiting ===")
	rateLimiting(apiKey)
}

func customHTTPClient() {
	// Create a custom HTTP client with specific settings
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			// Uncomment to use proxy:
			// Proxy: http.ProxyURL(proxyURL),
		},
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey),
		openrouter.WithHTTPClient(httpClient),
		openrouter.WithTimeout(60*time.Second),
	)

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What is 2+2?"),
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

func advancedConfiguration(apiKey string) {
	// Create a client with all configuration options
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey),
		openrouter.WithBaseURL("https://openrouter.ai/api/v1"),
		openrouter.WithDefaultModel("openai/gpt-3.5-turbo"),
		openrouter.WithTimeout(30*time.Second),
		openrouter.WithRetry(5, 2*time.Second),
		openrouter.WithReferer("https://myapp.com"),
		openrouter.WithAppName("AdvancedExample"),
		openrouter.WithHeader("X-Custom-Header", "custom-value"),
	)

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Say 'Configuration successful!'"),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		// Model will use default from client config
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func functionCalling(apiKey string) {
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Define available functions
	tools := []openrouter.Tool{
		{
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
						"unit": map[string]interface{}{
							"type": "string",
							"enum": []string{"celsius", "fahrenheit"},
						},
					},
					"required": []string{"location"},
				},
			},
		},
	}

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What's the weather like in Tokyo?"),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-4"),
		openrouter.WithTools(tools...),
		openrouter.WithToolChoice("auto"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Check if the model wants to call a function
	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			fmt.Printf("Function call: %s\n", toolCall.Function.Name)
			fmt.Printf("Arguments: %s\n", toolCall.Function.Arguments)

			// Here you would actually call the function and add the result to messages
			// Then make another request with the function result
		}
	} else {
		fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
	}
}

func jsonResponseFormat(apiKey string) {
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	responseFormat := openrouter.ResponseFormat{
		Type: "json_schema",
		JSONSchema: &openrouter.JSONSchema{
			Name:   "person",
			Strict: true,
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
					"age": map[string]interface{}{
						"type": "number",
					},
					"skills": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
				"required": []string{"name", "age", "skills"},
			},
		},
	}

	messages := []openrouter.Message{
		openrouter.CreateSystemMessage("You are a helpful assistant that always responds in valid JSON."),
		openrouter.CreateUserMessage("Create a fictional programmer profile with name, age, and 3 skills."),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-4"),
		openrouter.WithResponseFormat(responseFormat),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Parse the JSON response
	var profile map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content.(string)), &profile); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	fmt.Printf("Parsed profile:\n")
	fmt.Printf("  Name: %v\n", profile["name"])
	fmt.Printf("  Age: %v\n", profile["age"])
	fmt.Printf("  Skills: %v\n", profile["skills"])
}

func multiModalMessages(apiKey string) {
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Create a message with both text and image
	messages := []openrouter.Message{
		openrouter.CreateMultiModalMessage(
			"user",
			"What do you see in this image?",
			"https://example.com/sample-image.jpg",
		),
	}

	// Note: You need to use a model that supports vision
	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-4-vision-preview"),
		openrouter.WithMaxTokens(300),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		// This will likely fail without a valid image URL
		fmt.Println("(Multi-modal example requires a valid image URL and vision-capable model)")
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func providerPreferences(apiKey string) {
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Configure provider preferences
	trueVal := true
	provider := openrouter.Provider{
		Order:             []string{"OpenAI", "Anthropic"},
		RequireParameters: &trueVal,
		DataCollection:    "deny",
		AllowFallbacks:    &trueVal,
		Ignore:            []string{"Cohere"},
	}

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Hello, which provider are you using?"),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-3.5-turbo"),
		openrouter.WithProvider(provider),
		openrouter.WithMetadata(map[string]interface{}{
			"request_id": "example-123",
			"user_id":    "user-456",
		}),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
	fmt.Printf("Model used: %s\n", resp.Model)
}

func rateLimiting(apiKey string) {
	// Create a client with aggressive retry settings for rate limit handling
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey),
		openrouter.WithRetry(10, 1*time.Second),
	)

	// Create a rate limiter to control request rate
	limiter := openrouter.NewRateLimiter(2.0, 5) // 2 requests per second, burst of 5
	defer limiter.Close()

	// Make multiple requests with rate limiting
	for i := 0; i < 5; i++ {
		// Wait for rate limiter
		if err := limiter.Wait(context.Background()); err != nil {
			log.Printf("Rate limiter error: %v", err)
			continue
		}

		messages := []openrouter.Message{
			openrouter.CreateUserMessage(fmt.Sprintf("Request %d: What is %d + %d?", i+1, i, i+1)),
		}

		start := time.Now()
		resp, err := client.ChatComplete(
			context.Background(),
			messages,
			openrouter.WithModel("openai/gpt-3.5-turbo"),
			openrouter.WithMaxTokens(50),
		)

		elapsed := time.Since(start)

		if err != nil {
			// Check if it's a rate limit error
			if reqErr, ok := err.(*openrouter.RequestError); ok && reqErr.IsRateLimitError() {
				fmt.Printf("Request %d: Rate limited (will retry)\n", i+1)
			} else {
				log.Printf("Request %d error: %v", i+1, err)
			}
			continue
		}

		fmt.Printf("Request %d (%.2fs): %s\n", i+1, elapsed.Seconds(), resp.Choices[0].Message.Content)
	}
}