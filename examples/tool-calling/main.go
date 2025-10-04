// Package main demonstrates tool calling with the OpenRouter API.
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

// searchGutenbergBooks simulates searching the Project Gutenberg library
func searchGutenbergBooks(searchTerms []string) (string, error) {
	// Simulated search results for demonstration
	// In a real implementation, this would call the actual Gutenberg API
	results := []map[string]interface{}{
		{
			"id":    4300,
			"title": "Ulysses",
			"authors": []map[string]string{
				{"name": "Joyce, James"},
			},
		},
		{
			"id":    2814,
			"title": "Dubliners",
			"authors": []map[string]string{
				{"name": "Joyce, James"},
			},
		},
		{
			"id":    4217,
			"title": "A Portrait of the Artist as a Young Man",
			"authors": []map[string]string{
				{"name": "Joyce, James"},
			},
		},
	}

	jsonResult, err := json.Marshal(results)
	if err != nil {
		return "", err
	}
	return string(jsonResult), nil
}

// getCurrentWeather simulates getting current weather for a location
func getCurrentWeather(location string, unit string) (string, error) {
	// Simulated weather data for demonstration
	weatherData := map[string]interface{}{
		"location":       location,
		"temperature":    72,
		"unit":           unit,
		"conditions":     "Partly cloudy",
		"humidity":       65,
		"wind_speed":     10,
		"wind_direction": "NW",
	}

	jsonResult, err := json.Marshal(weatherData)
	if err != nil {
		return "", err
	}
	return string(jsonResult), nil
}

func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	// Create client with custom configuration
	client := openrouter.NewClient(
		openrouter.WithAPIKey(apiKey),
		openrouter.WithTimeout(30*time.Second),
		openrouter.WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
		}),
	)

	ctx := context.Background()

	// Example 1: Simple tool calling for book search
	fmt.Println("=== Example 1: Book Search with Tool Calling ===")
	runBookSearchExample(ctx, client)

	// Example 2: Weather tool with different parameter types
	fmt.Println("\n=== Example 2: Weather Tool with Parameters ===")
	runWeatherExample(ctx, client)

	// Example 3: Multiple tools available
	fmt.Println("\n=== Example 3: Multiple Tools Available ===")
	runMultipleToolsExample(ctx, client)

	// Example 4: Forced tool choice
	fmt.Println("\n=== Example 4: Forced Tool Choice ===")
	runForcedToolChoiceExample(ctx, client)
}

func runBookSearchExample(ctx context.Context, client *openrouter.Client) {
	// Define the tool
	tools := []openrouter.Tool{
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "search_gutenberg_books",
				Description: "Search for books in the Project Gutenberg library",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"search_terms": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
							"description": "List of search terms to find books",
						},
					},
					"required": []string{"search_terms"},
				},
			},
		},
	}

	// Initial message asking about books
	messages := []openrouter.Message{
		{
			Role:    "user",
			Content: "What are the titles of some James Joyce books?",
		},
	}

	// Make the first request with tools
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o-mini"),
		openrouter.WithTools(tools...),
		openrouter.WithMaxTokens(1000),
	)
	if err != nil {
		log.Printf("Error in first request: %v", err)
		return
	}

	// Check if the model wants to use a tool
	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		// Add the assistant's response to messages
		messages = append(messages, resp.Choices[0].Message)

		// Process each tool call
		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			fmt.Printf("Model requested tool: %s\n", toolCall.Function.Name)
			fmt.Printf("With arguments: %s\n", toolCall.Function.Arguments)

			// Parse the arguments
			var args struct {
				SearchTerms []string `json:"search_terms"`
			}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				log.Printf("Error parsing arguments: %v", err)
				continue
			}

			// Execute the tool
			result, err := searchGutenbergBooks(args.SearchTerms)
			if err != nil {
				log.Printf("Error executing tool: %v", err)
				continue
			}

			// Add the tool result to messages
			messages = append(messages, openrouter.Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: toolCall.ID,
			})
		}

		// Make a second request with the tool results
		resp, err = client.ChatComplete(ctx, messages,
			openrouter.WithModel("openai/gpt-4o-mini"),
			openrouter.WithTools(tools...),
			openrouter.WithMaxTokens(1000),
		)
		if err != nil {
			log.Printf("Error in second request: %v", err)
			return
		}
	}

	// Print the final response
	if len(resp.Choices) > 0 {
		fmt.Printf("Assistant: %s\n", resp.Choices[0].Message.Content)
	}
}

func runWeatherExample(ctx context.Context, client *openrouter.Client) {
	// Define the weather tool with enum parameter
	tools := []openrouter.Tool{
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "get_current_weather",
				Description: "Get the current weather for a specific location",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "City name, zip code, or coordinates",
						},
						"unit": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"celsius", "fahrenheit"},
							"description": "Temperature unit preference",
							"default":     "celsius",
						},
					},
					"required": []string{"location"},
				},
			},
		},
	}

	messages := []openrouter.Message{
		{
			Role:    "user",
			Content: "What's the weather like in San Francisco? Use fahrenheit.",
		},
	}

	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o-mini"),
		openrouter.WithTools(tools...),
		openrouter.WithMaxTokens(500),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Process tool calls if any
	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		messages = append(messages, resp.Choices[0].Message)

		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			var args struct {
				Location string `json:"location"`
				Unit     string `json:"unit"`
			}
			json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

			result, _ := getCurrentWeather(args.Location, args.Unit)

			messages = append(messages, openrouter.Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: toolCall.ID,
			})
		}

		// Final request with tool results
		resp, _ = client.ChatComplete(ctx, messages,
			openrouter.WithModel("openai/gpt-4o-mini"),
			openrouter.WithTools(tools...),
			openrouter.WithMaxTokens(500),
		)
	}

	if len(resp.Choices) > 0 {
		fmt.Printf("Assistant: %s\n", resp.Choices[0].Message.Content)
	}
}

func runMultipleToolsExample(ctx context.Context, client *openrouter.Client) {
	// Define multiple tools
	tools := []openrouter.Tool{
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "search_gutenberg_books",
				Description: "Search for books in the Project Gutenberg library",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"search_terms": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
							"description": "List of search terms",
						},
					},
					"required": []string{"search_terms"},
				},
			},
		},
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "get_current_weather",
				Description: "Get the current weather for a location",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "City name",
						},
					},
					"required": []string{"location"},
				},
			},
		},
	}

	messages := []openrouter.Message{
		{
			Role:    "user",
			Content: "What's the weather in Dublin, the city James Joyce wrote about?",
		},
	}

	// Disable parallel tool calls to ensure sequential processing
	parallelCalls := false
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("anthropic/claude-3.5-sonnet"),
		openrouter.WithTools(tools...),
		openrouter.WithParallelToolCalls(&parallelCalls),
		openrouter.WithMaxTokens(1000),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// The model might call multiple tools
	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		messages = append(messages, resp.Choices[0].Message)

		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			fmt.Printf("Tool called: %s\n", toolCall.Function.Name)

			var result string
			switch toolCall.Function.Name {
			case "search_gutenberg_books":
				var args struct {
					SearchTerms []string `json:"search_terms"`
				}
				json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
				result, _ = searchGutenbergBooks(args.SearchTerms)

			case "get_current_weather":
				var args struct {
					Location string `json:"location"`
				}
				json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
				result, _ = getCurrentWeather(args.Location, "celsius")
			}

			messages = append(messages, openrouter.Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: toolCall.ID,
			})
		}

		// Final response
		resp, _ = client.ChatComplete(ctx, messages,
			openrouter.WithModel("anthropic/claude-3.5-sonnet"),
			openrouter.WithTools(tools...),
			openrouter.WithMaxTokens(1000),
		)
	}

	if len(resp.Choices) > 0 {
		fmt.Printf("Assistant: %s\n", resp.Choices[0].Message.Content)
	}
}

func runForcedToolChoiceExample(ctx context.Context, client *openrouter.Client) {
	tools := []openrouter.Tool{
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "get_current_weather",
				Description: "Get weather information",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"location"},
				},
			},
		},
	}

	messages := []openrouter.Message{
		{
			Role:    "user",
			Content: "I'm thinking about visiting Paris.",
		},
	}

	// Force the model to use the weather tool
	toolChoice := map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name": "get_current_weather",
		},
	}

	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithTools(tools...),
		openrouter.WithToolChoice(toolChoice),
		openrouter.WithMaxTokens(500),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Process the forced tool call
	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		messages = append(messages, resp.Choices[0].Message)

		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			fmt.Printf("Forced tool call for: %s\n", toolCall.Function.Arguments)

			var args struct {
				Location string `json:"location"`
			}
			json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

			result, _ := getCurrentWeather(args.Location, "celsius")

			messages = append(messages, openrouter.Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: toolCall.ID,
			})
		}

		// Final response
		resp, _ = client.ChatComplete(ctx, messages,
			openrouter.WithModel("openai/gpt-4o"),
			openrouter.WithTools(tools...),
			openrouter.WithMaxTokens(500),
		)
	}

	if len(resp.Choices) > 0 {
		fmt.Printf("Assistant: %s\n", resp.Choices[0].Message.Content)
	}
}