package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunToolCallingTest tests tool/function calling functionality including basic tool calls,
// multiple tools, parallel calls, and streaming with tools.
func RunToolCallingTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Tool/Function Calling\n")

	// Test 1: Basic tool calling
	fmt.Printf("   Testing basic tool calling...\n")

	// Define a simple calculator tool
	tools := []openrouter.Tool{
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "calculate",
				Description: "Perform basic mathematical calculations",
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
							"description": "First operand",
						},
						"b": map[string]interface{}{
							"type":        "number",
							"description": "Second operand",
						},
					},
					"required": []string{"operation", "a", "b"},
				},
			},
		},
	}

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What is 15 multiplied by 7?"),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithTools(tools...),
		openrouter.WithMaxTokens(100),
	)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed initial request", err)
		return false
	}

	// Check if the model requested tool calls
	if len(resp.Choices) == 0 || len(resp.Choices[0].Message.ToolCalls) == 0 {
		fmt.Printf("   ❌ Model didn't request any tool calls\n")
		return false
	}

	fmt.Printf("   ✅ Tool request received (%.2fs)\n", elapsed.Seconds())

	// Process tool calls
	toolCall := resp.Choices[0].Message.ToolCalls[0]
	fmt.Printf("   Tool: %s\n", toolCall.Function.Name)
	fmt.Printf("   Arguments: %s\n", toolCall.Function.Arguments)

	// Parse arguments and simulate tool execution
	var args struct {
		Operation string  `json:"operation"`
		A         float64 `json:"a"`
		B         float64 `json:"b"`
	}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		printError("Failed to parse arguments", err)
		return false
	}

	// Execute the tool (simulated)
	var result float64
	switch args.Operation {
	case "add":
		result = args.A + args.B
	case "subtract":
		result = args.A - args.B
	case "multiply":
		result = args.A * args.B
	case "divide":
		if args.B != 0 {
			result = args.A / args.B
		} else {
			result = 0
		}
	}

	toolResult := fmt.Sprintf(`{"result": %f}`, result)
	fmt.Printf("   Tool result: %s\n", toolResult)

	// Add tool response to messages
	messages = append(messages, resp.Choices[0].Message)
	messages = append(messages, openrouter.Message{
		Role:       "tool",
		Content:    toolResult,
		ToolCallID: toolCall.ID,
	})

	// Get final response
	start = time.Now()
	finalResp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithTools(tools...),
		openrouter.WithMaxTokens(100),
	)
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed final request", err)
		return false
	}

	fmt.Printf("   ✅ Final response received (%.2fs)\n", elapsed.Seconds())
	if verbose {
		fmt.Printf("   Final answer: %s\n", strings.TrimSpace(finalResp.Choices[0].Message.Content.(string)))
	}

	// Test 2: Multiple tools and tool choice
	fmt.Printf("   Testing multiple tools and tool choice...\n")

	multiTools := []openrouter.Tool{
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "get_weather",
				Description: "Get current weather for a location",
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
		{
			Type: "function",
			Function: openrouter.Function{
				Name:        "get_time",
				Description: "Get current time for a timezone",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"timezone": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"timezone"},
				},
			},
		},
	}

	// Test with auto tool choice (default)
	weatherMessages := []openrouter.Message{
		openrouter.CreateUserMessage("What's the weather in Tokyo?"),
	}

	resp, err = client.ChatComplete(ctx, weatherMessages,
		openrouter.WithModel(model),
		openrouter.WithTools(multiTools...),
		openrouter.WithToolChoice("auto"),
		openrouter.WithMaxTokens(100),
	)

	if err != nil {
		fmt.Printf("   ⚠️  Multi-tool test failed: %v\n", err)
	} else if len(resp.Choices[0].Message.ToolCalls) > 0 {
		fmt.Printf("   ✅ Auto tool choice worked\n")
		fmt.Printf("      Tool selected: %s\n", resp.Choices[0].Message.ToolCalls[0].Function.Name)
	}

	// Test 3: Parallel tool calls
	fmt.Printf("   Testing parallel tool calls control...\n")

	parallelMessages := []openrouter.Message{
		openrouter.CreateUserMessage("What's the weather in Paris and the time in New York?"),
	}

	// Test with parallel tool calls disabled
	parallelCalls := false
	resp, err = client.ChatComplete(ctx, parallelMessages,
		openrouter.WithModel(model),
		openrouter.WithTools(multiTools...),
		openrouter.WithParallelToolCalls(&parallelCalls),
		openrouter.WithMaxTokens(100),
	)

	if err != nil {
		fmt.Printf("   ⚠️  Parallel tool calls test failed: %v\n", err)
	} else {
		toolCount := len(resp.Choices[0].Message.ToolCalls)
		fmt.Printf("   ✅ Parallel control test completed\n")
		fmt.Printf("      Tools requested: %d (parallel_tool_calls=false)\n", toolCount)
		if verbose && toolCount > 0 {
			for i, tc := range resp.Choices[0].Message.ToolCalls {
				fmt.Printf("      Tool %d: %s\n", i+1, tc.Function.Name)
			}
		}
	}

	// Test 4: Streaming with tool calls
	fmt.Printf("   Testing streaming with tool calls...\n")

	streamMessages := []openrouter.Message{
		openrouter.CreateUserMessage("Calculate 25 plus 17"),
	}

	stream, err := client.ChatCompleteStream(ctx, streamMessages,
		openrouter.WithModel(model),
		openrouter.WithTools(tools...),
		openrouter.WithMaxTokens(100),
	)

	if err != nil {
		printError("Failed to create stream", err)
		return false
	}
	defer stream.Close()

	hasToolCalls := false
	eventCount := 0

	for event := range stream.Events() {
		eventCount++
		// Check for tool calls in the streaming response
		for _, choice := range event.Choices {
			// Check if delta contains tool calls
			if choice.Delta != nil && len(choice.Delta.ToolCalls) > 0 {
				hasToolCalls = true
			}
			// Check finish reason
			if choice.FinishReason == "tool_calls" {
				hasToolCalls = true
			}
		}
	}

	if err := stream.Err(); err != nil {
		printError("Stream error", err)
		return false
	}

	if hasToolCalls {
		fmt.Printf("   ✅ Streaming with tool calls worked (%d events)\n", eventCount)
	} else {
		fmt.Printf("   ⚠️  No tool calls in stream (model may have answered directly)\n")
	}

	printSuccess("Tool calling tests completed")
	return true
}
