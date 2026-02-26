package openrouter

import (
	"encoding/json"
	"testing"
)

func TestConvertMCPTool_Simple(t *testing.T) {
	mcpTool := MCPTool{
		Name:        "read_file",
		Description: "Reads a file from disk",
		InputSchema: &MCPInputSchema{
			Type: "object",
			Properties: map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Path to the file to read",
				},
			},
			Required: []string{"path"},
		},
	}

	tool := ConvertMCPTool(mcpTool)

	if tool.Type != "function" {
		t.Errorf("expected tool type 'function', got %q", tool.Type)
	}
	if tool.Function.Name != "read_file" {
		t.Errorf("expected function name 'read_file', got %q", tool.Function.Name)
	}
	if tool.Function.Description != "Reads a file from disk" {
		t.Errorf("expected description 'Reads a file from disk', got %q", tool.Function.Description)
	}
	if tool.Function.Parameters["type"] != "object" {
		t.Errorf("expected parameters type 'object', got %v", tool.Function.Parameters["type"])
	}

	// Check properties are copied
	props, ok := tool.Function.Parameters["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties to be a map")
	}
	if _, exists := props["path"]; !exists {
		t.Error("expected 'path' property to exist")
	}

	// Check required is copied
	req, ok := tool.Function.Parameters["required"].([]string)
	if !ok {
		t.Fatal("expected required to be a string slice")
	}
	if len(req) != 1 || req[0] != "path" {
		t.Errorf("expected required ['path'], got %v", req)
	}
}

func TestConvertMCPTool_EmptySchema(t *testing.T) {
	mcpTool := MCPTool{
		Name:        "get_time",
		Description: "Gets the current time",
	}

	tool := ConvertMCPTool(mcpTool)

	if tool.Type != "function" {
		t.Errorf("expected tool type 'function', got %q", tool.Type)
	}
	if tool.Function.Name != "get_time" {
		t.Errorf("expected function name 'get_time', got %q", tool.Function.Name)
	}
	if tool.Function.Parameters["type"] != "object" {
		t.Errorf("expected parameters type 'object', got %v", tool.Function.Parameters["type"])
	}

	props, ok := tool.Function.Parameters["properties"].(map[string]any)
	if !ok || len(props) != 0 {
		t.Errorf("expected empty properties map, got %v", tool.Function.Parameters["properties"])
	}
}

func TestConvertMCPTool_NilProperties(t *testing.T) {
	mcpTool := MCPTool{
		Name: "test_tool",
		InputSchema: &MCPInputSchema{
			Type: "object",
			// Properties is nil
		},
	}

	tool := ConvertMCPTool(mcpTool)

	props, ok := tool.Function.Parameters["properties"].(map[string]any)
	if !ok || len(props) != 0 {
		t.Errorf("expected empty properties map when nil, got %v", tool.Function.Parameters["properties"])
	}
}

func TestConvertMCPTool_AdditionalProperties(t *testing.T) {
	additionalProps := false
	mcpTool := MCPTool{
		Name: "strict_tool",
		InputSchema: &MCPInputSchema{
			Type:                 "object",
			Properties:           map[string]any{},
			AdditionalProperties: &additionalProps,
		},
	}

	tool := ConvertMCPTool(mcpTool)

	if val, ok := tool.Function.Parameters["additionalProperties"].(bool); !ok || val != false {
		t.Errorf("expected additionalProperties false, got %v", tool.Function.Parameters["additionalProperties"])
	}
}

func TestConvertMCPTools_Multiple(t *testing.T) {
	mcpTools := []MCPTool{
		{Name: "tool1", Description: "First tool"},
		{Name: "tool2", Description: "Second tool"},
		{Name: "tool3", Description: "Third tool"},
	}

	tools := ConvertMCPTools(mcpTools)

	if len(tools) != 3 {
		t.Fatalf("expected 3 tools, got %d", len(tools))
	}

	for i, tool := range tools {
		expectedName := mcpTools[i].Name
		if tool.Function.Name != expectedName {
			t.Errorf("tool %d: expected name %q, got %q", i, expectedName, tool.Function.Name)
		}
		if tool.Type != "function" {
			t.Errorf("tool %d: expected type 'function', got %q", i, tool.Type)
		}
	}
}

func TestConvertMCPTools_Empty(t *testing.T) {
	tools := ConvertMCPTools([]MCPTool{})

	if len(tools) != 0 {
		t.Errorf("expected 0 tools, got %d", len(tools))
	}
}

func TestMCPToolSerialization(t *testing.T) {
	mcpTool := MCPTool{
		Name:        "search",
		Description: "Search for items",
		InputSchema: &MCPInputSchema{
			Type: "object",
			Properties: map[string]any{
				"query": map[string]any{"type": "string"},
				"limit": map[string]any{"type": "integer"},
			},
			Required: []string{"query"},
		},
	}

	data, err := json.Marshal(mcpTool)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded MCPTool
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Name != mcpTool.Name {
		t.Errorf("name mismatch: got %q, want %q", decoded.Name, mcpTool.Name)
	}
	if decoded.Description != mcpTool.Description {
		t.Errorf("description mismatch: got %q, want %q", decoded.Description, mcpTool.Description)
	}
}

func TestParseMCPToolFromJSON(t *testing.T) {
	jsonData := []byte(`{
		"name": "calculator",
		"description": "Performs calculations",
		"inputSchema": {
			"type": "object",
			"properties": {
				"expression": {"type": "string"}
			},
			"required": ["expression"]
		}
	}`)

	tool, err := ParseMCPToolFromJSON(jsonData)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if tool.Name != "calculator" {
		t.Errorf("expected name 'calculator', got %q", tool.Name)
	}
	if tool.Description != "Performs calculations" {
		t.Errorf("expected description 'Performs calculations', got %q", tool.Description)
	}
	if tool.InputSchema == nil {
		t.Fatal("expected InputSchema to be non-nil")
	}
	if tool.InputSchema.Type != "object" {
		t.Errorf("expected type 'object', got %q", tool.InputSchema.Type)
	}
}

func TestParseMCPToolFromJSON_Invalid(t *testing.T) {
	jsonData := []byte(`{invalid json}`)

	_, err := ParseMCPToolFromJSON(jsonData)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseMCPToolsFromJSON(t *testing.T) {
	jsonData := []byte(`[
		{"name": "tool1", "description": "First tool"},
		{"name": "tool2", "description": "Second tool"}
	]`)

	tools, err := ParseMCPToolsFromJSON(jsonData)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if len(tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(tools))
	}
	if tools[0].Name != "tool1" {
		t.Errorf("expected first tool name 'tool1', got %q", tools[0].Name)
	}
	if tools[1].Name != "tool2" {
		t.Errorf("expected second tool name 'tool2', got %q", tools[1].Name)
	}
}

func TestParseMCPToolsFromJSON_Invalid(t *testing.T) {
	jsonData := []byte(`not an array`)

	_, err := ParseMCPToolsFromJSON(jsonData)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestConvertToolResultToMCP_SimpleText(t *testing.T) {
	result := MCPToolResult{
		Content: []MCPContent{
			{Type: "text", Text: "Hello, world!"},
		},
	}

	output := ConvertToolResultToMCP(result)
	if output != "Hello, world!" {
		t.Errorf("expected 'Hello, world!', got %q", output)
	}
}

func TestConvertToolResultToMCP_Empty(t *testing.T) {
	result := MCPToolResult{}

	output := ConvertToolResultToMCP(result)
	if output != "" {
		t.Errorf("expected empty string, got %q", output)
	}
}

func TestConvertToolResultToMCP_EmptyContent(t *testing.T) {
	result := MCPToolResult{
		Content: []MCPContent{},
	}

	output := ConvertToolResultToMCP(result)
	if output != "" {
		t.Errorf("expected empty string, got %q", output)
	}
}

func TestConvertToolResultToMCP_Complex(t *testing.T) {
	result := MCPToolResult{
		Content: []MCPContent{
			{Type: "text", Text: "Result 1"},
			{Type: "text", Text: "Result 2"},
		},
	}

	output := ConvertToolResultToMCP(result)

	// Complex results should be JSON
	var parsed MCPToolResult
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}
	if len(parsed.Content) != 2 {
		t.Errorf("expected 2 content items, got %d", len(parsed.Content))
	}
}

func TestConvertToolResultToMCP_NonTextType(t *testing.T) {
	result := MCPToolResult{
		Content: []MCPContent{
			{Type: "image", Data: "base64data", MimeType: "image/png"},
		},
	}

	output := ConvertToolResultToMCP(result)

	// Non-text single items should still be JSON
	var parsed MCPToolResult
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}
	if parsed.Content[0].Type != "image" {
		t.Errorf("expected type 'image', got %q", parsed.Content[0].Type)
	}
}

func TestConvertToolResultToMCP_WithError(t *testing.T) {
	result := MCPToolResult{
		Content: []MCPContent{
			{Type: "text", Text: "Error occurred"},
		},
		IsError: true,
	}

	output := ConvertToolResultToMCP(result)
	// Simple text result should still return just the text
	if output != "Error occurred" {
		t.Errorf("expected 'Error occurred', got %q", output)
	}
}

func TestMCPContentSerialization(t *testing.T) {
	content := MCPContent{
		Type:     "image",
		Data:     "base64encodeddata",
		MimeType: "image/jpeg",
	}

	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded MCPContent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Type != "image" {
		t.Errorf("type mismatch: got %q, want 'image'", decoded.Type)
	}
	if decoded.Data != "base64encodeddata" {
		t.Errorf("data mismatch: got %q, want 'base64encodeddata'", decoded.Data)
	}
	if decoded.MimeType != "image/jpeg" {
		t.Errorf("mimeType mismatch: got %q, want 'image/jpeg'", decoded.MimeType)
	}
}

func TestConvertedToolJSONFormat(t *testing.T) {
	// Test that the converted tool produces the expected JSON format
	mcpTool := MCPTool{
		Name:        "get_weather",
		Description: "Get weather for a location",
		InputSchema: &MCPInputSchema{
			Type: "object",
			Properties: map[string]any{
				"location": map[string]any{
					"type":        "string",
					"description": "City name",
				},
			},
			Required: []string{"location"},
		},
	}

	tool := ConvertMCPTool(mcpTool)

	data, err := json.Marshal(tool)
	if err != nil {
		t.Fatalf("failed to marshal converted tool: %v", err)
	}

	// Verify it can be unmarshaled back
	var decoded Tool
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Type != "function" {
		t.Errorf("expected type 'function', got %q", decoded.Type)
	}
	if decoded.Function.Name != "get_weather" {
		t.Errorf("expected name 'get_weather', got %q", decoded.Function.Name)
	}
}
