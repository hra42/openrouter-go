// Package main demonstrates a simpler streaming with tool calling example.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// calculateExpression evaluates a simple mathematical expression
func calculateExpression(operation string, a, b float64) (string, error) {
	var result float64
	switch operation {
	case "add":
		result = a + b
	case "multiply":
		result = a * b
	case "subtract":
		result = a - b
	case "divide":
		if b != 0 {
			result = a / b
		} else {
			return "", fmt.Errorf("division by zero")
		}
	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}

	response := map[string]interface{}{
		"operation": operation,
		"a":         a,
		"b":         b,
		"result":    result,
	}

	jsonResult, err := json.Marshal(response)
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

	client := openrouter.NewClient(
		openrouter.WithAPIKey(apiKey),
		openrouter.WithTimeout(60*time.Second),
	)

	ctx := context.Background()

	// Define a calculator tool
	tools := []openrouter.Tool{
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "calculate",
				Description: "Perform mathematical calculations",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"operation": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"add", "subtract", "multiply", "divide"},
							"description": "Mathematical operation to perform",
						},
						"a": map[string]interface{}{
							"type":        "number",
							"description": "First number",
						},
						"b": map[string]interface{}{
							"type":        "number",
							"description": "Second number",
						},
					},
					"required": []string{"operation", "a", "b"},
				},
			},
		},
	}

	// Initial message
	messages := []openrouter.Message{
		{
			Role:    "user",
			Content: "What's 25 multiplied by 17?",
		},
	}

	fmt.Println("=== Streaming Tool Call Example ===")
	fmt.Println("\nUser: What's 25 multiplied by 17?")
	fmt.Println("\nStreaming response...")

	// Start streaming request
	stream, err := client.ChatCompleteStream(ctx, messages,
		openrouter.WithModel("openai/gpt-4o-mini"),
		openrouter.WithTools(tools...),
		openrouter.WithMaxTokens(100),
	)
	if err != nil {
		log.Fatalf("Error starting stream: %v", err)
	}
	defer stream.Close()

	// Collect the streaming response
	var fullContent strings.Builder
	toolCallsMap := make(map[string]*openrouter.ToolCall)
	currentToolCallID := "" // Track current tool call for argument accumulation
	hasToolCalls := false

	for event := range stream.Events() {
		for _, choice := range event.Choices {
			// Check for content
			if choice.Delta != nil && choice.Delta.Content != nil {
				if content, ok := choice.Delta.Content.(string); ok {
					fullContent.WriteString(content)
					fmt.Print(content)
				}
			}

			// Check for tool calls in delta
			if choice.Delta != nil && len(choice.Delta.ToolCalls) > 0 {
				hasToolCalls = true
				// Process each tool call delta
				for _, tc := range choice.Delta.ToolCalls {
					// If we have an ID, this is either a new tool call or updating an existing one
					if tc.ID != "" {
						currentToolCallID = tc.ID
						if existing, ok := toolCallsMap[tc.ID]; ok {
							// Update existing tool call
							if tc.Type != "" {
								existing.Type = tc.Type
							}
							if tc.Function.Name != "" {
								existing.Function.Name = tc.Function.Name
							}
							if tc.Function.Arguments != "" {
								existing.Function.Arguments += tc.Function.Arguments
							}
						} else {
							// Create new tool call
							newTC := tc
							toolCallsMap[tc.ID] = &newTC
						}
					} else if currentToolCallID != "" {
						// No ID means we're appending arguments to the current tool call
						if existing, ok := toolCallsMap[currentToolCallID]; ok {
							existing.Function.Arguments += tc.Function.Arguments
						}
					}
				}
			}

			// Check finish reason
			if choice.FinishReason == "tool_calls" {
				hasToolCalls = true
			}
		}
	}

	// Convert map to slice
	var toolCalls []openrouter.ToolCall
	for _, tc := range toolCallsMap {
		if tc.Function.Name != "" { // Only include valid tool calls
			toolCalls = append(toolCalls, *tc)
		}
	}

	// Check for errors
	if err := stream.Err(); err != nil {
		log.Fatalf("Stream error: %v", err)
	}

	// If we have tool calls, execute them
	if hasToolCalls && len(toolCalls) > 0 {
		fmt.Println("\n\nModel requested tool calls:")

		// Create assistant message with tool calls
		assistantMsg := openrouter.Message{
			Role:      "assistant",
			Content:   fullContent.String(),
			ToolCalls: toolCalls,
		}
		messages = append(messages, assistantMsg)

		// Execute each tool call
		for _, toolCall := range toolCalls {
			fmt.Printf("- Tool: %s (ID: %s)\n", toolCall.Function.Name, toolCall.ID)
			fmt.Printf("  Args: %s\n", toolCall.Function.Arguments)

			// Skip if arguments are empty
			if toolCall.Function.Arguments == "" {
				log.Printf("Skipping tool call with empty arguments")
				continue
			}

			// Parse arguments
			var args struct {
				Operation string  `json:"operation"`
				A         float64 `json:"a"`
				B         float64 `json:"b"`
			}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				log.Printf("Error parsing arguments: %v (args: %s)", err, toolCall.Function.Arguments)
				continue
			}

			// Execute tool
			result, err := calculateExpression(args.Operation, args.A, args.B)
			if err != nil {
				result = fmt.Sprintf(`{"error": "%s"}`, err.Error())
			}

			fmt.Printf("  Result: %s\n", result)

			// Add tool result to messages
			messages = append(messages, openrouter.Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: toolCall.ID,
			})
		}

		// Get final response with tool results
		fmt.Println("\nFinal response:")
		finalStream, err := client.ChatCompleteStream(ctx, messages,
			openrouter.WithModel("openai/gpt-4o-mini"),
			openrouter.WithTools(tools...),
			openrouter.WithMaxTokens(200),
		)
		if err != nil {
			log.Fatalf("Error in final stream: %v", err)
		}
		defer finalStream.Close()

		// Stream the final response
		for event := range finalStream.Events() {
			for _, choice := range event.Choices {
				if choice.Delta != nil && choice.Delta.Content != nil {
					if content, ok := choice.Delta.Content.(string); ok {
						fmt.Print(content)
					}
				}
			}
		}

		if err := finalStream.Err(); err != nil {
			log.Printf("Final stream error: %v", err)
		}
	} else if fullContent.Len() > 0 {
		// Model responded directly without tool calls
		fmt.Println("\n\nModel responded directly without using tools.")
	}

	fmt.Println("\n\n--- End of Example ---")
}