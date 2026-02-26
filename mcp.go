package openrouter

import (
	"encoding/json"
	"fmt"
)

// MCPTool represents a tool in MCP (Model Context Protocol) format.
// This format is used by MCP servers and clients for tool definitions.
type MCPTool struct {
	// Name is the unique identifier for the tool
	Name string `json:"name"`
	// Description explains what the tool does
	Description string `json:"description,omitempty"`
	// InputSchema defines the expected input parameters using JSON Schema
	InputSchema *MCPInputSchema `json:"inputSchema,omitempty"`
}

// MCPInputSchema represents the input schema for an MCP tool.
// It follows JSON Schema Draft 7 format.
type MCPInputSchema struct {
	// Type is typically "object" for MCP tool schemas
	Type string `json:"type,omitempty"`
	// Properties defines the tool's input parameters
	Properties map[string]any `json:"properties,omitempty"`
	// Required lists parameter names that must be provided
	Required []string `json:"required,omitempty"`
	// AdditionalProperties controls whether extra properties are allowed
	AdditionalProperties *bool `json:"additionalProperties,omitempty"`
}

// MCPToolResult represents the result of an MCP tool execution.
type MCPToolResult struct {
	// Content contains the tool execution result(s)
	Content []MCPContent `json:"content"`
	// IsError indicates if the tool execution failed
	IsError bool `json:"isError,omitempty"`
}

// MCPContent represents content in an MCP response.
type MCPContent struct {
	// Type is the content type (e.g., "text", "image", "resource")
	Type string `json:"type"`
	// Text is the text content (when Type is "text")
	Text string `json:"text,omitempty"`
	// Data is base64-encoded data (when Type is "image" or "resource")
	Data string `json:"data,omitempty"`
	// MimeType is the MIME type for binary content
	MimeType string `json:"mimeType,omitempty"`
}

// ConvertMCPTool converts an MCP tool to OpenRouter's Tool format.
// The OpenRouter format follows the OpenAI function calling specification.
func ConvertMCPTool(mcpTool MCPTool) Tool {
	parameters := make(map[string]any)

	if mcpTool.InputSchema != nil {
		// Set the type to "object" as required by OpenAI format
		parameters["type"] = "object"

		// Copy properties
		if mcpTool.InputSchema.Properties != nil {
			parameters["properties"] = mcpTool.InputSchema.Properties
		} else {
			parameters["properties"] = make(map[string]any)
		}

		// Copy required fields
		if mcpTool.InputSchema.Required != nil {
			parameters["required"] = mcpTool.InputSchema.Required
		}

		// Handle additionalProperties
		if mcpTool.InputSchema.AdditionalProperties != nil {
			parameters["additionalProperties"] = *mcpTool.InputSchema.AdditionalProperties
		}
	} else {
		// Empty schema defaults
		parameters["type"] = "object"
		parameters["properties"] = make(map[string]any)
	}

	return Tool{
		Type: "function",
		Function: Function{
			Name:        mcpTool.Name,
			Description: mcpTool.Description,
			Parameters:  parameters,
		},
	}
}

// ConvertMCPTools converts a slice of MCP tools to OpenRouter Tool format.
func ConvertMCPTools(mcpTools []MCPTool) []Tool {
	tools := make([]Tool, len(mcpTools))
	for i, mcpTool := range mcpTools {
		tools[i] = ConvertMCPTool(mcpTool)
	}
	return tools
}

// ConvertToolResultToMCP converts an MCP tool result to a string suitable
// for use as tool response content in chat completions.
func ConvertToolResultToMCP(result MCPToolResult) string {
	if len(result.Content) == 0 {
		return ""
	}

	// For simple text results, return just the text
	if len(result.Content) == 1 && result.Content[0].Type == "text" {
		return result.Content[0].Text
	}

	// For complex results, serialize to JSON
	data, _ := json.Marshal(result)
	return string(data)
}

// ParseMCPToolFromJSON parses an MCP tool from JSON data.
// This is useful when receiving tool definitions from an MCP server.
func ParseMCPToolFromJSON(data []byte) (MCPTool, error) {
	var tool MCPTool
	if err := json.Unmarshal(data, &tool); err != nil {
		return MCPTool{}, fmt.Errorf("failed to parse MCP tool: %w", err)
	}
	return tool, nil
}

// ParseMCPToolsFromJSON parses multiple MCP tools from JSON data.
func ParseMCPToolsFromJSON(data []byte) ([]MCPTool, error) {
	var tools []MCPTool
	if err := json.Unmarshal(data, &tools); err != nil {
		return nil, fmt.Errorf("failed to parse MCP tools: %w", err)
	}
	return tools, nil
}
