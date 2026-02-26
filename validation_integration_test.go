package openrouter

import (
	"context"
	"testing"
)

// TestChatCompleteWithInvalidSchema verifies that invalid JSON schemas are caught before API calls
func TestChatCompleteWithInvalidSchema(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	tests := []struct {
		name        string
		schema      map[string]any
		expectError bool
		errorField  string
	}{
		{
			name: "missing type field",
			schema: map[string]any{
				"properties": map[string]any{
					"name": map[string]any{
						"type": "string",
					},
				},
			},
			expectError: true,
			errorField:  "json_schema.schema.type",
		},
		{
			name: "required field not in properties",
			schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type": "string",
					},
				},
				"required": []string{"age"},
			},
			expectError: true,
			errorField:  "json_schema.schema.required",
		},
		{
			name: "valid schema",
			schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type": "string",
					},
					"age": map[string]any{
						"type": "integer",
					},
				},
				"required":             []string{"name", "age"},
				"additionalProperties": false,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := []Message{
				CreateUserMessage("Test message"),
			}

			_, err := client.ChatComplete(
				context.Background(),
				messages,
				WithModel("openai/gpt-4"),
				WithJSONSchema("test", true, tt.schema),
			)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}

				valErr, ok := err.(*ValidationError)
				if !ok {
					t.Errorf("expected ValidationError but got %T: %v", err, err)
					return
				}

				if tt.errorField != "" && valErr.Field != tt.errorField {
					t.Errorf("expected error field '%s' but got '%s'", tt.errorField, valErr.Field)
				}
			} else {
				// For valid schemas, we expect either no error (if API call succeeds)
				// or a different error (like network error, auth error, etc.)
				// but NOT a validation error
				if err != nil {
					if valErr, ok := err.(*ValidationError); ok {
						t.Errorf("unexpected validation error: %v", valErr)
					}
					// Other errors (network, auth, etc.) are fine for this test
				}
			}
		})
	}
}

// TestCompleteWithInvalidSchema verifies validation for legacy completion endpoint
func TestCompleteWithInvalidSchema(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"response": map[string]any{
				"type": "string",
			},
		},
		"required": []string{"missing_field"}, // This field doesn't exist in properties
	}

	_, err := client.Complete(
		context.Background(),
		"Test prompt",
		WithCompletionModel("openai/gpt-4"),
		WithCompletionJSONSchema("test", true, schema),
	)

	if err == nil {
		t.Error("expected validation error but got none")
		return
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError but got %T: %v", err, err)
		return
	}

	if valErr.Field != "json_schema.schema.required" {
		t.Errorf("expected error field 'json_schema.schema.required' but got '%s'", valErr.Field)
	}
}

// TestChatCompleteStreamWithInvalidSchema verifies validation for streaming requests
func TestChatCompleteStreamWithInvalidSchema(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	schema := map[string]any{
		// Missing "type" field
		"properties": map[string]any{
			"content": map[string]any{
				"type": "string",
			},
		},
	}

	messages := []Message{
		CreateUserMessage("Test message"),
	}

	_, err := client.ChatCompleteStream(
		context.Background(),
		messages,
		WithModel("openai/gpt-4"),
		WithJSONSchema("test", true, schema),
	)

	if err == nil {
		t.Error("expected validation error but got none")
		return
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError but got %T: %v", err, err)
		return
	}

	if valErr.Field != "json_schema.schema.type" {
		t.Errorf("expected error field 'json_schema.schema.type' but got '%s'", valErr.Field)
	}
}

// TestWithResponseFormatDirectValidation tests using WithResponseFormat directly
func TestWithResponseFormatDirectValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Invalid response format
	format := ResponseFormat{
		Type: "json_schema",
		JSONSchema: &JSONSchema{
			Name: "", // Invalid: name is required
			Schema: map[string]any{
				"type": "object",
			},
		},
	}

	messages := []Message{
		CreateUserMessage("Test message"),
	}

	_, err := client.ChatComplete(
		context.Background(),
		messages,
		WithModel("openai/gpt-4"),
		WithResponseFormat(format),
	)

	if err == nil {
		t.Error("expected validation error but got none")
		return
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError but got %T: %v", err, err)
		return
	}

	if valErr.Field != "response_format.json_schema.name" {
		t.Errorf("expected error field 'response_format.json_schema.name' but got '%s'", valErr.Field)
	}
}

// TestNestedSchemaValidation tests validation of complex nested schemas
func TestNestedSchemaValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Schema with nested object that has invalid required field
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"user": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type": "string",
					},
				},
				"required": []string{"email"}, // email doesn't exist in nested properties
			},
		},
	}

	messages := []Message{
		CreateUserMessage("Test message"),
	}

	_, err := client.ChatComplete(
		context.Background(),
		messages,
		WithModel("openai/gpt-4"),
		WithJSONSchema("test", true, schema),
	)

	if err == nil {
		t.Error("expected validation error but got none")
		return
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError but got %T: %v", err, err)
		return
	}

	// The error should mention the nested property
	if valErr.Field != "json_schema.schema.properties.user" {
		t.Errorf("expected error to be in nested schema, got field '%s'", valErr.Field)
	}
}
