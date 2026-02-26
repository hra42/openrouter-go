package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunMCPToolConversionTest tests MCP tool conversion and usage with the API.
func RunMCPToolConversionTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: MCP Tool Conversion\n")

	// Test 1: Basic MCP tool conversion
	fmt.Printf("   Testing basic MCP tool conversion...\n")

	mcpTool := openrouter.MCPTool{
		Name:        "calculate",
		Description: "Perform a mathematical calculation",
		InputSchema: &openrouter.MCPInputSchema{
			Type: "object",
			Properties: map[string]any{
				"expression": map[string]any{
					"type":        "string",
					"description": "Mathematical expression to evaluate",
				},
			},
			Required: []string{"expression"},
		},
	}

	tool := openrouter.ConvertMCPTool(mcpTool)

	if tool.Type != "function" {
		printError("Expected tool type 'function'", fmt.Errorf("got %q", tool.Type))
		return false
	}
	if tool.Function.Name != "calculate" {
		printError("Expected function name 'calculate'", fmt.Errorf("got %q", tool.Function.Name))
		return false
	}

	fmt.Printf("   ✅ Conversion successful\n")

	// Test 2: Use converted MCP tool with API
	fmt.Printf("   Testing MCP tool with API...\n")

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What is 15 times 7?"),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithTools(tool),
		openrouter.WithMaxTokens(100),
	)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed API request with MCP tool", err)
		return false
	}

	if len(resp.Choices) == 0 {
		printError("No choices in response", nil)
		return false
	}

	fmt.Printf("   ✅ API request successful (%.2fs)\n", elapsed.Seconds())

	if verbose {
		if len(resp.Choices[0].Message.ToolCalls) > 0 {
			fmt.Printf("      Tool call requested: %s\n", resp.Choices[0].Message.ToolCalls[0].Function.Name)
		} else if content, ok := resp.Choices[0].Message.Content.(string); ok {
			fmt.Printf("      Response: %s\n", truncateString(content, 100))
		}
	}

	// Test 3: MCP tool result conversion
	fmt.Printf("   Testing MCP tool result conversion...\n")

	mcpResult := openrouter.MCPToolResult{
		Content: []openrouter.MCPContent{
			{Type: "text", Text: "The result is 105"},
		},
	}

	resultStr := openrouter.ConvertToolResultToMCP(mcpResult)
	if resultStr != "The result is 105" {
		printError("Unexpected result string", fmt.Errorf("got %q", resultStr))
		return false
	}

	fmt.Printf("   ✅ Result conversion successful\n")

	// Test 4: Multiple MCP tools conversion
	fmt.Printf("   Testing multiple MCP tools conversion...\n")

	mcpTools := []openrouter.MCPTool{
		{Name: "tool1", Description: "First tool"},
		{Name: "tool2", Description: "Second tool"},
		{Name: "tool3", Description: "Third tool"},
	}

	tools := openrouter.ConvertMCPTools(mcpTools)
	if len(tools) != 3 {
		printError("Expected 3 tools", fmt.Errorf("got %d", len(tools)))
		return false
	}

	fmt.Printf("   ✅ Multiple tools conversion successful\n")

	// Test 5: Parse MCP tools from JSON
	fmt.Printf("   Testing MCP tools JSON parsing...\n")

	jsonData := []byte(`[{"name": "test_tool", "description": "A test tool"}]`)
	parsedTools, err := openrouter.ParseMCPToolsFromJSON(jsonData)
	if err != nil {
		printError("Failed to parse MCP tools from JSON", err)
		return false
	}
	if len(parsedTools) != 1 || parsedTools[0].Name != "test_tool" {
		printError("Unexpected parsed tool", fmt.Errorf("got %+v", parsedTools))
		return false
	}

	fmt.Printf("   ✅ JSON parsing successful\n")

	// Test 6: Complex MCP tool with full schema
	fmt.Printf("   Testing complex MCP tool with full schema...\n")

	complexMCPTool := openrouter.MCPTool{
		Name:        "search_files",
		Description: "Search for files in a directory",
		InputSchema: &openrouter.MCPInputSchema{
			Type: "object",
			Properties: map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Directory path to search",
				},
				"pattern": map[string]any{
					"type":        "string",
					"description": "Search pattern",
				},
				"recursive": map[string]any{
					"type":        "boolean",
					"description": "Whether to search recursively",
				},
			},
			Required: []string{"path"},
		},
	}

	complexTool := openrouter.ConvertMCPTool(complexMCPTool)

	// Verify the parameters are correctly converted
	params := complexTool.Function.Parameters
	if params["type"] != "object" {
		printError("Expected parameters type 'object'", fmt.Errorf("got %v", params["type"]))
		return false
	}

	props, ok := params["properties"].(map[string]any)
	if !ok {
		printError("Expected properties to be a map", nil)
		return false
	}

	if _, exists := props["path"]; !exists {
		printError("Expected 'path' property", nil)
		return false
	}

	fmt.Printf("   ✅ Complex tool conversion successful\n")

	if verbose {
		printVerbose(true, "All MCP conversion tests passed")
	}

	printSuccess("MCP tool conversion tests completed")
	return true
}
