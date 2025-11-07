package openrouter

import (
	"fmt"
)

// validateJSONSchema validates a JSON Schema structure for structured outputs.
// It checks:
// 1. Required top-level fields (type)
// 2. Schema structure (properties, items, etc. are correct types)
// 3. Consistency checks (required fields exist in properties, etc.)
func validateJSONSchema(schema map[string]interface{}) error {
	if schema == nil {
		return &ValidationError{
			Field:   "json_schema.schema",
			Message: "schema cannot be nil",
		}
	}

	// Check 1: Required top-level fields
	schemaType, hasType := schema["type"]
	if !hasType {
		return &ValidationError{
			Field:   "json_schema.schema.type",
			Message: "type field is required",
		}
	}

	// Validate type is a string
	typeStr, ok := schemaType.(string)
	if !ok {
		return &ValidationError{
			Field:   "json_schema.schema.type",
			Message: "type must be a string",
		}
	}

	// Validate type value
	validTypes := map[string]bool{
		"object":  true,
		"array":   true,
		"string":  true,
		"number":  true,
		"integer": true,
		"boolean": true,
		"null":    true,
	}
	if !validTypes[typeStr] {
		return &ValidationError{
			Field:   "json_schema.schema.type",
			Message: fmt.Sprintf("invalid type '%s', must be one of: object, array, string, number, integer, boolean, null", typeStr),
		}
	}

	// Check 2 & 3: Type-specific validation
	switch typeStr {
	case "object":
		if err := validateObjectSchema(schema); err != nil {
			return err
		}
	case "array":
		if err := validateArraySchema(schema); err != nil {
			return err
		}
	}

	// Validate nested properties if present
	if properties, hasProps := schema["properties"]; hasProps {
		if err := validateProperties(properties); err != nil {
			return err
		}
	}

	return nil
}

// validateObjectSchema validates object-specific schema fields.
func validateObjectSchema(schema map[string]interface{}) error {
	// Check 2: Validate properties structure if present
	if properties, hasProps := schema["properties"]; hasProps {
		propsMap, ok := properties.(map[string]interface{})
		if !ok {
			return &ValidationError{
				Field:   "json_schema.schema.properties",
				Message: "properties must be an object",
			}
		}

		// Check 3: Validate that required fields exist in properties
		if required, hasRequired := schema["required"]; hasRequired {
			requiredArr, ok := required.([]interface{})
			if !ok {
				// Try []string as well
				requiredStrArr, ok := required.([]string)
				if !ok {
					return &ValidationError{
						Field:   "json_schema.schema.required",
						Message: "required must be an array",
					}
				}
				// Convert to []interface{} for uniform handling
				requiredArr = make([]interface{}, len(requiredStrArr))
				for i, v := range requiredStrArr {
					requiredArr[i] = v
				}
			}

			// Validate each required field exists in properties
			for i, reqField := range requiredArr {
				reqFieldStr, ok := reqField.(string)
				if !ok {
					return &ValidationError{
						Field:   fmt.Sprintf("json_schema.schema.required[%d]", i),
						Message: "required field names must be strings",
					}
				}

				if _, exists := propsMap[reqFieldStr]; !exists {
					return &ValidationError{
						Field:   "json_schema.schema.required",
						Message: fmt.Sprintf("required field '%s' does not exist in properties", reqFieldStr),
					}
				}
			}
		}

		// Validate additionalProperties if present
		if additionalProps, hasAdditional := schema["additionalProperties"]; hasAdditional {
			// additionalProperties can be boolean or object
			switch v := additionalProps.(type) {
			case bool:
				// Valid
			case map[string]interface{}:
				// Valid - it's a schema
				if err := validateJSONSchema(v); err != nil {
					return &ValidationError{
						Field:   "json_schema.schema.additionalProperties",
						Message: fmt.Sprintf("invalid schema: %v", err),
					}
				}
			default:
				return &ValidationError{
					Field:   "json_schema.schema.additionalProperties",
					Message: "additionalProperties must be a boolean or an object (schema)",
				}
			}
		}
	}

	return nil
}

// validateArraySchema validates array-specific schema fields.
func validateArraySchema(schema map[string]interface{}) error {
	// Arrays should have an 'items' field
	items, hasItems := schema["items"]
	if !hasItems {
		// items is not strictly required in JSON Schema, but recommended
		return nil
	}

	// items can be a schema object or array of schemas
	switch v := items.(type) {
	case map[string]interface{}:
		// Single schema for all items
		if err := validateJSONSchema(v); err != nil {
			return &ValidationError{
				Field:   "json_schema.schema.items",
				Message: fmt.Sprintf("invalid schema: %v", err),
			}
		}
	case []interface{}:
		// Tuple validation - array of schemas
		for i, item := range v {
			itemSchema, ok := item.(map[string]interface{})
			if !ok {
				return &ValidationError{
					Field:   fmt.Sprintf("json_schema.schema.items[%d]", i),
					Message: "each item must be a schema object",
				}
			}
			if err := validateJSONSchema(itemSchema); err != nil {
				return &ValidationError{
					Field:   fmt.Sprintf("json_schema.schema.items[%d]", i),
					Message: fmt.Sprintf("invalid schema: %v", err),
				}
			}
		}
	default:
		return &ValidationError{
			Field:   "json_schema.schema.items",
			Message: "items must be a schema object or array of schemas",
		}
	}

	return nil
}

// validateProperties validates all property schemas.
func validateProperties(properties interface{}) error {
	propsMap, ok := properties.(map[string]interface{})
	if !ok {
		return &ValidationError{
			Field:   "json_schema.schema.properties",
			Message: "properties must be an object",
		}
	}

	// Validate each property schema
	for propName, propSchema := range propsMap {
		propSchemaMap, ok := propSchema.(map[string]interface{})
		if !ok {
			return &ValidationError{
				Field:   fmt.Sprintf("json_schema.schema.properties.%s", propName),
				Message: "property schema must be an object",
			}
		}

		// Recursively validate the property schema
		if err := validateJSONSchema(propSchemaMap); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("json_schema.schema.properties.%s", propName),
				Message: fmt.Sprintf("invalid schema: %v", err),
			}
		}
	}

	return nil
}

// validateResponseFormat validates a ResponseFormat structure.
func validateResponseFormat(format *ResponseFormat) error {
	if format == nil {
		return nil
	}

	// Validate type
	validTypes := map[string]bool{
		"json_object": true,
		"json_schema": true,
		"text":        true,
	}

	if !validTypes[format.Type] {
		return &ValidationError{
			Field:   "response_format.type",
			Message: fmt.Sprintf("invalid type '%s', must be one of: json_object, json_schema, text", format.Type),
		}
	}

	// If type is json_schema, validate the schema
	if format.Type == "json_schema" {
		if format.JSONSchema == nil {
			return &ValidationError{
				Field:   "response_format.json_schema",
				Message: "json_schema is required when type is 'json_schema'",
			}
		}

		// Validate JSONSchema fields
		if format.JSONSchema.Name == "" {
			return &ValidationError{
				Field:   "response_format.json_schema.name",
				Message: "name is required",
			}
		}

		if format.JSONSchema.Schema == nil {
			return &ValidationError{
				Field:   "response_format.json_schema.schema",
				Message: "schema is required",
			}
		}

		// Validate the actual schema structure
		if err := validateJSONSchema(format.JSONSchema.Schema); err != nil {
			return err
		}
	}

	return nil
}
