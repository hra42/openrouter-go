// Responses API Example
//
// WARNING: The Responses API is in BETA and may have breaking changes at any time.
// Do not rely on this API for production workloads. The API structure, parameters,
// and behavior may change without notice as OpenRouter continues to develop this feature.
//
// For stable production use, consider using the Chat Completions API instead.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	// WARNING: Responses API is in BETA - expect breaking changes!
	fmt.Println("========================================")
	fmt.Println("WARNING: Responses API is in BETA")
	fmt.Println("This API may have breaking changes at any time.")
	fmt.Println("Do not use in production workloads.")
	fmt.Println("========================================")
	fmt.Println()

	// Get API key from environment variable
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	// Create a new client
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Example 1: Basic string input
	fmt.Println("=== Example 1: Basic String Input ===")
	basicStringInput(client)

	// Example 2: Structured message input
	fmt.Println("\n=== Example 2: Structured Message Input ===")
	structuredInput(client)

	// Example 3: With reasoning
	fmt.Println("\n=== Example 3: With Reasoning ===")
	withReasoning(client)

	// Example 4: Tool calling
	fmt.Println("\n=== Example 4: Tool Calling ===")
	toolCalling(client)

	// Example 5: Web search
	fmt.Println("\n=== Example 5: Web Search ===")
	webSearch(client)

	// Example 6: Streaming
	fmt.Println("\n=== Example 6: Streaming ===")
	streaming(client)
}

func basicStringInput(client *openrouter.Client) {
	// Simple string input - the most basic way to use the Responses API
	resp, err := client.CreateResponse(
		context.Background(),
		"What is the capital of France? Reply in one word.",
		openrouter.WithResponsesModel("openai/gpt-4o-mini"),
		openrouter.WithResponsesMaxOutputTokens(100),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.GetTextContent())
	fmt.Printf("Tokens used: %d input, %d output, %d total\n",
		resp.Usage.InputTokens, resp.Usage.OutputTokens, resp.Usage.TotalTokens)
}

func structuredInput(client *openrouter.Client) {
	// Use structured input for multi-turn conversations or more control
	input := []openrouter.ResponsesInputItem{
		openrouter.CreateResponsesSystemMessage("You are a helpful math tutor. Explain concepts clearly and concisely."),
		openrouter.CreateResponsesUserMessage("What is 15 + 27?"),
	}

	resp, err := client.CreateResponse(
		context.Background(),
		input,
		openrouter.WithResponsesModel("openai/gpt-4o-mini"),
		openrouter.WithResponsesMaxOutputTokens(200),
		openrouter.WithResponsesTemperature(0.7),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.GetTextContent())

	// Continue the conversation
	input = append(input, openrouter.ResponsesInputItem{
		Type:   "message",
		ID:     resp.Output[0].ID,
		Status: "completed",
		Role:   "assistant",
		Content: []openrouter.ResponsesInputContent{
			{Type: "input_text", Text: resp.GetTextContent()},
		},
	})
	input = append(input, openrouter.CreateResponsesUserMessage("Now multiply that by 2"))

	resp2, err := client.CreateResponse(
		context.Background(),
		input,
		openrouter.WithResponsesModel("openai/gpt-4o-mini"),
		openrouter.WithResponsesMaxOutputTokens(200),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Follow-up response: %s\n", resp2.GetTextContent())
}

func withReasoning(client *openrouter.Client) {
	// Enable reasoning for complex problems
	// Note: Best used with models that support reasoning like openai/o4-mini
	resp, err := client.CreateResponse(
		context.Background(),
		"A farmer has 17 sheep. All but 9 run away. How many sheep does the farmer have left?",
		openrouter.WithResponsesModel("openai/gpt-4o-mini"), // Use o4-mini for better reasoning
		openrouter.WithResponsesMaxOutputTokens(500),
		openrouter.WithResponsesReasoningEffort(openrouter.ReasoningEffortMedium),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.GetTextContent())

	// Check for reasoning summary if available
	if summary := resp.GetReasoningSummary(); len(summary) > 0 {
		fmt.Println("Reasoning steps:")
		for i, step := range summary {
			fmt.Printf("  %d. %s\n", i+1, step)
		}
	}

	if resp.Usage.ReasoningTokens > 0 {
		fmt.Printf("Reasoning tokens used: %d\n", resp.Usage.ReasoningTokens)
	}
}

func toolCalling(client *openrouter.Client) {
	// Define a tool that the model can call
	weatherTool := openrouter.Tool{
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
						"type":        "string",
						"enum":        []string{"celsius", "fahrenheit"},
						"description": "The temperature unit",
					},
				},
				"required": []string{"location"},
			},
		},
	}

	resp, err := client.CreateResponse(
		context.Background(),
		"What's the weather like in Tokyo?",
		openrouter.WithResponsesModel("openai/gpt-4o-mini"),
		openrouter.WithResponsesMaxOutputTokens(200),
		openrouter.WithResponsesTools(weatherTool),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Check if we got function calls
	calls := resp.GetFunctionCalls()
	if len(calls) > 0 {
		fmt.Println("Function calls:")
		for _, call := range calls {
			fmt.Printf("  Function: %s\n", call.Name)
			fmt.Printf("  Arguments: %s\n", call.Arguments)
			fmt.Printf("  Call ID: %s\n", call.CallID)
		}

		// In a real application, you would execute the function and send the result back
		// Example of sending function result:
		/*
			input := []openrouter.ResponsesInputItem{
				openrouter.CreateResponsesUserMessage("What's the weather like in Tokyo?"),
				{
					Type:      "function_call",
					ID:        calls[0].ID,
					CallID:    calls[0].CallID,
					Name:      calls[0].Name,
					Arguments: calls[0].Arguments,
				},
				openrouter.CreateResponsesFunctionOutput(calls[0].CallID, `{"temperature": 22, "unit": "celsius", "condition": "sunny"}`),
			}
		*/
	} else {
		// Model responded with text instead
		fmt.Printf("Response: %s\n", resp.GetTextContent())
	}
}

func webSearch(client *openrouter.Client) {
	// Enable web search to get up-to-date information
	resp, err := client.CreateResponse(
		context.Background(),
		"What are the latest developments in AI today?",
		openrouter.WithResponsesModel("openai/gpt-4o-mini"),
		openrouter.WithResponsesMaxOutputTokens(500),
		openrouter.WithResponsesWebSearch(3), // Get up to 3 search results
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.GetTextContent())

	// Check for citations
	annotations := resp.GetAnnotations()
	if len(annotations) > 0 {
		fmt.Println("\nSources:")
		for i, ann := range annotations {
			if ann.Type == "url_citation" {
				fmt.Printf("  [%d] %s\n", i+1, ann.URL)
			}
		}
	}
}

func streaming(client *openrouter.Client) {
	// Use streaming for real-time responses
	stream, err := client.CreateResponseStream(
		context.Background(),
		"Write a haiku about programming.",
		openrouter.WithResponsesModel("openai/gpt-4o-mini"),
		openrouter.WithResponsesMaxOutputTokens(100),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer func() {
		if err := stream.Close(); err != nil {
			log.Printf("Error closing stream: %v", err)
		}
	}()

	fmt.Print("Streaming response: ")
	var lastContent string

	for event := range stream.Events() {
		content := event.GetTextContent()
		// Only print new content
		if len(content) > len(lastContent) {
			fmt.Print(content[len(lastContent):])
			lastContent = content
		}
	}
	fmt.Println()

	if err := stream.Err(); err != nil {
		log.Printf("Stream error: %v", err)
	}
}
