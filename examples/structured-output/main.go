package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	client := openrouter.NewClient(apiKey)

	// Example 1: Simple structured output for weather information
	weatherExample(client)

	// Example 2: Complex structured output for function calling
	functionCallingExample(client)

	// Example 3: Streaming with structured output
	streamingExample(client)
}

func weatherExample(client *openrouter.Client) {
	fmt.Println("=== Weather Information Example ===")

	// Define the JSON schema for weather information
	weatherSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "City or location name",
			},
			"temperature": map[string]interface{}{
				"type":        "number",
				"description": "Temperature in Celsius",
			},
			"conditions": map[string]interface{}{
				"type":        "string",
				"description": "Weather conditions description",
			},
			"humidity": map[string]interface{}{
				"type":        "number",
				"description": "Humidity percentage",
			},
			"windSpeed": map[string]interface{}{
				"type":        "number",
				"description": "Wind speed in km/h",
			},
		},
		"required":             []string{"location", "temperature", "conditions"},
		"additionalProperties": false,
	}

	messages := []openrouter.Message{
		{
			Role:    "user",
			Content: "What's the weather like in London today? Please provide detailed information.",
		},
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithJSONSchema("weather", true, weatherSchema),
		openrouter.WithRequireParameters(true), // Ensure we only use models that support structured outputs
	)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Parse the JSON response
	var weather map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content.(string)), &weather); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	fmt.Printf("Weather Response: %+v\n\n", weather)
}

func functionCallingExample(client *openrouter.Client) {
	fmt.Println("=== Function Calling Example ===")

	// Define a schema for extracting function calls from natural language
	functionSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"search", "calculate", "translate", "summarize"},
				"description": "The type of action to perform",
			},
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The main input for the action",
					},
					"options": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"language": map[string]interface{}{
								"type":        "string",
								"description": "Target language for translation",
							},
							"limit": map[string]interface{}{
								"type":        "integer",
								"description": "Maximum number of results",
							},
							"format": map[string]interface{}{
								"type":        "string",
								"enum":        []string{"brief", "detailed", "bullets"},
								"description": "Output format preference",
							},
						},
					},
				},
				"required": []string{"query"},
			},
			"confidence": map[string]interface{}{
				"type":        "number",
				"description": "Confidence level in the interpretation (0-1)",
				"minimum":     0,
				"maximum":     1,
			},
		},
		"required":             []string{"action", "parameters"},
		"additionalProperties": false,
	}

	messages := []openrouter.Message{
		{
			Role:    "user",
			Content: "Can you help me translate 'Hello, how are you?' to French?",
		},
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithJSONSchema("function_call", true, functionSchema),
	)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Parse the structured response
	var functionCall map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content.(string)), &functionCall); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	fmt.Printf("Function Call: %+v\n\n", functionCall)
}

func streamingExample(client *openrouter.Client) {
	fmt.Println("=== Streaming with Structured Output ===")

	// Define a simple schema for a task list
	taskSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"tasks": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "integer",
							"description": "Task ID",
						},
						"title": map[string]interface{}{
							"type":        "string",
							"description": "Task title",
						},
						"priority": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"low", "medium", "high"},
							"description": "Task priority",
						},
						"completed": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether the task is completed",
						},
					},
					"required": []string{"id", "title", "priority"},
				},
				"description": "List of tasks",
			},
			"summary": map[string]interface{}{
				"type":        "string",
				"description": "Brief summary of the task list",
			},
		},
		"required":             []string{"tasks", "summary"},
		"additionalProperties": false,
	}

	messages := []openrouter.Message{
		{
			Role:    "user",
			Content: "Create a task list for planning a small birthday party.",
		},
	}

	stream, err := client.ChatCompleteStream(
		context.Background(),
		messages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithJSONSchema("task_list", true, taskSchema),
	)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Streaming response:")
	var fullContent string

	for event := range stream.Events() {
		if len(event.Choices) > 0 && event.Choices[0].Delta != nil {
			if content, ok := event.Choices[0].Delta.Content.(string); ok {
				fmt.Print(content)
				fullContent += content
			}
		}
	}

	if err := stream.Err(); err != nil {
		log.Printf("Stream error: %v", err)
		return
	}

	// Parse the complete JSON
	fmt.Println("\n\nParsed tasks:")
	var taskList map[string]interface{}
	if err := json.Unmarshal([]byte(fullContent), &taskList); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	// Pretty print the parsed JSON
	prettyJSON, _ := json.MarshalIndent(taskList, "", "  ")
	fmt.Println(string(prettyJSON))
}