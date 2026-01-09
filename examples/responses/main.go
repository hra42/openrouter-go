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
	"time"

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

	// Create a shared context with timeout for all examples
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Example 1: Basic string input
	fmt.Println("=== Example 1: Basic String Input ===")
	basicStringInput(ctx, client)

	// Example 2: Structured message input
	fmt.Println("\n=== Example 2: Structured Message Input ===")
	structuredInput(ctx, client)

	// Example 3: With reasoning
	fmt.Println("\n=== Example 3: With Reasoning ===")
	withReasoning(ctx, client)

	// Example 4: Tool calling
	fmt.Println("\n=== Example 4: Tool Calling ===")
	toolCalling(ctx, client)

	// Example 5: Web search
	fmt.Println("\n=== Example 5: Web Search ===")
	webSearch(ctx, client)

	// Example 6: Streaming
	fmt.Println("\n=== Example 6: Streaming ===")
	streaming(ctx, client)
}

func basicStringInput(ctx context.Context, client *openrouter.Client) {
	// Simple string input - the most basic way to use the Responses API
	resp, err := client.CreateResponse(
		ctx,
		"What is the capital of France? Reply in one word.",
		openrouter.WithResponsesModel("openai/gpt-4o-mini"),
		openrouter.WithResponsesMaxOutputTokens(100),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp == nil {
		log.Printf("Received nil response")
		return
	}

	if text := resp.GetTextContent(); text != "" {
		fmt.Printf("Response: %s\n", text)
	} else {
		fmt.Println("No text content in response")
	}

	if resp.Usage.InputTokens > 0 || resp.Usage.OutputTokens > 0 || resp.Usage.TotalTokens > 0 {
		fmt.Printf("Tokens used: %d input, %d output, %d total\n",
			resp.Usage.InputTokens, resp.Usage.OutputTokens, resp.Usage.TotalTokens)
	} else {
		fmt.Println("No usage information available")
	}
}

func structuredInput(ctx context.Context, client *openrouter.Client) {
	// Use structured input for multi-turn conversations or more control
	input := []openrouter.ResponsesInputItem{
		openrouter.CreateResponsesSystemMessage("You are a helpful math tutor. Explain concepts clearly and concisely."),
		openrouter.CreateResponsesUserMessage("What is 15 + 27?"),
	}

	resp, err := client.CreateResponse(
		ctx,
		input,
		openrouter.WithResponsesModel("openai/gpt-4o-mini"),
		openrouter.WithResponsesMaxOutputTokens(200),
		openrouter.WithResponsesTemperature(0.7),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp == nil || len(resp.Output) == 0 {
		log.Printf("No output in response")
		return
	}

	textContent := resp.GetTextContent()
	if textContent == "" {
		log.Printf("Empty text content in response")
		return
	}

	fmt.Printf("Response: %s\n", textContent)

	// Continue the conversation
	input = append(input, openrouter.ResponsesInputItem{
		Type:   "message",
		ID:     resp.Output[0].ID,
		Status: "completed",
		Role:   "assistant",
		Content: []openrouter.ResponsesInputContent{
			{Type: "input_text", Text: textContent},
		},
	})
	input = append(input, openrouter.CreateResponsesUserMessage("Now multiply that by 2"))

	resp2, err := client.CreateResponse(
		ctx,
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

func withReasoning(ctx context.Context, client *openrouter.Client) {
	// Enable reasoning for complex problems
	// Note: Best used with models that support reasoning like openai/gpt-4o-mini
	resp, err := client.CreateResponse(
		ctx,
		"A farmer has 17 sheep. All but 9 run away. How many sheep does the farmer have left?",
		openrouter.WithResponsesModel("openai/gpt-4o-mini"), // Use gpt-4o-mini for reasoning
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

func toolCalling(ctx context.Context, client *openrouter.Client) {
	// Define a tool that the model can call
	// Responses API uses a flat tool structure (name, description, parameters at top level)
	weatherTool := openrouter.CreateResponsesTool(
		"get_weather",
		"Get the current weather in a given location",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"location": map[string]any{
					"type":        "string",
					"description": "The city and state, e.g. San Francisco, CA",
				},
				"unit": map[string]any{
					"type":        "string",
					"enum":        []string{"celsius", "fahrenheit"},
					"description": "The temperature unit",
				},
			},
			"required": []string{"location"},
		},
	)

	resp, err := client.CreateResponse(
		ctx,
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

		// In a real application, you would execute the function and send the result back.
		// Use CreateResponsesFunctionOutput to send the function result:
		//
		//   input := []openrouter.ResponsesInputItem{
		//       openrouter.CreateResponsesUserMessage("What's the weather like in Tokyo?"),
		//       // Use "function_call_output" type with CallID and Output fields
		//       openrouter.CreateResponsesFunctionOutput(calls[0].CallID, `{"temperature": 22, "unit": "celsius", "condition": "sunny"}`),
		//   }
		//   resp2, err := client.CreateResponse(ctx, input, ...)
	} else {
		// Model responded with text instead
		fmt.Printf("Response: %s\n", resp.GetTextContent())
	}
}

func webSearch(ctx context.Context, client *openrouter.Client) {
	// Enable web search to get up-to-date information
	resp, err := client.CreateResponse(
		ctx,
		"What are recent developments in artificial intelligence? Please cite sources for each claim.",
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

func streaming(ctx context.Context, client *openrouter.Client) {
	// Use streaming for real-time responses
	stream, err := client.CreateResponseStream(
		ctx,
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
		// Only print new content (delta)
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
