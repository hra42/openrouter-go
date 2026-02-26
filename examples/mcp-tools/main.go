// Package main demonstrates using MCP tools with the OpenRouter API.
//
// MCP (Model Context Protocol) is a protocol for providing LLMs with tool calling
// abilities. This example shows how to convert MCP tool definitions to the
// OpenAI-compatible format used by OpenRouter.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hra42/openrouter-go"
)

func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	client := openrouter.NewClient(
		openrouter.WithAPIKey(apiKey),
		openrouter.WithTimeout(30*time.Second),
	)

	ctx := context.Background()

	// Example 1: Converting MCP tools to OpenRouter format
	fmt.Println("=== Example 1: MCP Tool Conversion ===")
	runMCPToolConversionExample(ctx, client)

	// Example 2: Simulating MCP server tool response
	fmt.Println("\n=== Example 2: MCP Tool Result Handling ===")
	runMCPToolResultExample()

	// Example 3: Parsing MCP tools from JSON
	fmt.Println("\n=== Example 3: Parsing MCP Tools from JSON ===")
	runMCPJSONParsingExample(ctx, client)
}

func runMCPToolConversionExample(ctx context.Context, client *openrouter.Client) {
	// Define MCP tools (as you might receive from an MCP server)
	mcpTools := []openrouter.MCPTool{
		{
			Name:        "read_file",
			Description: "Read the contents of a file from the filesystem",
			InputSchema: &openrouter.MCPInputSchema{
				Type: "object",
				Properties: map[string]any{
					"path": map[string]any{
						"type":        "string",
						"description": "The path to the file to read",
					},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "write_file",
			Description: "Write content to a file on the filesystem",
			InputSchema: &openrouter.MCPInputSchema{
				Type: "object",
				Properties: map[string]any{
					"path": map[string]any{
						"type":        "string",
						"description": "The path to the file to write",
					},
					"content": map[string]any{
						"type":        "string",
						"description": "The content to write to the file",
					},
				},
				Required: []string{"path", "content"},
			},
		},
	}

	// Convert MCP tools to OpenRouter format
	tools := openrouter.ConvertMCPTools(mcpTools)

	fmt.Printf("Converted %d MCP tools to OpenRouter format\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("  - %s: %s\n", tool.Function.Name, tool.Function.Description)
	}

	// Use the tools with chat completion
	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What files can you read or write for me?"),
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

	if len(resp.Choices) > 0 {
		fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
	}
}

func runMCPToolResultExample() {
	// Simulate an MCP tool result (as you might receive from an MCP server)
	mcpResult := openrouter.MCPToolResult{
		Content: []openrouter.MCPContent{
			{
				Type: "text",
				Text: "File contents: Hello, World!",
			},
		},
	}

	// Convert to string for tool response
	resultStr := openrouter.ConvertToolResultToMCP(mcpResult)
	fmt.Printf("Tool result as string: %s\n", resultStr)

	// Example with complex result (multiple content items)
	complexResult := openrouter.MCPToolResult{
		Content: []openrouter.MCPContent{
			{Type: "text", Text: "Found 3 files:"},
			{Type: "text", Text: "1. config.json"},
			{Type: "text", Text: "2. data.csv"},
			{Type: "text", Text: "3. readme.md"},
		},
	}

	complexStr := openrouter.ConvertToolResultToMCP(complexResult)
	fmt.Printf("Complex result as JSON: %s\n", complexStr)
}

func runMCPJSONParsingExample(ctx context.Context, client *openrouter.Client) {
	// JSON representation of MCP tools (as received from an MCP server)
	mcpToolsJSON := []byte(`[
		{
			"name": "get_weather",
			"description": "Get the current weather for a location",
			"inputSchema": {
				"type": "object",
				"properties": {
					"location": {
						"type": "string",
						"description": "City name or coordinates"
					},
					"units": {
						"type": "string",
						"enum": ["celsius", "fahrenheit"],
						"description": "Temperature units"
					}
				},
				"required": ["location"]
			}
		}
	]`)

	// Parse MCP tools from JSON
	mcpTools, err := openrouter.ParseMCPToolsFromJSON(mcpToolsJSON)
	if err != nil {
		log.Printf("Failed to parse MCP tools: %v", err)
		return
	}

	fmt.Printf("Parsed %d MCP tools from JSON\n", len(mcpTools))

	// Convert and use with OpenRouter
	tools := openrouter.ConvertMCPTools(mcpTools)

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What's the weather like in Tokyo?"),
	}

	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o-mini"),
		openrouter.WithTools(tools...),
		openrouter.WithMaxTokens(200),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Check for tool calls
	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		for _, tc := range resp.Choices[0].Message.ToolCalls {
			fmt.Printf("Tool call requested: %s\n", tc.Function.Name)
			fmt.Printf("Arguments: %s\n", tc.Function.Arguments)

			// Parse the arguments
			var args map[string]any
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err == nil {
				fmt.Printf("Parsed arguments: %v\n", args)
			}
		}
	} else if len(resp.Choices) > 0 {
		fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
	}
}
