package openrouter

import (
	"context"
	"fmt"
	"strings"
)

// Complete sends a legacy completion request to the OpenRouter API.
func (c *Client) Complete(ctx context.Context, prompt string, opts ...CompletionOption) (*CompletionResponse, error) {
	// Validate inputs
	if err := c.validateCompletionRequest(prompt); err != nil {
		return nil, err
	}

	// Build request
	req := &CompletionRequest{
		Model:  c.defaultModel,
		Prompt: prompt,
		Stream: false,
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Validate response format if present
	if req.ResponseFormat != nil {
		if err := validateResponseFormat(req.ResponseFormat); err != nil {
			return nil, err
		}
	}

	// Validate request parameters
	if err := validateCompletionParams(req); err != nil {
		return nil, err
	}

	// Handle model suffixes
	req.Model = processModelSuffix(req.Model, req)

	// Ensure model is set
	if req.Model == "" {
		return nil, ErrNoModel
	}

	// Make request
	var resp CompletionResponse
	err := c.doRequest(ctx, "POST", "/completions", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CompleteStream sends a streaming completion request to the OpenRouter API.
// This method returns a stream that can be used to receive events as they arrive.
func (c *Client) CompleteStream(ctx context.Context, prompt string, opts ...CompletionOption) (*CompletionStream, error) {
	// Validate inputs
	if err := c.validateCompletionRequest(prompt); err != nil {
		return nil, err
	}

	// Build request
	req := &CompletionRequest{
		Model:  c.defaultModel,
		Prompt: prompt,
		Stream: true,
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Validate response format if present
	if req.ResponseFormat != nil {
		if err := validateResponseFormat(req.ResponseFormat); err != nil {
			return nil, err
		}
	}

	// Validate request parameters
	if err := validateCompletionParams(req); err != nil {
		return nil, err
	}

	// Handle model suffixes
	req.Model = processModelSuffix(req.Model, req)

	// Ensure model is set
	if req.Model == "" {
		return nil, ErrNoModel
	}

	// Apply per-request timeout to connection setup only
	connectCtx := ctx
	if req.requestTimeout != nil {
		var cancel context.CancelFunc
		connectCtx, cancel = context.WithTimeout(ctx, *req.requestTimeout)
		defer cancel()
	}

	// Create stream (timeout applies to connection, not stream lifetime)
	stream, err := c.createStream(connectCtx, "/completions", req)
	if err != nil {
		return nil, err
	}

	return &CompletionStream{
		stream: stream,
	}, nil
}

// validateCompletionRequest validates the completion request parameters.
func (c *Client) validateCompletionRequest(prompt string) error {
	if c.apiKey == "" {
		return ErrNoAPIKey
	}

	if prompt == "" {
		return ErrNoPrompt
	}

	return nil
}

// CompleteWithContext is a convenience method that combines prompt completion with context.
func (c *Client) CompleteWithContext(ctx context.Context, contextPrompt, userPrompt string, opts ...CompletionOption) (*CompletionResponse, error) {
	fullPrompt := fmt.Sprintf("%s\n\n%s", contextPrompt, userPrompt)
	return c.Complete(ctx, fullPrompt, opts...)
}

// CompleteWithExamples is a convenience method for few-shot prompting.
func (c *Client) CompleteWithExamples(ctx context.Context, instruction string, examples []string, prompt string, opts ...CompletionOption) (*CompletionResponse, error) {
	var fullPrompt strings.Builder
	fullPrompt.WriteString(instruction)

	if len(examples) > 0 {
		fullPrompt.WriteString("\n\nExamples:\n")
		for i, example := range examples {
			fmt.Fprintf(&fullPrompt, "%d. %s\n", i+1, example)
		}
	}

	fmt.Fprintf(&fullPrompt, "\n\nNow: %s", prompt)

	return c.Complete(ctx, fullPrompt.String(), opts...)
}
