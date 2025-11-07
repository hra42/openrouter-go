package openrouter

import (
	"testing"
)

func TestValidateJSONSchema(t *testing.T) {
	tests := []struct {
		name        string
		schema      map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name:        "nil schema",
			schema:      nil,
			expectError: true,
			errorField:  "json_schema.schema",
		},
		{
			name:        "missing type field",
			schema:      map[string]interface{}{},
			expectError: true,
			errorField:  "json_schema.schema.type",
		},
		{
			name: "type is not a string",
			schema: map[string]interface{}{
				"type": 123,
			},
			expectError: true,
			errorField:  "json_schema.schema.type",
		},
		{
			name: "invalid type value",
			schema: map[string]interface{}{
				"type": "invalid_type",
			},
			expectError: true,
			errorField:  "json_schema.schema.type",
		},
		{
			name: "valid simple object schema",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid object with required fields",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
					"age": map[string]interface{}{
						"type": "integer",
					},
				},
				"required": []string{"name", "age"},
			},
			expectError: false,
		},
		{
			name: "required field not in properties",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"name", "age"},
			},
			expectError: true,
			errorField:  "json_schema.schema.required",
		},
		{
			name: "properties is not an object",
			schema: map[string]interface{}{
				"type":       "object",
				"properties": "invalid",
			},
			expectError: true,
			errorField:  "json_schema.schema.properties",
		},
		{
			name: "required is not an array",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
				},
				"required": "name",
			},
			expectError: true,
			errorField:  "json_schema.schema.required",
		},
		{
			name: "required contains non-string",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []interface{}{123},
			},
			expectError: true,
			errorField:  "json_schema.schema.required[0]",
		},
		{
			name: "additionalProperties as boolean",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
				},
				"additionalProperties": false,
			},
			expectError: false,
		},
		{
			name: "additionalProperties as schema",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
				},
				"additionalProperties": map[string]interface{}{
					"type": "string",
				},
			},
			expectError: false,
		},
		{
			name: "additionalProperties invalid",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
				},
				"additionalProperties": "invalid",
			},
			expectError: true,
			errorField:  "json_schema.schema.additionalProperties",
		},
		{
			name: "valid array schema",
			schema: map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			expectError: false,
		},
		{
			name: "array items is invalid",
			schema: map[string]interface{}{
				"type":  "array",
				"items": "invalid",
			},
			expectError: true,
			errorField:  "json_schema.schema.items",
		},
		{
			name: "nested object schema",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"user": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type": "string",
							},
						},
						"required": []string{"name"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "nested object invalid required",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"user": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type": "string",
							},
						},
						"required": []string{"age"},
					},
				},
			},
			expectError: true,
			errorField:  "json_schema.schema.properties.user",
		},
		{
			name: "valid simple types",
			schema: map[string]interface{}{
				"type": "string",
			},
			expectError: false,
		},
		{
			name: "property schema is not an object",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": "invalid",
				},
			},
			expectError: true,
			errorField:  "json_schema.schema.properties.name",
		},
		{
			name: "array with tuple validation",
			schema: map[string]interface{}{
				"type": "array",
				"items": []interface{}{
					map[string]interface{}{
						"type": "string",
					},
					map[string]interface{}{
						"type": "number",
					},
				},
			},
			expectError: false,
		},
		{
			name: "array tuple item not an object",
			schema: map[string]interface{}{
				"type": "array",
				"items": []interface{}{
					"invalid",
				},
			},
			expectError: true,
			errorField:  "json_schema.schema.items[0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJSONSchema(tt.schema)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				valErr, ok := err.(*ValidationError)
				if !ok {
					t.Errorf("expected ValidationError but got %T", err)
					return
				}
				if tt.errorField != "" && valErr.Field != tt.errorField {
					t.Errorf("expected error field '%s' but got '%s'", tt.errorField, valErr.Field)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateResponseFormat(t *testing.T) {
	tests := []struct {
		name        string
		format      *ResponseFormat
		expectError bool
		errorField  string
	}{
		{
			name:        "nil format",
			format:      nil,
			expectError: false,
		},
		{
			name: "invalid type",
			format: &ResponseFormat{
				Type: "invalid",
			},
			expectError: true,
			errorField:  "response_format.type",
		},
		{
			name: "json_object type",
			format: &ResponseFormat{
				Type: "json_object",
			},
			expectError: false,
		},
		{
			name: "text type",
			format: &ResponseFormat{
				Type: "text",
			},
			expectError: false,
		},
		{
			name: "json_schema without schema",
			format: &ResponseFormat{
				Type: "json_schema",
			},
			expectError: true,
			errorField:  "response_format.json_schema",
		},
		{
			name: "json_schema without name",
			format: &ResponseFormat{
				Type: "json_schema",
				JSONSchema: &JSONSchema{
					Schema: map[string]interface{}{
						"type": "object",
					},
				},
			},
			expectError: true,
			errorField:  "response_format.json_schema.name",
		},
		{
			name: "json_schema without schema field",
			format: &ResponseFormat{
				Type: "json_schema",
				JSONSchema: &JSONSchema{
					Name: "test",
				},
			},
			expectError: true,
			errorField:  "response_format.json_schema.schema",
		},
		{
			name: "valid json_schema",
			format: &ResponseFormat{
				Type: "json_schema",
				JSONSchema: &JSONSchema{
					Name:   "test",
					Strict: true,
					Schema: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type": "string",
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "json_schema with invalid schema",
			format: &ResponseFormat{
				Type: "json_schema",
				JSONSchema: &JSONSchema{
					Name:   "test",
					Strict: true,
					Schema: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type": "string",
							},
						},
						"required": []string{"age"},
					},
				},
			},
			expectError: true,
			errorField:  "json_schema.schema.required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateResponseFormat(tt.format)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				valErr, ok := err.(*ValidationError)
				if !ok {
					t.Errorf("expected ValidationError but got %T", err)
					return
				}
				if tt.errorField != "" && valErr.Field != tt.errorField {
					t.Errorf("expected error field '%s' but got '%s'", tt.errorField, valErr.Field)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateObjectSchema(t *testing.T) {
	tests := []struct {
		name        string
		schema      map[string]interface{}
		expectError bool
	}{
		{
			name: "valid object with properties",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
				},
			},
			expectError: false,
		},
		{
			name: "properties is not a map",
			schema: map[string]interface{}{
				"type":       "object",
				"properties": []string{"name"},
			},
			expectError: true,
		},
		{
			name: "required field exists",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"name"},
			},
			expectError: false,
		},
		{
			name: "required field missing",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"missing"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateObjectSchema(tt.schema)
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateArraySchema(t *testing.T) {
	tests := []struct {
		name        string
		schema      map[string]interface{}
		expectError bool
	}{
		{
			name: "valid array with items",
			schema: map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			expectError: false,
		},
		{
			name: "array without items",
			schema: map[string]interface{}{
				"type": "array",
			},
			expectError: false,
		},
		{
			name: "items is invalid type",
			schema: map[string]interface{}{
				"type":  "array",
				"items": "invalid",
			},
			expectError: true,
		},
		{
			name: "tuple validation",
			schema: map[string]interface{}{
				"type": "array",
				"items": []interface{}{
					map[string]interface{}{
						"type": "string",
					},
					map[string]interface{}{
						"type": "number",
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateArraySchema(tt.schema)
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}
