// Package openrouter provides Go bindings for the OpenRouter API.
package openrouter

import (
	"context"
	"fmt"
)

// ResponsesStream represents a streaming response from the Responses API.
type ResponsesStream = Stream[ResponsesResponse]

// CreateResponse sends a request to the OpenRouter Responses API (beta).
// The input can be a simple string or a structured array of ResponsesInputItem.
//
// Note: This API is in beta and may have breaking changes. Use with caution in production.
//
// Example with simple string input:
//
//	resp, err := client.CreateResponse(ctx, "Hello, world!",
//	    openrouter.WithResponsesModel("openai/o4-mini"),
//	    openrouter.WithResponsesMaxOutputTokens(100),
//	)
//
// Example with structured input:
//
//	input := []openrouter.ResponsesInputItem{
//	    openrouter.CreateResponsesUserMessage("Hello!"),
//	}
//	resp, err := client.CreateResponse(ctx, input,
//	    openrouter.WithResponsesModel("openai/o4-mini"),
//	)
func (c *Client) CreateResponse(ctx context.Context, input interface{}, opts ...ResponsesOption) (*ResponsesResponse, error) {
	// Validate input
	if err := c.validateResponsesInput(input); err != nil {
		return nil, err
	}

	// Build request
	req := &ResponsesRequest{
		Model:  c.defaultModel,
		Input:  input,
		Stream: false,
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Validate request
	if err := c.validateResponsesRequest(req); err != nil {
		return nil, err
	}

	// Make request
	var resp ResponsesResponse
	err := c.doRequest(ctx, "POST", "/responses", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateResponseStream sends a streaming request to the OpenRouter Responses API (beta).
// Returns a stream that can be used to receive response events as they arrive.
//
// Note: This API is in beta and may have breaking changes. Use with caution in production.
//
// Example:
//
//	stream, err := client.CreateResponseStream(ctx, "Tell me a story",
//	    openrouter.WithResponsesModel("openai/o4-mini"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer stream.Close()
//
//	for event := range stream.Events() {
//	    // Process streaming event
//	    fmt.Print(event.GetTextContent())
//	}
//
//	if err := stream.Err(); err != nil {
//	    log.Printf("Stream error: %v", err)
//	}
func (c *Client) CreateResponseStream(ctx context.Context, input interface{}, opts ...ResponsesOption) (*ResponsesStream, error) {
	// Validate input
	if err := c.validateResponsesInput(input); err != nil {
		return nil, err
	}

	// Build request
	req := &ResponsesRequest{
		Model:  c.defaultModel,
		Input:  input,
		Stream: true,
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Validate request
	if err := c.validateResponsesRequest(req); err != nil {
		return nil, err
	}

	// Create stream
	stream, err := c.createStream(ctx, "/responses", req)
	if err != nil {
		return nil, err
	}

	return &ResponsesStream{
		stream: stream,
	}, nil
}

// validateResponsesInput validates the input parameter for Responses API requests.
func (c *Client) validateResponsesInput(input interface{}) error {
	if c.apiKey == "" {
		return ErrNoAPIKey
	}

	if input == nil {
		return &ValidationError{
			Field:   "input",
			Message: "input is required",
		}
	}

	switch v := input.(type) {
	case string:
		if v == "" {
			return &ValidationError{
				Field:   "input",
				Message: "input string cannot be empty",
			}
		}
	case []ResponsesInputItem:
		if len(v) == 0 {
			return &ValidationError{
				Field:   "input",
				Message: "input array cannot be empty",
			}
		}
		// Validate each input item
		for i, item := range v {
			if item.Type == "" {
				return &ValidationError{
					Field:   fmt.Sprintf("input[%d].type", i),
					Message: "type is required",
				}
			}
			if item.Type == "message" {
				if item.Role == "" {
					return &ValidationError{
						Field:   fmt.Sprintf("input[%d].role", i),
						Message: "role is required for message type",
					}
				}
				validRoles := map[string]bool{
					"system":    true,
					"user":      true,
					"assistant": true,
				}
				if !validRoles[item.Role] {
					return &ValidationError{
						Field:   fmt.Sprintf("input[%d].role", i),
						Message: fmt.Sprintf("invalid role '%s', must be one of: system, user, assistant", item.Role),
					}
				}
			} else if item.Type == "function_call_output" {
				if item.CallID == "" {
					return &ValidationError{
						Field:   fmt.Sprintf("input[%d].call_id", i),
						Message: "call_id is required for function_call_output type",
					}
				}
			}
		}
	default:
		return &ValidationError{
			Field:   "input",
			Message: "input must be a string or []ResponsesInputItem",
		}
	}

	return nil
}

// validateResponsesRequest validates a ResponsesRequest.
func (c *Client) validateResponsesRequest(req *ResponsesRequest) error {
	if req.Model == "" {
		return ErrNoModel
	}

	// Validate reasoning effort if specified
	if req.Reasoning != nil && req.Reasoning.Effort != "" {
		validEfforts := map[string]bool{
			ReasoningEffortMinimal: true,
			ReasoningEffortLow:     true,
			ReasoningEffortMedium:  true,
			ReasoningEffortHigh:    true,
		}
		if !validEfforts[req.Reasoning.Effort] {
			return &ValidationError{
				Field:   "reasoning.effort",
				Message: fmt.Sprintf("invalid reasoning effort '%s', must be one of: minimal, low, medium, high", req.Reasoning.Effort),
			}
		}
	}

	// Validate tools if specified
	for i, tool := range req.Tools {
		if tool.Type == "" {
			return &ValidationError{
				Field:   fmt.Sprintf("tools[%d].type", i),
				Message: "tool type is required",
			}
		}
		if tool.Function.Name == "" {
			return &ValidationError{
				Field:   fmt.Sprintf("tools[%d].function.name", i),
				Message: "function name is required",
			}
		}
	}

	return nil
}

// Predefined validation errors for Responses API.
var (
	// ErrNoResponsesInput is returned when no input is provided.
	ErrNoResponsesInput = &ValidationError{
		Field:   "input",
		Message: "input is required",
	}
)
