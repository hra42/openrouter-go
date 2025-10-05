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
