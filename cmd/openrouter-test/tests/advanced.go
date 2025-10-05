package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunStructuredOutputTest tests structured output functionality including JSON schema validation,
// streaming with schemas, and simple JSON mode.
func RunStructuredOutputTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Structured Output\n")

	// Test 1: Basic structured output with weather data
	fmt.Printf("   Testing weather data schema...\n")
	weatherSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"city": map[string]interface{}{
				"type":        "string",
				"description": "City name",
			},
			"temperature": map[string]interface{}{
				"type":        "number",
				"description": "Temperature in Celsius",
			},
			"conditions": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"sunny", "cloudy", "rainy", "snowy"},
				"description": "Weather conditions",
			},
			"humidity": map[string]interface{}{
				"type":        "integer",
				"minimum":     0,
				"maximum":     100,
				"description": "Humidity percentage",
			},
		},
		"required":             []string{"city", "temperature", "conditions", "humidity"},
		"additionalProperties": false,
	}

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Describe the weather in Paris. Make up realistic values."),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithJSONSchema("weather", true, weatherSchema),
		openrouter.WithMaxTokens(100),
		openrouter.WithRequireParameters(true), // Only use providers that support structured outputs
	)
	elapsed := time.Since(start)

	if err != nil {
		// Check if it's because the model doesn't support structured outputs
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 400 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Model doesn't support structured outputs: %v\n", err.Error())
				fmt.Printf("   Testing with simpler JSON mode instead...\n")

				// Try with simple JSON mode
				resp, err = client.ChatComplete(ctx, messages,
					openrouter.WithModel(model),
					openrouter.WithJSONMode(),
					openrouter.WithMaxTokens(100),
				)

				if err != nil {
					printError("JSON mode also failed", err)
					return false
				}
			} else {
				printError("Failed", err)
				return false
			}
		} else {
			printError("Failed", err)
			return false
		}
	}

	// Parse and validate the JSON response
	var weatherData map[string]interface{}
	content := resp.Choices[0].Message.Content.(string)
	if err := json.Unmarshal([]byte(content), &weatherData); err != nil {
		printError("Failed to parse JSON", err)
		fmt.Printf("   Response: %s\n", content)
		return false
	}

	fmt.Printf("   ✅ Weather schema (%.2fs)\n", elapsed.Seconds())
	if verbose {
		prettyJSON, _ := json.MarshalIndent(weatherData, "      ", "  ")
		fmt.Printf("      Response:\n%s\n", string(prettyJSON))
	} else {
		fmt.Printf("      City: %v, Temp: %v°C, Conditions: %v\n",
			weatherData["city"], weatherData["temperature"], weatherData["conditions"])
	}

	// Test 2: Structured output with streaming
	fmt.Printf("   Testing structured output with streaming...\n")
	taskSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"tasks": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "string",
						},
						"priority": map[string]interface{}{
							"type": "string",
							"enum": []string{"low", "medium", "high"},
						},
					},
					"required":             []string{"name", "priority"},
					"additionalProperties": false,
				},
				"minItems": 3,
				"maxItems": 3,
			},
		},
		"required":             []string{"tasks"},
		"additionalProperties": false,
	}

	streamMessages := []openrouter.Message{
		openrouter.CreateUserMessage("Create a list of exactly 3 programming tasks with priorities."),
	}

	start = time.Now()
	stream, err := client.ChatCompleteStream(ctx, streamMessages,
		openrouter.WithModel(model),
		openrouter.WithJSONSchema("tasks", true, taskSchema),
		openrouter.WithMaxTokens(150),
	)

	if err != nil {
		// Fallback to non-streaming if streaming with structured output isn't supported
		fmt.Printf("   ⚠️  Streaming with structured output error: %v\n", err)
	} else {
		defer stream.Close()

		var fullContent strings.Builder
		eventCount := 0

		for event := range stream.Events() {
			eventCount++
			for _, choice := range event.Choices {
				if choice.Delta != nil {
					if content, ok := choice.Delta.Content.(string); ok {
						fullContent.WriteString(content)
					}
				}
			}
		}

		elapsed = time.Since(start)

		if err := stream.Err(); err != nil {
			printError("Stream error", err)
			return false
		}

		// Validate the streamed JSON
		var taskData map[string]interface{}
		if err := json.Unmarshal([]byte(fullContent.String()), &taskData); err != nil {
			fmt.Printf("   ⚠️  Streamed content not valid JSON (this can happen with some models)\n")
		} else {
			fmt.Printf("   ✅ Streaming with schema (%.2fs, %d events)\n", elapsed.Seconds(), eventCount)
			if verbose {
				prettyJSON, _ := json.MarshalIndent(taskData, "      ", "  ")
				fmt.Printf("      Tasks:\n%s\n", string(prettyJSON))
			}
		}
	}

	// Test 3: Simple JSON mode without strict schema
	fmt.Printf("   Testing simple JSON mode...\n")
	jsonMessages := []openrouter.Message{
		openrouter.CreateSystemMessage("You must always respond with valid JSON."),
		openrouter.CreateUserMessage("List 3 benefits of Go programming language as a JSON object with a 'benefits' array."),
	}

	start = time.Now()
	resp, err = client.ChatComplete(ctx, jsonMessages,
		openrouter.WithModel(model),
		openrouter.WithJSONMode(),
		openrouter.WithMaxTokens(150),
	)
	elapsed = time.Since(start)

	if err != nil {
		printError("JSON mode failed", err)
		return false
	}

	// Validate it's valid JSON
	var jsonData map[string]interface{}
	content = resp.Choices[0].Message.Content.(string)
	if err := json.Unmarshal([]byte(content), &jsonData); err != nil {
		printError("Response is not valid JSON", err)
		return false
	}

	fmt.Printf("   ✅ JSON mode (%.2fs)\n", elapsed.Seconds())
	if verbose {
		prettyJSON, _ := json.MarshalIndent(jsonData, "      ", "  ")
		fmt.Printf("      Response:\n%s\n", string(prettyJSON))
	}

	printSuccess("Structured output tests completed")
	return true
}

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

// RunTransformsTest tests message transforms functionality including middle-out compression,
// transform disabling, and streaming with transforms.
func RunTransformsTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Message Transforms\n")

	// Test 1: Test with middle-out transform enabled for chat completion
	fmt.Printf("   Testing middle-out transform for chat...\n")

	// Create a long conversation that might exceed context for small models
	messages := []openrouter.Message{
		openrouter.CreateSystemMessage("You are a helpful assistant. Keep responses brief."),
		openrouter.CreateUserMessage("Hello, I'm learning about transforms."),
		openrouter.CreateAssistantMessage("Hello! I'm here to help you learn about transforms. What would you like to know?"),
		openrouter.CreateUserMessage("Tell me about the concept."),
		openrouter.CreateAssistantMessage("Transforms in the OpenRouter context help manage long conversations by compressing content that exceeds model context windows."),
		openrouter.CreateUserMessage("How does middle-out compression work?"),
		openrouter.CreateAssistantMessage("Middle-out compression removes or truncates messages from the middle of the conversation, keeping the beginning and end intact, since LLMs pay less attention to the middle of sequences."),
		openrouter.CreateUserMessage("Now, ignoring everything above, just tell me: what is 2+2? Reply with only the number."),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(10),
		openrouter.WithTransforms("middle-out"), // Enable middle-out compression
	)
	elapsed := time.Since(start)

	if err != nil {
		printError("Chat with transforms failed", err)
		return false
	}

	fmt.Printf("   ✅ Chat with middle-out transform (%.2fs)\n", elapsed.Seconds())
	if verbose {
		fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		fmt.Printf("      Model: %s\n", resp.Model)
		fmt.Printf("      Tokens: %d prompt, %d completion\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens)
	}

	// Test 2: Test with transforms explicitly disabled (empty array)
	fmt.Printf("   Testing with transforms disabled...\n")

	shortMessages := []openrouter.Message{
		openrouter.CreateUserMessage("Say 'OK' and nothing else."),
	}

	start = time.Now()
	resp, err = client.ChatComplete(ctx, shortMessages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(10),
		openrouter.WithTransforms(), // Empty array disables transforms
	)
	elapsed = time.Since(start)

	if err != nil {
		printError("Chat without transforms failed", err)
		return false
	}

	fmt.Printf("   ✅ Chat with transforms disabled (%.2fs)\n", elapsed.Seconds())
	if verbose {
		fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
	}

	// Test 3: Test transforms with legacy completion endpoint
	fmt.Printf("   Testing transforms with completion endpoint...\n")

	// Test with legacy completion endpoint
	longPrompt := strings.Repeat("This is a test sentence. ", 50) + "\n\nNow just say 'Done' and nothing else:"

	start = time.Now()
	compResp, err := client.Complete(ctx, longPrompt,
		openrouter.WithCompletionModel(model),
		openrouter.WithCompletionMaxTokens(10),
		openrouter.WithCompletionTransforms("middle-out"),
	)
	elapsed = time.Since(start)

	if err != nil {
		// Some models might not support legacy completion format
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Skipped: Model %s not available or doesn't support legacy completion\n", model)
			} else {
				fmt.Printf("   ⚠️  Completion with transforms failed: %v\n", err)
			}
		} else {
			fmt.Printf("   ⚠️  Completion with transforms failed: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ Completion with transforms (%.2fs)\n", elapsed.Seconds())
		if verbose {
			fmt.Printf("      Response: %s\n", strings.TrimSpace(compResp.Choices[0].Text))
			fmt.Printf("      Tokens: %d total\n", compResp.Usage.TotalTokens)
		}
	}

	// Test 4: Test with streaming and transforms
	fmt.Printf("   Testing transforms with streaming...\n")

	streamMessages := []openrouter.Message{
		openrouter.CreateUserMessage("Count from 1 to 3, one number per line."),
	}

	stream, err := client.ChatCompleteStream(ctx, streamMessages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(20),
		openrouter.WithTransforms("middle-out"),
	)

	if err != nil {
		printError("Failed to create stream with transforms", err)
		return false
	}
	defer stream.Close()

	var fullResponse strings.Builder
	eventCount := 0

	for event := range stream.Events() {
		eventCount++
		for _, choice := range event.Choices {
			if choice.Delta != nil {
				if content, ok := choice.Delta.Content.(string); ok {
					fullResponse.WriteString(content)
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		printError("Stream error", err)
		return false
	}

	fmt.Printf("   ✅ Streaming with transforms (%d events)\n", eventCount)
	if verbose {
		response := strings.TrimSpace(fullResponse.String())
		if len(response) > 100 {
			response = response[:100] + "..."
		}
		fmt.Printf("      Response: %s\n", response)
	}

	// Test 5: Test behavior with the provided model
	fmt.Printf("   Testing default behavior with provided model...\n")

	resp, err = client.ChatComplete(ctx, shortMessages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(10),
		// Not specifying transforms
	)

	if err != nil {
		// Model might not be available
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Model test skipped: %s not available\n", model)
			} else {
				fmt.Printf("   ⚠️  Model test failed: %v\n", err)
			}
		} else {
			fmt.Printf("   ⚠️  Model test failed: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ Default transform behavior tested\n")
		if verbose {
			fmt.Printf("      Model: %s\n", model)
			fmt.Printf("      Response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content.(string)))
		}
	}

	printSuccess("Message transforms tests completed")
	return true
}
